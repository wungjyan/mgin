// 路由器负责将 METHOD-PATH 映射到具体的处理函数，并在请求到来时进行分发
package mgin

import "net/http"

// router 保存路由处理器的映射关系
// handlers 的键使用 "METHOD-PATH" 组合（例如 "GET-/"）来标识唯一的路由
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter 初始化一个空的路由器实例
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRouter 注册路由，将 HTTP 方法与路径模式组合为键，绑定到处理器
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle 根据请求的 Method 与 Path 查找处理器，匹配成功则执行
// 未匹配到时返回 404 提示
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
