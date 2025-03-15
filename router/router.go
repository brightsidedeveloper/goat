package goatRouter

import (
	"strings"
	"syscall/js"
)

type Router struct {
	routes map[string]func(params map[string]string)
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]func(params map[string]string)),
	}
}

func (r *Router) Register(path string, handler func(params map[string]string)) {
	cleanPath := strings.TrimPrefix(path, "/")
	r.routes[cleanPath] = handler
}

func (r *Router) RegisterGuarded(path string, handler func(params map[string]string), guard func() bool, notAuthorizedPath string) {
	r.routes[path] = func(params map[string]string) {
		if guard() {
			handler(params)
		} else {
			r.Navigate(notAuthorizedPath)
		}
	}
}

func (r *Router) Navigate(path string) {
	cleanPath := strings.TrimPrefix(path, "/")

	if handler, exists := r.routes[cleanPath]; exists {
		handler(nil)
		return
	}

	for routePattern, handler := range r.routes {
		if !strings.Contains(routePattern, ":") {
			continue
		}

		patternParts := strings.Split(routePattern, "/")
		pathParts := strings.Split(cleanPath, "/")

		if len(patternParts) != len(pathParts) {
			continue
		}

		params := make(map[string]string)
		matches := true

		for i, patternPart := range patternParts {
			if strings.HasPrefix(patternPart, ":") {
				paramName := strings.TrimPrefix(patternPart, ":")
				params[paramName] = pathParts[i]
			} else if patternPart != pathParts[i] {
				matches = false
				break
			}
		}

		if matches {
			handler(params)
			return
		}
	}

	if handler, exists := r.routes["404"]; exists {
		handler(nil)
	}
}

func (r *Router) SetupEventListeners() {

	// Handle Links
	js.Global().Get("document").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		target := e.Call("target").Call("closest", "a")
		if !target.IsNull() {
			href := target.Get("pathname").String()
			if strings.HasPrefix(href, "/") {
				e.Call("preventDefault")
				js.Global().Get("history").Call("pushState", nil, "", href)
				r.Navigate(href)
			}
		}
		return nil
	}))

	// Handle Pops
	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		path := js.Global().Get("location").Get("pathname").String()
		r.Navigate(path)
		return nil
	}))

	// Init
	path := js.Global().Get("location").Get("pathname").String()
	r.Navigate(path)
}
