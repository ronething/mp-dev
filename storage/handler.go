package storage

import (
	"time"

	"github.com/imroc/req"
	"github.com/ronething/mp-dev/storage/trie"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
)

var wechatRouter *trie.Router
var httpClient *req.Req

func init() {
	InitWechatHandlerRouter()
	InitHttpClient()
}

func InitHttpClient() {
	httpClient = req.New()
	// 4.5 秒超时 微信会控制 5s 超时 返回服务不可用 这里应该是对单个连接有效
	httpClient.SetTimeout(4500 * time.Millisecond)
}

// 根据消息内容进行路由判断
func InitWechatHandlerRouter() {
	wechatRouter = trie.NewRouter()
	wechatRouter.AddRoute("GET", "/help", HelpUsage)
	wechatRouter.AddRoute("GET", "/music/play/:sid", PlayMusicBySongId)
	wechatRouter.AddRoute("GET", "/music/url/:sid", GetSongURL)
	wechatRouter.AddRoute("GET", "/music/search/:keywords", SearchMusicByKeyword)
	wechatRouter.AddRoute("GET", "/music/search/:keywords/:page", SearchMusicByKeyword)
	wechatRouter.AddRoute("GET", "/music/:name", PlayMusicByName)
}

//MsgHandler http 实际处理函数
func MsgHandler(msg message.MixMessage) *message.Reply {
	switch msg.MsgType {
	case message.MsgTypeText:
		log.Debugf("用户发送消息正文: %v", msg.Content)
		return wechatRouter.Handle(msg.Content)
	case message.MsgTypeEvent:
		if msg.Event == message.EventSubscribe {
			return &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText("感谢关注🙏"),
			}
		} else {
			return &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText(""),
			}
		}
	default:
		text := message.NewText("暂不支持其他类型消息，请输入文本")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
