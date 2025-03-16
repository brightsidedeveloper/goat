package goat

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// ComponentWrapper adapts a templ.Component to our Component interface.
type ComponentWrapper struct {
	renderFunc func(props any) templ.Component
}

func (t *ComponentWrapper) Render(ctx context.Context, w io.Writer, props any) error {
	comp := t.renderFunc(props)
	return comp.Render(ctx, w)
}

// WrapTemplComponent creates a Component from a function that returns a templ.Component.
func WrapComponent(renderFunc func(props any) templ.Component) Component {
	return &ComponentWrapper{renderFunc: renderFunc}
}
