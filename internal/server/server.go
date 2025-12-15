package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/katungi/edon/internal/runtime"
)

type EvalRequest struct {
	Code string `json:"code"`
}

type EvalResponse struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Server struct {
	rt     *runtime.Runtime
	port   string
	server *http.Server
	evalMu sync.Mutex // Serializes eval requests to avoid stdout race
}

func New() (*Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rt, err := runtime.New()
	if err != nil {
		return nil, err
	}

	return &Server{
		rt:   rt,
		port: port,
	}, nil
}

func (s *Server) Close() {
	if s.server != nil {
		_ = s.server.Close()
	}
	if s.rt != nil {
		s.rt.Close()
	}
}

func (s *Server) Start(staticDir string) error {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}
		http.NotFound(w, r)
	})

	// Handle REPL evaluation
	mux.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req EvalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Serialize eval requests to avoid stdout race condition
		s.evalMu.Lock()
		defer s.evalMu.Unlock()

		// Capture stdout
		oldStdout := os.Stdout
		pipeReader, pipeWriter, err := os.Pipe()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		os.Stdout = pipeWriter

		// Evaluate the code
		evalErr := s.rt.Eval(req.Code)

		// Read the output
		_ = pipeWriter.Close()
		var buf strings.Builder
		_, _ = io.Copy(&buf, pipeReader)
		os.Stdout = oldStdout
		_ = pipeReader.Close()

		response := EvalResponse{}
		if evalErr != nil {
			response.Error = evalErr.Error()
		} else {
			output := buf.String()
			if output == "" {
				output = "=> " + req.Code
			}
			response.Output = output
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})

	// #81: Don't use default HTTP server - configure timeouts
	s.server = &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on http://localhost:%s\n", s.port)
	return s.server.ListenAndServe()
}
