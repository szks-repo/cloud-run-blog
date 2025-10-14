package blog

import (
	"context"
	"errors"
)

// ErrNotFound indicates that a post does not exist in the repository.
var ErrNotFound = errors.New("post not found")

// Repository defines how posts are persisted.
type Repository interface {
	List(ctx context.Context) ([]Post, error)
	Get(ctx context.Context, id string) (Post, error)
	Create(ctx context.Context, post Post) (Post, error)
	Update(ctx context.Context, post Post) (Post, error)
	Delete(ctx context.Context, id string) error
}
