package blog

import (
	"time"
)

// Post represents a blog entry stored in the system.
type Post struct {
	ID        string
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
