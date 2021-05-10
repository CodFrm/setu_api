package service

import (
	"bytes"
	"github.com/CodFrm/setu_api/internal/errs"
	"github.com/CodFrm/setu_api/pkg/cache"
	pixiv2 "github.com/CodFrm/setu_api/pkg/pixiv"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var PicIsNil = errs.NewRespondError(http.StatusOK, 1000, "我真的一张都没有了")

type Pixiv interface {
	GetPicInfo(user string, keyword string) (*pixiv2.PixivPicItem, error)
	GetRelatePicInfo(user, id string) (*pixiv2.PixivPicItem, error)
	Download(user, id string, small bool, io io.Writer) error
}

type pixiv struct {
	pixiv *pixiv2.Pixiv
	cache cache.Cache
}

func NewPixiv(cache cache.Cache) Pixiv {
	// 创建缓存目录
	if err := os.MkdirAll("./runtime/pic/small", 0644); err != nil {
		glog.Fatalf("create cache dir error: %v", err)
	}
	if err := os.MkdirAll("./runtime/pic/original", 0644); err != nil {
		glog.Fatalf("create cache dir error: %v", err)
	}
	return &pixiv{
		pixiv: pixiv2.NewPixiv("", http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}),
		cache: cache,
	}
}

func (p *pixiv) GetPicInfo(user string, keyword string) (*pixiv2.PixivPicItem, error) {
	list := make([]*pixiv2.PixivPicItem, 0)
	page := 0
	if err := p.cache.GetOrSet(cache.Join("pixiv:user:rank:page", user, keyword), &page, func() (interface{}, error) {
		return 1, nil
	}); err != nil {
		return nil, err
	}
	if keyword == "" {
		// 查询当前用户上一次的页面
		if err := p.cache.GetOrSet(cache.Join("pixiv:user:rank:page:list", strconv.Itoa(page)), &list, func() (interface{}, error) {
			list, err := p.pixiv.RankList("weekly", page)
			if err != nil {
				return nil, err
			}
			return list, nil
		}); err != nil {
			return nil, err
		}
	} else {
		if err := p.cache.GetOrSet(cache.Join("pixiv:user:rank:keyword:page:list", keyword, strconv.Itoa(page)), &list, func() (interface{}, error) {
			list, err := p.pixiv.RankList("weekly", page)
			if err != nil {
				return nil, err
			}
			return list, nil
		}); err != nil {
			return nil, err
		}
	}
	if len(list) == 0 {
		return nil, PicIsNil
	}
	ret, err := p.uniqueRand(user, list)
	if err != nil {
		if err == PicIsNil {
			if err := p.cache.Set(cache.Join("pixiv:user:rank:page", user, keyword), page+1); err != nil {
				return nil, err
			}
			return p.GetPicInfo(user, keyword)
		}
		return nil, err
	}
	return ret, nil
}

//5天内不再重复
func (p *pixiv) uniqueRand(user string, data []*pixiv2.PixivPicItem) (*pixiv2.PixivPicItem, error) {
	randList := make([]*pixiv2.PixivPicItem, 0)
	for _, v := range data {
		if _, err := p.cache.Get(cache.Join("pixiv:unique", user, v.Id)); err == cache.ErrNotExist {
			randList = append(randList, v)
		}
	}
	if len(randList) == 0 {
		return nil, PicIsNil
	}
	ret := randList[rand.Intn(len(randList))]
	if err := p.cache.Set(cache.Join("pixiv:unique", user, ret.Id), "1", cache.Expire(86400*5)); err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *pixiv) GetRelatePicInfo(user, id string) (*pixiv2.PixivPicItem, error) {
	list := make([]*pixiv2.PixivPicItem, 0)
	if err := p.cache.GetOrSet(cache.Join("pixiv:relate:id", id), &list, func() (interface{}, error) {
		return p.pixiv.GetRelatedPicInfo(id)
	}, cache.Expire(86400)); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, PicIsNil
	}
	ret, err := p.uniqueRand(user, list)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *pixiv) Download(user, id string, small bool, w io.Writer) error {
	info := &pixiv2.PixivIllust{}
	err := p.cache.GetOrSet(cache.Join("pixiv:pic:info", id), info, func() (interface{}, error) {
		return p.pixiv.GetPicInfo(id)
	})
	if err != nil {
		return err
	}
	path := picDri(id, small)
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		buf := bytes.NewBuffer(nil)
		// 不存在,下载到本地
		if small {
			err = p.pixiv.DownloadSmallPic(info, buf)
		} else {
			err = p.pixiv.DownloadOriginalPic(info, buf)
		}
		if err != nil {
			return err
		}
		err := ioutil.WriteFile(path, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	return err
}
