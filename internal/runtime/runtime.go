package runtime

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/buke/quickjs-go"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/katungi/edon/internal/errors"
	"github.com/katungi/edon/internal/modules/console"
)

type Runtime struct {
	jsRuntime *quickjs.Runtime
	context   *quickjs.Context
}

const (
	prompt      = ">>"
	multiPrompt = "..."
)

// REPL specific errors (re-exported from errors package)
var (
	ErrInterrupt = errors.ErrInterrupt
	ErrExit      = errors.ErrExit
)

func New() (*Runtime, error) {
	rt := quickjs.NewRuntime()
	ctx := rt.NewContext()

	r := &Runtime{
		jsRuntime: rt,
		context:   ctx,
	}

	// Initialize built-in modules
	if err := r.initializeBuiltins(); err != nil {
		ctx.Close()
		rt.Close()
		return nil, errors.WrapWith(errors.ErrBuiltinInit, err, "initialize builtins")
	}
	return r, nil
}

func (r *Runtime) initializeBuiltins() error {
	// Add console module
	if err := console.Init(r.context); err != nil {
		return errors.WrapWith(errors.ErrConsoleInit, err, "console module")
	}
	return nil
}

func (r *Runtime) Eval(script string) error {
	result := r.context.Eval(script)
	if result.IsException() {
		return fmt.Errorf("%s", result.String())
	}
	if !result.IsUndefined() {
		fmt.Println(result.String())
	}
	return nil
}

func (r *Runtime) ExecuteFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.WrapWith(errors.ErrFileRead, err, "")
	}

	result := r.context.Eval(string(data))
	if result.IsException() {
		return fmt.Errorf("%s", result.String())
	}
	if !result.IsUndefined() {
		fmt.Println(result.String())
	}
	return nil
}

func (r *Runtime) Close() {
	if r.context != nil {
		r.context.Close()
		r.context = nil
	}
	if r.jsRuntime != nil {
		r.jsRuntime.Close()
		r.jsRuntime = nil
	}
}

func (r *Runtime) StartREPL() error {
	fmt.Println("Starting REPL...")
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            prompt,
		HistoryFile:       "/tmp/edon_history",
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		HistoryLimit:      500,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
	})

	if err != nil {
		return errors.WrapWith(errors.ErrRuntimeInit, err, "create readline instance")
	}
	defer rl.Close()

	printWelcome()

	// The Fucking Loop
	multiline := false
	var code strings.Builder

	fmt.Println("Entering REPL loop...")
	for {
		currentPrompt := prompt
		if multiline {
			currentPrompt = multiPrompt
		}

		rl.SetPrompt(currentPrompt)

		line, err := rl.Readline()
		if err != nil {
			fmt.Printf("Readline error: %v\n", err)
			if err == readline.ErrInterrupt {
				if multiline {
					// Cancel multi-line input
					multiline = false
					code.Reset()
					continue
				}
				return ErrInterrupt
			}
			// #2: Avoid unnecessary else after return
			if err == io.EOF {
				fmt.Println("Exiting...")
				return nil
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle special commands
		if !multiline {
			switch line {
			case ".exit", "exit", "quit":
				fmt.Println("Exiting...")
				return nil
			case ".help", "help":
				printHelp()
				continue
			case ".clear":
				code.Reset()
				multiline = false
				continue
			}
		}

		// Append the line to our code buffer
		_, _ = code.WriteString(line)
		_, _ = code.WriteString("\n")

		// check if we need to continue reading more lines
		if isIncomplete(line) {
			multiline = true
			continue
		}

		// Execute the code
		fmt.Printf("Executing code: %s\n", code.String())
		result := r.context.Eval(code.String())
		if result.IsException() {
			color.Red("Error: %v", result.String())
		} else {
			if !result.IsUndefined() && !result.IsNull() {
				// Convert result to string and print
				str := result.String()
				if str != "undefined" && str != "" {
					color.Green("=> %s", str)
				}
			}
		}

		// Reset for next input
		code.Reset()
		multiline = false
	}
}

// isIncomplete checks if the input code block is incomplete and needs more lines
func isIncomplete(line string) bool {
	line = strings.TrimSpace(line)

	// Check for obvious continuation cases
	if strings.HasSuffix(line, "{") ||
		strings.HasSuffix(line, "\\") ||
		strings.HasSuffix(line, ".") {
		return true
	}

	// Count brackets/braces/parentheses
	brackets := 0
	braces := 0
	parens := 0

	for _, ch := range line {
		switch ch {
		case '[':
			brackets++
		case ']':
			brackets--
		case '{':
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		}
	}

	return brackets > 0 || braces > 0 || parens > 0
}

// printWelcome prints the REPL welcome message
func printWelcome() {
	color.Cyan("Welcome to Halo REPL!")
	color.Cyan("Type .help for more information")
	fmt.Println()
}

// printHelp prints the help information
func printHelp() {
	help := `
Commands:
  .help, help    Show this help message
  .exit, exit    Exit the REPL
  .clear         Clear the current input
  
Special Keys:
  Ctrl+C         Cancel current input / Exit REPL
  Ctrl+D         Exit REPL
  Up/Down        Navigate through history
  Tab            Auto-complete
`
	color.Yellow(help)
}
