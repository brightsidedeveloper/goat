package goat

func RenderRoot(id string, comp Component, props any) {
	r := NewRenderer(id, comp, props)
	r.Render()
	select {}
}
