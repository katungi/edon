package console

import (
	"fmt"
	"github.com/buke/quickjs-go"
)

func Init(ctx *quickjs.Context) error {
	global := ctx.Globals()
	console := ctx.Object()

	// Implement console.log
	log := ctx.Function(func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			str := arg.String()
			fmt.Print(str)
		}
		fmt.Println()
		return ctx.Undefined()
	})

	// Add the function to console object and console to global
	console.Set("log", log)
	global.Set("console", console)

	// Let QuickJS handle the cleanup of these values
	return nil
}
