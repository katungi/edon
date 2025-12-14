package unit

import (
	"testing"

	"github.com/katungi/edon/internal/runtime"
)

func TestBasicExecution(t *testing.T) {
	test := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name:    "simple console log",
			script:  `console.log("test")`,
			wantErr: false,
		},
		{
			name:    "syntax error",
			script:  `console.log("test"`,
			wantErr: true,
		},
		{
			name:    "basic math",
			script:  `2 + 2 === 4`,
			wantErr: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := runtime.New()
			if err != nil {
				t.Fatalf("Failed to create runtime: %v", err)
			}
			defer rt.Close()

			err = rt.Eval(tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

}
