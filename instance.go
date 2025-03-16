package goat

import (
	"context"
	"sync"
)

type ComponentInstance struct {
	states     map[int]any
	stateOrder []int
	callIndex  int
	mu         sync.Mutex
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

func (ci *ComponentInstance) UseState(initialValue any) (func() any, func(any)) {
	ci.mu.Lock()
	defer ci.mu.Unlock()

	if ci.callIndex >= len(ci.stateOrder) {
		stateKey := len(ci.stateOrder)
		ci.stateOrder = append(ci.stateOrder, stateKey)
		ci.states[stateKey] = initialValue
	}

	stateKey := ci.stateOrder[ci.callIndex]
	ci.callIndex++

	getState := func() any {
		ci.mu.Lock()
		defer ci.mu.Unlock()
		return ci.states[stateKey]
	}

	setState := func(newValue any) {
		ci.mu.Lock()
		ci.states[stateKey] = newValue
		ci.mu.Unlock()
		// Trigger re-render (placeholder for now)
		renderer := getRendererForInstance(ci)
		if renderer != nil {
			go renderer.Render()
		}
	}

	return getState, setState
}

// Temporary placeholder; to be implemented in Step 3
func getRendererForInstance(ci *ComponentInstance) *Renderer {
	return nil
}
