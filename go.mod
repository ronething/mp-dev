module github.com/ronething/mp-dev

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/imroc/req v0.3.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/pkg/errors v0.8.1
	github.com/silenceper/wechat/v2 v2.0.2
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
)

replace github.com/silenceper/wechat/v2 v2.0.2 => github.com/ronething/wechat/v2 v2.0.4-beta.1
