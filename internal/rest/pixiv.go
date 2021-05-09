package rest

import (
	"github.com/CodFrm/setu_api/internal/service"
	"github.com/gin-gonic/gin"
)

type pixiv struct {
	service.Pixiv
}

func newPixiv(api service.Pixiv) *pixiv {
	return &pixiv{
		Pixiv: api,
	}
}

// 通过关键字,随机获取一张图片
func (p *pixiv) pic(ctx *gin.Context) {
	handel(ctx, func() (interface{}, error) {
		keyword := ctx.Param("param")
		u, err := p.userid(ctx)
		if err != nil {
			return nil, err
		}
		item, err := p.GetPicInfo(u, keyword)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}

// 通过id获取相关图片
func (p *pixiv) relatePic(ctx *gin.Context) {
	handel(ctx, func() (interface{}, error) {
		pid := ctx.Param("param")
		u, err := p.userid(ctx)
		if err != nil {
			return nil, err
		}
		item, err := p.GetRelatePicInfo(u, pid)
		if err != nil {
			return nil, err
		}
		return item, nil
	})
}

// 获取用户id,防止重复
func (p *pixiv) userid(ctx *gin.Context) (string, error) {
	token := ctx.Request.Header.Get("X-Setu-Token")
	if token != "" {
		return token, nil
	}
	return ctx.ClientIP(), nil
}

func (p *pixiv) relateKeyword(ctx *gin.Context) {

}

func (p *pixiv) rand(ctx *gin.Context) {

}

func (p *pixiv) small(ctx *gin.Context) {
	handel(ctx, func() (interface{}, error) {
		pid := ctx.Param("param")
		u, err := p.userid(ctx)
		if err != nil {
			return nil, err
		}
		return nil, p.Pixiv.Download(u, pid, true, ctx.Writer)
	})
}

func (p *pixiv) original(ctx *gin.Context) {
	handel(ctx, func() (interface{}, error) {
		pid := ctx.Param("param")
		u, err := p.userid(ctx)
		if err != nil {
			return nil, err
		}
		return nil, p.Pixiv.Download(u, pid, false, ctx.Writer)
	})
}
