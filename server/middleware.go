package server

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/ronething/mp-dev/storage"
	log "github.com/sirupsen/logrus"
)

//AllowIPList
//DONE: 生产环境 IP 白名单过滤
func AllowIPList(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		realIP := c.RealIP()
		log.Debugf("real ip is %s", realIP)
		ips, err := storage.GetWechatIPList()
		if err != nil {
			log.Errorf("获取 ip 白名单失败")
			return err
		}
		log.Debugf("ips is %+v", ips)
		if len(ips) == 0 {
			return c.JSON(http.StatusForbidden, "not allow")
		}
		sort.Strings(ips)
		index := sort.SearchStrings(ips, realIP)
		found := index < len(ips) && ips[index] == realIP
		log.Debugf("found or not %v", found)
		if found {
			return next(c)
		} else {
			// 走到这里说明不匹配
			return c.JSON(http.StatusForbidden, "not allow")
		}
	}
}
