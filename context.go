// Context 封装了请求和响应的上下文，提供读取参数、设置响应头与写出多种格式的便捷方法
package mgin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 是一个便捷的 Map 类型，常用于构造 JSON 响应体
type H map[string]interface{}

type Context struct {
	// 原始的 ResponseWriter 和 Request
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息，便于快速访问无需每次从 Req 中解析
	Path   string
	Method string
	// 响应状态码
	StatusCode int
}

// newContext 从原始的 http 请求/响应构造 Context，并缓存常用的路径与方法
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm 读取表单字段（application/x-www-form-urlencoded 或 multipart/form-data）
// 若同名字段同时存在于 Body 与 Query，FormValue 会优先返回 Body 中的值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 读取 URL 查询参数（例如 /path?key=value）
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应状态码，并记录到 Context 中
// 注意：首次调用 WriteHeader 后，状态码与响应头将被发送，后续无法再修改
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头字段，应在写入 Body 或 WriteHeader 之前调用
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 输出纯文本内容，并设置合适的 Content-Type 与状态码
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 将对象序列化为 JSON 并写入响应体；序列化失败时返回 500 错误
func (c *Context) JSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.Status(code)
	c.Writer.Write(data)
}

// Data 直接写出字节数据，Content-Type 由调用方自行控制（可结合 SetHeader 使用）
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 写出 HTML 字符串响应，并设置 text/html 的 Content-Type
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
