package trie

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type HandlerFunc func(*Context) (*message.Reply, error)

type Context struct {
	Path   string
	Method string
	Params map[string]string
}

func newContext(path string) *Context {
	return &Context{
		Path:   path,
		Method: "GET",
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Text(text string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(text),
	}
}

func (c *Context) Data(reply message.Reply) *message.Reply {
	return &reply
}
