package domain

import (
	"time"
)

// Example2 struct, models of example2 table
type Example2 struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at" pg:",soft_delete"`
}
