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
	id       string
	vdom     VNode
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
	if !r.vdom.DOMNode.Truthy() {
		// Initial render
		r.vdom = createDOM(newVdom, doc)
		container.Call("appendChild", r.vdom.DOMNode)
	} else {
		// Diff and patch
		diffAndPatch(r.vdom, newVdom, r.vdom.DOMNode, doc)
	}
}

func createDOM(v VNode, doc js.Value) VNode {
	var node js.Value
	if v.Tag == "" {
		node = doc.Call("createTextNode", v.Text)
	} else {
		node = doc.Call("createElement", v.Tag)
		for key, value := range v.Attrs {
			node.Call("setAttribute", key, value)
		}
		for event, handler := range v.Events {
			node.Call("addEventListener", event, js.FuncOf(handler))
		}
		children := make([]VNode, len(v.Children))
		for i, child := range v.Children {
			childNode := createDOM(child, doc)
			node.Call("appendChild", childNode.DOMNode)
			children[i] = childNode
		}
		v.Children = children
		if v.Text != "" {
			node.Set("textContent", v.Text)
		}
	}
	v.DOMNode = node
	return v
}

func diffAndPatch(oldVdom, newVdom VNode, parent js.Value, doc js.Value) {
	if oldVdom.Tag != newVdom.Tag || oldVdom.Text != newVdom.Text {
		// Replace node if tag or text differs
		newNode := createDOM(newVdom, doc)
		parent.Call("replaceChild", newNode.DOMNode, oldVdom.DOMNode)
		oldVdom.DOMNode = newNode.DOMNode
		return
	}

	// Update attributes
	for key, value := range newVdom.Attrs {
		if oldVdom.Attrs[key] != value {
			oldVdom.DOMNode.Call("setAttribute", key, value)
		}
	}
	for key := range oldVdom.Attrs {
		if _, exists := newVdom.Attrs[key]; !exists {
			oldVdom.DOMNode.Call("removeAttribute", key)
		}
	}

	// Update Events
	for event, handler := range newVdom.Events {
		if oldHandler, exists := oldVdom.Events[event]; exists {
			oldVdom.DOMNode.Call("removeEventListener", event, js.FuncOf(oldHandler))
		}
		oldVdom.DOMNode.Call("addEventListener", event, js.FuncOf(handler))
	}
	for event := range oldVdom.Events {
		if _, exists := newVdom.Events[event]; !exists {
			oldVdom.DOMNode.Call("removeEventListener", event, js.FuncOf(oldVdom.Events[event]))
		}
	}

	// Update children
	maxLen := len(oldVdom.Children)
	if len(newVdom.Children) > maxLen {
		maxLen = len(newVdom.Children)
	}
	for i := 0; i < maxLen; i++ {
		if i >= len(oldVdom.Children) {
			// Add new child
			newChild := createDOM(newVdom.Children[i], doc)
			oldVdom.DOMNode.Call("appendChild", newChild.DOMNode)
			oldVdom.Children = append(oldVdom.Children, newChild)
		} else if i >= len(newVdom.Children) {
			// Remove excess child
			oldVdom.DOMNode.Call("removeChild", oldVdom.Children[i].DOMNode)
			oldVdom.Children = oldVdom.Children[:i]
			break
		} else {
			// Diff existing child
			diffAndPatch(oldVdom.Children[i], newVdom.Children[i], oldVdom.DOMNode, doc)
		}
	}
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
