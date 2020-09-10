package trie

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type Router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

//NewRouter
func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) AddRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *Router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *Router) PrintRoutes(method string) string {
	nodes := r.getRoutes(method)
	var text bytes.Buffer
	for _, node := range nodes {
		text.WriteString(fmt.Sprintf("- %s\n", node.pattern))
	}
	return text.String()
}

func (r *Router) Handle(path string) *message.Reply {
	c := newContext(path)
	n, params := r.getRoute("GET", path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		reply, err := r.handlers[key](c)
		if err != nil {
			return c.Text("服务端发生错误")
		}
		return reply
	} else {
		return c.Text("没有匹配到路由")
	}
}
