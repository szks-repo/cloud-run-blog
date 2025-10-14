package server

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/szks-repo/cloud-run-blog/internal/blog"
	templatefs "github.com/szks-repo/cloud-run-blog/web/templates"
)

// Server wraps the HTTP mux and dependencies for the blog API/UI.
type Server struct {
	repo      blog.Repository
	mux       *http.ServeMux
	templates *template.Template
	logger    *log.Logger
}

// New constructs a Server instance with routes wired up.
func New(repo blog.Repository) (*Server, error) {
	tpl, err := template.New("base").ParseFS(templatefs.FS, "*.tmpl")
	if err != nil {
		return nil, err
	}

	s := &Server{
		repo:      repo,
		mux:       http.NewServeMux(),
		templates: tpl,
		logger:    log.Default(),
	}

	s.registerRoutes()

	return s, nil
}

// Run starts the HTTP server until the provided context is cancelled.
func (s *Server) Run(ctx context.Context, addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return ctx.Err()
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// ServeHTTP delegates to the internal mux.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux.HandleFunc("/", s.handleIndex)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := s.repo.List(r.Context())
	if err != nil {
		s.logger.Printf("list posts: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts []blog.Post
	}{
		Posts: posts,
	}

	if err := s.templates.ExecuteTemplate(w, "index.tmpl", data); err != nil {
		s.logger.Printf("render index: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
