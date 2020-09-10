package storage

import (
	"time"

	"github.com/imroc/req"
	config2 "github.com/ronething/mp-dev/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	log "github.com/sirupsen/logrus"
)

var WechatOfficialAccount *officialaccount.OfficialAccount
var memory cache.Cache

const callbackIPList = "callbackIPList"

type (
	IPList struct {
		IPList []string `json:"ip_list"`
	}
)

func InitWechatConfig() {
	wc := wechat.NewWechat()
	//这里本地内存保存 access_token，也可选择 redis，memcache或者自定cache
	memory = cache.NewMemory()
	cfg := &config.Config{
		AppID:          config2.Config.GetString("wechat.AppID"),
		AppSecret:      config2.Config.GetString("wechat.AppSecret"),
		Token:          config2.Config.GetString("wechat.Token"),
		EncodingAESKey: config2.Config.GetString("wechat.EncodingAESKey"),
		Cache:          memory,
	}
	WechatOfficialAccount = wc.GetOfficialAccount(cfg)
}

//GetWechatIPList
func GetWechatIPList() ([]string, error) {
	if memory.IsExist(callbackIPList) {
		log.Debugf("进入缓存")
		return memory.Get(callbackIPList).([]string), nil
	}
	var (
		token string
		resp  *req.Resp
		err   error
	)
	token, err = WechatOfficialAccount.GetAccessToken()
	if err != nil {
		log.Errorf("获取 token 失败, err:%+v", err)
		return nil, err
	}

	if resp, err = httpClient.Get(
		"https://api.weixin.qq.com/cgi-bin/getcallbackip", req.QueryParam{
			"access_token": token,
		}); err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("get wechat ip list resp is %s", resp.String())
	var ipList IPList
	if err = resp.ToJSON(&ipList); err != nil {
		log.Error(err)
		return nil, err
	}
	// 缓存
	log.Debugf("设置缓存")
	if err = memory.Set(callbackIPList, ipList.IPList, 1*time.Hour); err != nil {
		log.Warnf("设置缓存失败")
		//return nil, err
	}
	return ipList.IPList, nil
}
