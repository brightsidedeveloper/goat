package goat

import "syscall/js"

type GoatNode struct {
	Tag      string
	Attrs    map[string]string
	Events   map[string]js.Func
	Children []GoatNode
	Text     string
	DOMNode  js.Value
}
