package runtime

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/buke/quickjs-go"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/katungi/edon/internals/modules/console"
)

type Runtime struct {
	jsRuntime *quickjs.Runtime
	context   *quickjs.Context
}

const (
	prompt      = ">> (ﾉ˚Д˚)ﾉ ⌒┻━┻"
	multiPrompt = "..."
)

// REPL specific errors
var (
	ErrInterrupt = fmt.Errorf("interrupted")
	ErrExit      = fmt.Errorf("exit")
)

func New() (*Runtime, error) {
	fmt.Println("Creating new runtime...")
	rt := quickjs.NewRuntime()
	fmt.Println("Creating new context...")
	ctx := rt.NewContext()

	r := &Runtime{
		jsRuntime: &rt,
		context:   ctx,
	}

	// Initialize built-in modules
	fmt.Println("Initializing built-in modules...")
	if err := r.initializeBuiltins(); err != nil {
		fmt.Printf("Failed to initialize builtins: %v\n", err)
		ctx.Close()
		rt.Close()
		return nil, fmt.Errorf("failed to initialize builtins: %w", err)
	}
	fmt.Println("Runtime initialization complete")

	return r, nil
}

func (r *Runtime) initializeBuiltins() error {
	// Add console module
	if err := console.Init(r.context); err != nil {
		return fmt.Errorf("failed to initialize console: %w", err)
	}
	return nil
}

func (r *Runtime) Eval(script string) error {
	result, err := r.context.Eval(script)
	if err != nil {
		return err
	}
	if !result.IsUndefined() {
		fmt.Println(result.String())
	}
	return nil
}

func (r *Runtime) ExecuteFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	result, err := r.context.Eval(string(data))
	if err != nil {
		return err
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
		return fmt.Errorf("failed to create readline instance: %w", err)
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
			} else if err == io.EOF {
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
		code.WriteString(line)
		code.WriteString("\n")

		// check if we need to continue reading more lines
		if isIncomplete(line) {
			multiline = true
			continue
		}

		// Execute the code
		fmt.Printf("Executing code: %s\n", code.String())
		result, err := r.context.Eval(code.String())
		if err != nil {
			color.Red("Error: %v", err)
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


// createCompleter creates an autocomplete handler
func createCompleter() *readline.PrefixCompleter {
    return readline.NewPrefixCompleter(
        readline.PcItem("console.log"),
        readline.PcItem("let"),
        readline.PcItem("const"),
        readline.PcItem("function"),
        readline.PcItem("return"),
        readline.PcItem("if"),
        readline.PcItem("else"),
        readline.PcItem("for"),
        readline.PcItem("while"),
        readline.PcItem(".help"),
        readline.PcItem(".exit"),
        readline.PcItem(".clear"),
    )
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