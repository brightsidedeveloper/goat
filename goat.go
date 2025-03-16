package goat

import "syscall/js"

func Log(args ...any) {
	js.Global().Get("console").Call("log", args...)
}

func Alert(msg string) {
	js.Global().Call("alert", msg)
}

func EventCB(f func(el js.Value, event js.Value)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		f(this, args[0])
		return nil
	})
}
