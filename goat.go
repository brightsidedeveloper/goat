package goat

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"syscall/js"

	"github.com/a-h/templ"
)

type TemplJoint interface {
	Render(context.Context, io.Writer) error
}

func HTML(j TemplJoint) string {
	return HTMLWithContext(j, context.Background())
}

func HTMLWithContext(j TemplJoint, c context.Context) string {
	var buf bytes.Buffer
	err := j.Render(c, &buf)
	if err != nil {
		js.Global().Get("console").Call("error", "Error rendering template:", err.Error())
		return ""
	}

	return buf.String()
}

func Log(args ...any) {
	js.Global().Get("console").Call("log", args...)
}

func RenderRootHTML(html string) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", html)
}

func RenderRoot(j TemplJoint) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", HTML(j))
}

func RenderRootWithContext(j TemplJoint, ctx context.Context) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", HTMLWithContext(j, ctx))
}

func RenderHTML(id string, html string) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", id)
	output.Set("outerHTML", html)
}

func Render(id string, j TemplJoint) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", id)
	output.Set("outerHTML", HTML(j))
}

func RenderWithContext(id string, j TemplJoint, ctx context.Context) {
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", id)
	output.Set("outerHTML", HTMLWithContext(j, ctx))
}

func jsFunc(name string, f func(this js.Value, args []js.Value) any) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return f(this, args)
	}))
}

func JSVar(name string, v any) {
	js.Global().Set(name, js.ValueOf(v))
}

func Callback(name string, f func(this js.Value, args []js.Value) any) func(args []js.Value) {
	jsFunc(name, f)
	return func(args []js.Value) {
		templ.JSFuncCall(name, args)
	}
}

type Driver struct {
	id         string
	states     map[string]any
	stateOrder []string
	mu         sync.Mutex
	callIndex  int
}

type DriverManager struct {
	drivers map[string]*Driver
	mu      sync.Mutex
}

var manager = &DriverManager{
	drivers: make(map[string]*Driver),
}

func UseDriver(driverID string) *Driver {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if drv, exists := manager.drivers[driverID]; exists {
		return drv
	}

	drv := &Driver{
		id:         driverID,
		states:     make(map[string]any),
		stateOrder: []string{},
	}
	manager.drivers[driverID] = drv
	return drv
}

func (d *Driver) State(initialValue any) (func() any, func(any)) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.callIndex >= len(d.stateOrder) {
		d.callIndex = 0
	}

	if d.callIndex >= len(d.stateOrder) {
		stateKey := generateStateKey(d.id, len(d.stateOrder))
		d.stateOrder = append(d.stateOrder, stateKey)
		d.states[stateKey] = initialValue
	}

	stateKey := d.stateOrder[d.callIndex]
	d.callIndex++

	getState := func() any {
		d.mu.Lock()
		defer d.mu.Unlock()
		return d.states[stateKey]
	}

	setState := func(newValue any) {
		d.mu.Lock()
		defer d.mu.Unlock()
		d.states[stateKey] = newValue
	}

	return getState, setState
}

func generateStateKey(driverID string, index int) string {
	return fmt.Sprintf("%s_state_%d", driverID, index)
}
