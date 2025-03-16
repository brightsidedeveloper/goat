package goat

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

// ComponentWrapper adapts a templ.Component to our Component interface.
type ComponentWrapper struct {
	renderFunc func(props any) templ.Component
}

func (t *ComponentWrapper) Render(ctx context.Context, props any) VNode {
	comp := t.renderFunc(props)
	// Convert templ.Component to VNode (temporary hack until templ supports VNode directly)
	return templToVNode(comp, ctx)
}

// WrapTemplComponent creates a Component from a function that returns a templ.Component.
func WrapComponent(renderFunc func(props any) templ.Component) Component {
	return &ComponentWrapper{renderFunc: renderFunc}
}

func templToVNode(comp templ.Component, ctx context.Context) VNode {
	// For now, render to string and parse (inefficient, but works)
	var buf bytes.Buffer
	comp.Render(ctx, &buf)
	html := buf.String()
	return VNode{Tag: "div", Text: html} // Simplified; assumes a div wrapper
}
