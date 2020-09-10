package storage

import (
	"github.com/ronething/mp-dev/storage/trie"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func HelpUsage(c *trie.Context) (*message.Reply, error) {
	text := wechatRouter.PrintRoutes("GET")
	usage := "usage:\n" + text
	return c.Text(usage), nil
}
