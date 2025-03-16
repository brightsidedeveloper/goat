package goat

import (
	"bytes"
	"context"
	"syscall/js"
)

func RenderRoot(j Component, props any) {
	ci := &ComponentInstance{}
	ctx := context.WithValue(context.Background(), componentInstanceKey, ci)
	doc := js.Global().Get("document")
	output := doc.Call("getElementById", "root")
	output.Set("innerHTML", html(j, ctx, props))
}

func html(j Component, ctx context.Context, props any) string {
	var buf bytes.Buffer
	err := j.Render(ctx, &buf, props)
	if err != nil {
		js.Global().Get("console").Call("error", "Error rendering template:", err.Error())
		return ""
	}
	return buf.String()
}
