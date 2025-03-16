package goat

import "context"

type ComponentInstance struct {
}

type componentInstanceKeyType struct{}

var componentInstanceKey = componentInstanceKeyType{}

// GetComponentInstance retrieves the ComponentInstance from the context.
func GetComponentInstance(ctx context.Context) *ComponentInstance {
	if ci, ok := ctx.Value(componentInstanceKey).(*ComponentInstance); ok {
		return ci
	}
	panic("No component instance found in context")
}
