package goat

import (
	"context"
	"sync"
	"syscall/js"
)

type Renderer struct {
	id       string
	vdom     VNode
	instance *ComponentInstance
	comp     Component
	props    any
	mu       sync.Mutex
}

func NewRenderer(id string, comp Component, props any) *Renderer {
	instance := &ComponentInstance{
		states:     make(map[int]any),
		stateOrder: []int{},
		callbacks:  make(map[int]string),
		callIndex:  0,
	}
	r := &Renderer{
		instance: instance,
		comp:     comp,
		props:    props,
		id:       id,
	}
	registerRenderer(instance, r)
	return r
}

func (r *Renderer) Render() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instance.Reset()
	ctx := context.WithValue(context.Background(), componentInstanceKey, r.instance)
	ctx = context.WithValue(ctx, propsKey, r.props)
	newVdom := r.comp.Render(ctx, r.props)
	r.updateDOM(newVdom)
	r.vdom = newVdom
}

func (r *Renderer) updateDOM(newVdom VNode) {
	doc := js.Global().Get("document")
	container := doc.Call("getElementById", r.id)
	if !container.Truthy() {
		// Initial render: create the element
		container = doc.Call("createElement", "div")
		container.Set("id", r.id)
		doc.Get("body").Call("appendChild", container)
	}
	// For now, just set innerHTML (no diffing yet)
	html := vdomToHTML(newVdom)
	container.Set("innerHTML", html)
}

func vdomToHTML(v VNode) string {
	if v.Text != "" {
		return v.Text
	}
	return "<div>" + v.Text + "</div>" // Placeholder
}

type propsKeyType struct{}

var propsKey = propsKeyType{}

func GetProps(ctx context.Context) any {
	if props, ok := ctx.Value(propsKey).(any); ok {
		return props
	}
	panic("No props found in context")
}
