package model

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrBadIDFormat   = errors.New("bad_id_format")
	ErrNotFound      = errors.New("not_found")

	// Errors map error codes to error messages
	Errors = map[string]string{
		"bad_id_format": "bad ID format",
		"not_found":     "post not found",
	}
)

type Post struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt uint32 `json:"created_at"`
	User      User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
