package goat

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
)

func Callback(f func(this js.Value, args []js.Value) any) func(args ...any) templ.ComponentScript {
	name := fmt.Sprintf("fn%d", time.Now().UnixNano()) // Unique function name
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return f(this, args)
	}))
	return func(args ...any) templ.ComponentScript {
		jsArgs := make([]js.Value, len(args))
		for i, arg := range args {
			jsArgs[i] = js.ValueOf(arg)
		}
		return templ.JSFuncCall(name, jsArgs)
	}
}
