package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ronething/mp-dev/config"
	"github.com/ronething/mp-dev/controller"
	log "github.com/sirupsen/logrus"
)

func useMiddleWare(e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

func addApi(e *echo.Echo) {
	// wechat 开发者配置 url
	ipValidate := config.Config.GetBool("wechat.ipValidate")
	log.Debugf("ip 白名单校验是否开启 %v", ipValidate)
	if ipValidate {
		e.Any("", controller.ServeWechat, AllowIPList)
	} else {
		e.Any("", controller.ServeWechat)
	}
	// TODO: api sign
	if config.Config.GetBool("server.admin") {
		adminV1 := e.Group("/admin/v1")
		{
			adminV1.GET("/images", controller.ListImages)
		}
	} else {
		log.Infof("admin api 已关闭")
	}
	// 其他服务 需要进行签名验证
}

func CreateEngine() *echo.Echo {
	e := echo.New()
	useMiddleWare(e)
	addApi(e)
	return e
}
