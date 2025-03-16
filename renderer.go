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
	return &Renderer{
		instance: instance,
		comp:     comp,
		props:    props,
	}
}

func (r *Renderer) Render() {
	r.mu.Lock()
	defer r.mu.Unlock()
	ctx := context.WithValue(context.Background(), componentInstanceKey, r.instance)
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", html(r.comp, ctx, r.props))
}
