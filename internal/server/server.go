package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

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

type indexData struct {
	Posts []blog.Post
}

type formData struct {
	Post   blog.Post
	Action string
	IsEdit bool
	Error  string
}

// New constructs a Server instance with routes wired up.
func New(repo blog.Repository) (*Server, error) {
	tpl, err := template.New("base").Funcs(template.FuncMap{
		"formatTime": formatTime,
	}).ParseFS(templatefs.FS, "*.tmpl")
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
	s.mux.HandleFunc("/posts", s.handlePostCreate)
	s.mux.HandleFunc("/posts/new", s.handlePostNew)
	s.mux.HandleFunc("/posts/", s.handlePostByID)
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

	data := indexData{
		Posts: posts,
	}

	if err := s.renderTemplate(w, "index.tmpl", data); err != nil {
		s.logger.Printf("render index: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handlePostNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := formData{
		Post:   blog.Post{},
		Action: "/posts",
		IsEdit: false,
	}

	if err := s.renderTemplate(w, "form.tmpl", data); err != nil {
		s.logger.Printf("render new form: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handlePostCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	post := blog.Post{
		Title: strings.TrimSpace(r.FormValue("title")),
		Body:  strings.TrimSpace(r.FormValue("body")),
	}

	if post.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		data := formData{
			Post:   post,
			Action: "/posts",
			IsEdit: false,
			Error:  "title is required",
		}
		if err := s.renderTemplate(w, "form.tmpl", data); err != nil {
			s.logger.Printf("render new form with error: %v", err)
		}
		return
	}

	saved, err := s.repo.Create(r.Context(), post)
	if err != nil {
		s.logger.Printf("create post: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.logger.Printf("post created: %s", saved.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handlePostByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	parts := strings.Split(path, "/")
	id := parts[0]
	if id == "" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 1 {
		http.NotFound(w, r)
		return
	}

	switch parts[1] {
	case "edit":
		if r.Method == http.MethodGet {
			s.handlePostEditForm(w, r, id)
			return
		}
		if r.Method == http.MethodPost {
			s.handlePostUpdate(w, r, id)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	case "delete":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handlePostDelete(w, r, id)
	case "markdown":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handlePostMarkdown(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handlePostEditForm(w http.ResponseWriter, r *http.Request, id string) {
	post, err := s.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, blog.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		s.logger.Printf("get post for edit: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := formData{
		Post:   post,
		Action: fmt.Sprintf("/posts/%s/edit", id),
		IsEdit: true,
	}

	if err := s.renderTemplate(w, "form.tmpl", data); err != nil {
		s.logger.Printf("render edit form: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) handlePostUpdate(w http.ResponseWriter, r *http.Request, id string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	post := blog.Post{
		ID:    id,
		Title: strings.TrimSpace(r.FormValue("title")),
		Body:  strings.TrimSpace(r.FormValue("body")),
	}

	if post.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		data := formData{
			Post:   post,
			Action: fmt.Sprintf("/posts/%s/edit", id),
			IsEdit: true,
			Error:  "title is required",
		}
		if err := s.renderTemplate(w, "form.tmpl", data); err != nil {
			s.logger.Printf("render edit form with error: %v", err)
		}
		return
	}

	updated, err := s.repo.Update(r.Context(), post)
	if err != nil {
		if errors.Is(err, blog.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		s.logger.Printf("update post: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.logger.Printf("post updated: %s", updated.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handlePostDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, blog.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		s.logger.Printf("delete post: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.logger.Printf("post deleted: %s", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handlePostMarkdown(w http.ResponseWriter, r *http.Request, id string) {
	post, err := s.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, blog.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		s.logger.Printf("get post for markdown: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s.md", slugify(post.Title))
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

	content := fmt.Sprintf("# %s\n\n%s\n", post.Title, post.Body)
	if _, err := w.Write([]byte(content)); err != nil {
		s.logger.Printf("write markdown response: %v", err)
	}
}

func (s *Server) renderTemplate(w http.ResponseWriter, name string, data any) error {
	return s.templates.ExecuteTemplate(w, name, data)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("2006-01-02 15:04 MST")
}

func slugify(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return "post"
	}

	var normalized strings.Builder
	lastWasHyphen := false
	for _, r := range input {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			normalized.WriteRune(r)
			lastWasHyphen = false
		case r == ' ' || r == '-' || r == '_':
			if !lastWasHyphen {
				normalized.WriteRune('-')
				lastWasHyphen = true
			}
		}
	}

	result := strings.Trim(normalized.String(), "-")
	if result == "" {
		return "post"
	}
	return result
}
