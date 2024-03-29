package pondhttp

import (
	"github.com/e-fish/api/pkg/common/helper/ctxutil"
	pondconfig "github.com/e-fish/api/pond_http/pond_config"
	pondhandler "github.com/e-fish/api/pond_http/pond_handler"
	pondservice "github.com/e-fish/api/pond_http/pond_service"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type route struct {
	conf pondconfig.PondConfig
	gin  *gin.Engine
}

func newRoute(ro route) {
	ginEngine := ro.gin

	service := pondservice.NewService(ro.conf)
	handler := pondhandler.Handler{
		Conf:    ro.conf,
		Service: service,
	}

	ginEngine.POST("/create-pond", ctxutil.Authorization(), handler.CreatePond)
	ginEngine.POST("/update-pond", ctxutil.Authorization(), handler.UpdatePond)
	ginEngine.POST("/resubmission-pond", ctxutil.Authorization(), handler.ResubmissionPond)
	ginEngine.GET("/pond", ctxutil.Authorization(), handler.GetPondByUserAdmin)
	ginEngine.GET("/list-pond", handler.GetAllPond)

	ginEngine.GET("/list-pool", handler.GetListPool)

	ginEngine.GET("/all-pond-submission", ctxutil.Authorization(), handler.GetAllPondSubmission)
	ginEngine.POST("/update-pond-status", ctxutil.Authorization(), handler.UpdatePondStatus)

	ginEngine.POST("/upload-pond-photo", handler.SaveImagePond)
	ginEngine.Use(static.Serve("/assets/image/pond", static.LocalFile(ro.conf.PondImageConfig.Path, false)))

	ginEngine.POST("/upload-pool-photo", handler.SaveImagePool)
	ginEngine.POST("/upload-file", handler.SaveFilePond)
	ginEngine.Use(static.Serve("/assets/image/pool", static.LocalFile(ro.conf.PoolImageConfig.Path, false)))

}
