package goat

import "syscall/js"

type VNode struct {
	Tag      string
	Attrs    map[string]string
	Children []VNode
	Text     string
	DOMNode  js.Value
	// TODO: Need event listeners and complex props
}

func NewVNode(tag string, attrs map[string]string, children []VNode, text string) VNode {
	return VNode{
		Tag:      tag,
		Attrs:    attrs,
		Children: children,
		Text:     text,
	}
}
