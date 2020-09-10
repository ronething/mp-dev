package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ronething/mp-dev/storage"
	log "github.com/sirupsen/logrus"
)

//ServeWechat http api
func ServeWechat(c echo.Context) (err error) {
	// 传入 request 和 responseWriter
	server := storage.WechatOfficialAccount.GetServer(c.Request(), c.Response().Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(storage.MsgHandler)

	//处理消息接收以及回复
	if err = server.Serve(); err != nil {
		log.Errorf("[serveWechat] err: %s", err.Error())
		return c.JSON(http.StatusForbidden, nil)
	}
	//发送回复的消息
	return server.Send()
}

//ListImages
func ListImages(c echo.Context) (err error) {
	m := storage.WechatOfficialAccount.GetMaterial()
	list, err := m.BatchGetMaterial("image", 0, 10)
	if err != nil {
		log.Errorf("err:%v", err)
		return c.JSON(http.StatusOK, "error")
	}
	return c.JSON(http.StatusOK, list)
}
