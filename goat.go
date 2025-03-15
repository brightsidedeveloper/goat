package goat

import (
	"bytes"
	"context"
	"io"
	"syscall/js"
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

func JSFunc(name string, f func(this js.Value, args []js.Value) any) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return f(this, args)
	}))
}

func JSVar(name string, v any) {
	js.Global().Set(name, js.ValueOf(v))
}
