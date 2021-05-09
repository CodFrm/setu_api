package pixiv

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Pixiv struct {
	cookie string
	client http.Client
}

func NewPixiv(cookie string, client http.Client) *Pixiv {
	return &Pixiv{
		cookie: cookie,
		client: client,
	}
}

func (p *Pixiv) get(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// GetRelatedPicInfo 通过通配id获取相关的图片
func (p *Pixiv) GetRelatedPicInfo(id string) ([]*PixivPicItem, error) {
	ret := make([]*PixivPicItem, 0)
	b, err := p.get("https://www.pixiv.net/ajax/illust/"+id+"/recommend/init?limit=18&lang=zh", map[string]string{
		"Cookie":  p.cookie,
		"Referer": "https://www.pixiv.net/artworks/" + id,
	})
	if err != nil {
		return nil, err
	}
	m := &PixivRecommend{}
	if err := json.Unmarshal(b, m); err != nil {
		return nil, err
	}
	for _, v := range m.Body.Illusts {
		ret = append(ret, &PixivPicItem{Id: v.Id})
	}
	for _, v := range m.Body.NextIds {
		ret = append(ret, &PixivPicItem{Id: v})
	}
	return ret, nil
}

// GetPicInfo 获取图片
func (p *Pixiv) GetPicInfo(id string) (*PixivIllust, error) {
	data, err := p.get("https://www.pixiv.net/ajax/illust/"+id+"?lang=zh", map[string]string{
		"Cookie": p.cookie,
	})
	if err != nil {
		return nil, err
	}
	m := &PixivIllust{}
	if err := json.Unmarshal(data, m); err != nil {
		return nil, err
	}
	return m, nil
}

// DownloadSmallPic 下载图片
func (p *Pixiv) DownloadSmallPic(info *PixivIllust, w io.Writer) error {
	data, err := p.get(info.Body.Urls.Small, map[string]string{
		"Cookie":  p.cookie,
		"Referer": "https://www.pixiv.net/artworks/" + info.Id,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(data))
	return err
}

// DownloadOriginalPic 下载图片
func (p *Pixiv) DownloadOriginalPic(info *PixivIllust, w io.Writer) error {
	data, err := p.get(info.Body.Urls.Original, map[string]string{
		"Cookie":  p.cookie,
		"Referer": "https://www.pixiv.net/artworks/" + info.Id,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(data))
	return err
}

var Hots = []string{"30000", "20000", "10000", "5000", "1000", "500"}

func (p *Pixiv) Tagurlencode(tag string, hot int) string {
	return strings.ReplaceAll(url.QueryEscape(tag+" "+Hots[hot]+"users入り"), "+", "%20")
}

func (p *Pixiv) List(tag string, page int) ([]*PixivPicItem, error) {
	str, err := p.get("https://www.pixiv.net/ajax/search/illustrations/"+tag+
		"?word="+tag+"&order=date_d&mode=safe&p="+strconv.Itoa(page)+"&s_mode=s_tag&type=illust_and_ugoira&lang=zh",
		map[string]string{
			"Cookie":  p.cookie,
			"Referer": "https://www.pixiv.net/tags/" + tag + "/illustrations?s_mode=s_tag",
		})
	if err != nil {
		return nil, err
	}
	m := &IllustRespond{}
	if err := json.Unmarshal(str, m); err != nil {
		return nil, err
	}
	return m.Body.Illust.Data, nil
}

func (p *Pixiv) RankList(mode string, page int) ([]*PixivPicItem, error) {
	str, err := p.get("https://www.pixiv.net/ranking.php?mode="+mode+"&content=illust&p="+strconv.Itoa(page)+"&format=json",
		map[string]string{
			"Cookie":  p.cookie,
			"Referer": "https://www.pixiv.net/ranking.php?mode=" + mode + "&content=illust",
		})
	if err != nil {
		return nil, err
	}
	m := &PixivRankList{}
	if err := json.Unmarshal(str, m); err != nil {
		return nil, err
	}
	ret := make([]*PixivPicItem, 0)
	for _, v := range m.Contents {
		ret = append(ret, &PixivPicItem{
			Id:              strconv.Itoa(v.IllustId),
			ProfileImageUrl: v.ProfileImg,
			Url:             v.Url,
			UserId:          strconv.Itoa(v.UserId),
			UserName:        v.UserName,
			Title:           v.Title,
		})
	}
	return ret, nil
}

func (p *Pixiv) GetRelateTags(tag string) ([]string, error) {
	str, err := p.get("https://www.pixiv.net/rpc/cps.php?keyword="+url.QueryEscape(tag)+"&lang=zh",
		map[string]string{
			"Cookie":  p.cookie,
			"Referer": "https://www.pixiv.net/",
		})
	if err != nil {
		return nil, err
	}
	m := &PixivTags{}
	if err := json.Unmarshal(str, m); err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, v := range m.Candidates {
		ret = append(ret, v.TagName)
	}
	return ret, nil
}
