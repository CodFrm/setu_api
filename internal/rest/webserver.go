package rest

import (
	"github.com/CodFrm/setu_api/internal/service"
	"github.com/CodFrm/setu_api/pkg/cache"
	"github.com/gin-gonic/gin"
)

func NewWebServer(r *gin.Engine) error {
	rg := r.Group("/v1/api/pixiv")
	pixiv := newPixiv(service.NewPixiv(cache.NewMapCache(86400 * 2)))
	rg.GET("", pixiv.pic)
	rg.GET("/:param/keyword", pixiv.pic)
	rg.GET("/:param/keyword/relate", pixiv.relateKeyword)
	rg.GET("/:param/relate", pixiv.relatePic)
	rg.GET("/rand", pixiv.rand)
	rg.GET("/:param/small", pixiv.small)
	rg.GET("/:param/original", pixiv.original)

	return nil
}
