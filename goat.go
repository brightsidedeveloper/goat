package goat

import "syscall/js"

func Log(args ...any) {
	js.Global().Get("console").Call("log", args...)
}
