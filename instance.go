package goat

import (
	"context"
	"fmt"
	"sync"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
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
		if renderer := getRendererForInstance(ci); renderer != nil {
			go renderer.Render() // Non-blocking re-render
		}
	}

	return getState, setState
}

func (ci *ComponentInstance) UseCallback(f func(this js.Value, args []js.Value) any) func(args ...any) templ.ComponentScript {
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

	return func(args ...any) templ.ComponentScript {
		jsArgs := make([]js.Value, len(args))
		for i, arg := range args {
			jsArgs[i] = js.ValueOf(arg)
		}
		return templ.JSFuncCall(name, jsArgs)
	}
}

func (ci *ComponentInstance) Reset() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.callIndex = 0
}
