package blog

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// InMemoryRepository stores posts in-memory. Suitable for local development and tests.
type InMemoryRepository struct {
	mu    sync.RWMutex
	posts map[string]Post
}

// NewInMemoryRepository creates a repository backed by an in-memory map.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts: make(map[string]Post),
	}
}

// List returns all posts ordered by creation time (newest first).
func (r *InMemoryRepository) List(ctx context.Context) ([]Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]Post, 0, len(r.posts))
	for _, post := range r.posts {
		items = append(items, post)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	return items, nil
}

// Get returns a post by ID.
func (r *InMemoryRepository) Get(ctx context.Context, id string) (Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return Post{}, ErrNotFound
	}

	return post, nil
}

// Create stores a new post, assigning an ID if needed.
func (r *InMemoryRepository) Create(ctx context.Context, post Post) (Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	if post.ID == "" {
		post.ID = fmt.Sprintf("%d", now.UnixNano())
	}
	post.CreatedAt = now
	post.UpdatedAt = now

	r.posts[post.ID] = post
	return post, nil
}

// Update replaces an existing post.
func (r *InMemoryRepository) Update(ctx context.Context, post Post) (Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.posts[post.ID]
	if !ok {
		return Post{}, ErrNotFound
	}

	post.CreatedAt = current.CreatedAt
	post.UpdatedAt = time.Now().UTC()

	r.posts[post.ID] = post
	return post, nil
}

// Delete removes a post by ID.
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[id]; !ok {
		return ErrNotFound
	}

	delete(r.posts, id)
	return nil
}
