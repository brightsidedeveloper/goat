package goat

import (
	"context"
	"fmt"
	"sync"
	"syscall/js"
	"time"
)

type ComponentInstance struct {
	states     map[int]any
	callbacks  map[int]string
	stateOrder []int
	callIndex  int
	mu         sync.Mutex
}

type componentInstanceKeyType struct{}

var componentInstanceKey = componentInstanceKeyType{}

func getInstanceFromContext(ctx context.Context) *ComponentInstance {
	if ci, ok := ctx.Value(componentInstanceKey).(*ComponentInstance); ok {
		return ci
	}
	panic("No component instance found in context")
}

func UseState[T any](ctx context.Context, initialValue T) (func() T, func(T)) {
	ci := getInstanceFromContext(ctx)
	ci.mu.Lock()
	defer ci.mu.Unlock()

	if ci.callIndex >= len(ci.stateOrder) {
		stateKey := len(ci.stateOrder)
		ci.stateOrder = append(ci.stateOrder, stateKey)
		ci.states[stateKey] = initialValue
	}

	stateKey := ci.stateOrder[ci.callIndex]
	ci.callIndex++

	getState := func() T {
		ci.mu.Lock()
		defer ci.mu.Unlock()
		return ci.states[stateKey].(T)
	}

	setState := func(newValue T) {
		ci.mu.Lock()
		ci.states[stateKey] = newValue
		ci.mu.Unlock()
		if renderer := getRendererForInstance(ci); renderer != nil {
			go renderer.Render()
		}
	}

	return getState, setState
}

func UseCallback(ctx context.Context, f func(this js.Value, args []js.Value) any) func(js.Value, []js.Value) any {
	ci := getInstanceFromContext(ctx)
	ci.mu.Lock()
	defer ci.mu.Unlock()

	callbackIndex := ci.callIndex
	ci.callIndex++

	if oldName, exists := ci.callbacks[callbackIndex]; exists {
		js.Global().Delete(oldName)
	}

	name := fmt.Sprintf("fn%d_%d", time.Now().UnixNano(), callbackIndex)
	js.Global().Set(name, js.FuncOf(f))
	ci.callbacks[callbackIndex] = name

	return f
}

func UseEffect(ctx context.Context, effect func()) {
	ci := getInstanceFromContext(ctx)
	ci.mu.Lock()
	defer ci.mu.Unlock()

	// TODO: Queue this to fire after the re-render
	effect()
}

func (ci *ComponentInstance) Reset() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.callIndex = 0
}
