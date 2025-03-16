package goat

import "syscall/js"

type VNode struct {
	Tag      string
	Attrs    map[string]string
	Events   map[string]func(js.Value, []js.Value) any // Event handlers
	Children []VNode
	Text     string
	DOMNode  js.Value
}

func NewVNode(tag string, attrs map[string]string, events map[string]func(js.Value, []js.Value) any, children []VNode, text string) VNode {
	return VNode{
		Tag:      tag,
		Attrs:    attrs,
		Events:   events,
		Children: children,
		Text:     text,
	}
}
