package goat

import (
	"context"
	"sync"
	"syscall/js"
)

type Renderer struct {
	instance *ComponentInstance
	comp     Component
	props    any
	mu       sync.Mutex
}

func NewRenderer(comp Component, props any) *Renderer {
	instance := &ComponentInstance{
		states:     make(map[int]any),
		stateOrder: []int{},
		callIndex:  0,
	}
	r := &Renderer{
		instance: instance,
		comp:     comp,
		props:    props,
	}
	registerRenderer(instance, r)
	return r
}

func (r *Renderer) Render() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instance.mu.Lock()
	r.instance.callIndex = 0
	r.instance.mu.Unlock()
	ctx := context.WithValue(context.Background(), componentInstanceKey, r.instance)
	ctx = context.WithValue(ctx, propsKey, r.props)
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", html(r.comp, ctx, r.props))
}

type propsKeyType struct{}

var propsKey = propsKeyType{}

func GetProps(ctx context.Context) any {
	if props, ok := ctx.Value(propsKey).(any); ok {
		return props
	}
	panic("No props found in context")
}
