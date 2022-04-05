package neo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type M map[string]interface{}

type Context struct {
	Writer       http.ResponseWriter
	Req          *http.Request
	Path         string            // router path
	Method       string            // request method
	StatusCode   int               // http response status code
	Params       map[string]string // mark wildcard path
	handlers     HandlerChain      // current context handle chain
	handlerIndex int               // current handler chain index
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Req = r
	c.Path = r.URL.Path
	c.Method = r.Method
	c.StatusCode = http.StatusOK
	c.Params = map[string]string{}
	c.handlers = nil
	c.handlerIndex = -1
}

func (c *Context) Next() {
	c.handlerIndex++
	if c.handlerIndex < len(c.handlers) {
		c.handlers[c.handlerIndex](c)
		c.handlerIndex++
	}
}

// PostForm acquire specific field from form
// 获取表达你对应字段的值
func (c *Context) PostForm(key string) string {
	// acquire the first field from Request.Form
	return c.Req.FormValue(key)
}

// Query 获取参数对应字段的值
// acquire first value use r.URL.Query().Get("ParamName")
// acquire all value use r.URL.Query()["ParamName"]
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetHeader(key, value string) {
	// A Header represents the key-value pairs in an HTTP header.
	c.Writer.Header().Set(key, value)
}

// Status store status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	// WriteHeader sends an HTTP response header with the provided status code.
	c.Writer.WriteHeader(code)
}

// JSON 封装Json类型的响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// String 封装String类型的响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data 封装Data类型的响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 封装Html类型的响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
