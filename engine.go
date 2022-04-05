package neo

import (
	"net/http"
	"sync"
)

type HandlerFunc func(ctx *Context)

type HandlerChain []HandlerFunc

type Engine struct {
	*router
	ctxPool sync.Pool
}

func NewEngine() *Engine {
	return &Engine{
		router: newRouter(),
		ctxPool: sync.Pool{
			New: func() interface{} {
				return &Context{}
			},
		},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.allocateContext(w, r)
	defer e.ctxPool.Put(ctx)
	e.handle(ctx)
}

func (e *Engine) Use(middlewares ...HandlerFunc) *Engine {
	e.router.Use(middlewares...)
	return e
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) allocateContext(w http.ResponseWriter, r *http.Request) (ctx *Context) {
	obj := e.ctxPool.Get()
	ctx = obj.(*Context)
	ctx.reset(w, r)
	return
}
