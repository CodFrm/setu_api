package main

import (
	"github.com/CodFrm/setu_api/internal/pkg/config"
	"github.com/CodFrm/setu_api/internal/rest"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func main() {

	if err := config.InitConfig("config.yaml"); err != nil {
		glog.Fatalf("config init error: %v", err)
	}

	r := gin.New()

	if err := rest.NewWebServer(r); err != nil {
		glog.Fatalf("web server init error: %v", err)
	}

	if err := r.Run(config.AppConfig.WebAddr); err != nil {
		glog.Fatalf("web server error: %v", err)
	}
}
