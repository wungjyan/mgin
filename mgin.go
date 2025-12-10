// mgin 是一个极简的 Web 框架雏形，仿照 gin 的路由注册与启动流程
package mgin

import "net/http"

// HandlerFunc 定义了请求处理函数的签名，约定处理器接收 ResponseWriter 与 *Request
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 是框架的核心结构体，保存路由表等运行时信息
// router 使用 "METHOD-PATH" 组合作为键，映射到对应的处理函数
type Engine struct {
	router map[string]HandlerFunc
}

// New 创建并返回一个 Engine 实例，同时初始化路由表
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRouter 将指定的 HTTP 方法与路径模式绑定到处理器，并写入路由表
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router[method+"-"+pattern] = handler
}

// GET 注册一个 GET 路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

// POST 注册一个 POST 路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// Run 启动 HTTP 服务器，传入 Engine 作为 http.Handler 接管所有请求
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

/*
ServeHTTP 是 Engine 实现 `http.Handler` 接口的关键方法。
`http.ListenAndServe` 的第二个参数需要一个 `http.Handler`：

	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}

当 Engine 提供了 `ServeHTTP` 并作为参数传入 `ListenAndServe` 时，
它会在每个请求到来时被调用，从而接管所有 HTTP 请求的处理流程。
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 构造路由匹配键：方法与路径拼接，例如 "GET-/"
	key := req.Method + "-" + req.URL.Path
	// 查找路由表，若存在则调用对应处理器
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		// 未匹配到处理器时返回 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 NOT FOUND"))
	}
}
