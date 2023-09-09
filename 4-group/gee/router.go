package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	ps := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range ps {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *router) handle(ctx *Context) {
	node, params := r.getRoute(ctx.Method, ctx.Path)
	if node != nil {
		ctx.Params = params
		key := ctx.Method + "-" + node.pattern
		handler, ok := r.handlers[key]
		if ok {
			handler(ctx)
		}
		//r.handlers[key](ctx)
	} else {
		ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
	}
}
