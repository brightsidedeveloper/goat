package goat

import (
	"bytes"
	"context"
	"sync"
)

var (
	renderersMu sync.Mutex
	renderers   = make(map[string]*Renderer)
)

func RenderRoot(id string, comp Component, props any) {
	renderersMu.Lock()
	defer renderersMu.Unlock()
	r := NewRenderer(id, comp, props)
	renderers[id] = r
	r.Render()
}

func html(j Component, ctx context.Context, props any) string {
	vdom := j.Render(ctx, props)
	var buf bytes.Buffer
	buf.WriteString(vdomToHTML(vdom))
	return buf.String()
}
