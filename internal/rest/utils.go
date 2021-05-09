package rest

import (
	"github.com/CodFrm/setu_api/internal/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handel(ctx *gin.Context, f func() (interface{}, error)) {
	resp, err := f()
	if err != nil {
		switch err.(type) {
		case *errs.RespondError:
			e := err.(*errs.RespondError)
			ctx.JSON(e.StatusCode, gin.H{
				"code": e.Code,
				"msg":  e.Msg,
			})
		case error:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  err.Error(),
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": resp,
	})
}
