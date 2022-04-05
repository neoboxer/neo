package neo

import (
	"net/http"
)

type router struct {
	roots         map[string]*trie
	middlewares   HandlerChain
	handlerChains map[string]HandlerChain
}

func newRouter() *router {
	return &router{
		roots:         map[string]*trie{},
		middlewares:   nil,
		handlerChains: map[string]HandlerChain{},
	}
}

// addRoute with request method and path to trie
func (r *router) addRoute(method, path string, handlers ...HandlerFunc) {
	assert(method != "", "HTTP method can not be empty")
	assert(len(handlers) != 0, "there must be at least one handler")
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trie{
			son: map[string]*trie{},
		}
	}
	root := r.roots[method]
	root.insert(path)
	key := method + "-" + path
	r.handlerChains[key] = r.getChain(handlers)
}

// getRoute from the trie with request method
func (r *router) getRoute(method, path string) (*trie, map[string]string) {
	if root, ok := r.roots[method]; !ok {
		return nil, nil
	} else {
		return root.search(path)
	}
}

func (r *router) getChain(handlers HandlerChain) (handlerChain HandlerChain) {
	mdSize := len(r.middlewares)
	totalSize := mdSize + len(handlers)
	handlerChain = make(HandlerChain, totalSize)
	copy(handlerChain, r.middlewares)
	copy(handlerChain[mdSize:], handlers)
	return
}

func (r *router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = params
		key := ctx.Method + "-" + node.path
		if handlers := r.handlerChains[key]; len(handlers) != 0 {
			ctx.handlers = handlers
			ctx.Next()
		}
	} else {
		ctx.String(http.StatusNotFound, "404 not found")
	}
}

func (r *router) Use(middlewares ...HandlerFunc) *router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *router) GET(pattern string, handlers ...HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handlers...)
}

func (r *router) PUT(pattern string, handlers ...HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handlers...)
}

func (r *router) POST(pattern string, handlers ...HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handlers...)
}

func (r *router) DELETE(pattern string, handlers ...HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handlers...)
}
