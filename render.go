package goat

import (
	"bytes"
	"context"
	"syscall/js"
)

func RenderRoot(r *Renderer) {
	r.Render()
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
