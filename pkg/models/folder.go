package models

import (
	"time"
)

type Folder struct {
	ID         uint64     `json:"id" db:"id"`
	UserID     uint64     `json:"user_id" db:"user_id"`
	ParentID   *uint64    `json:"parent_id,omitempty" db:"parent_id"`
	Name       string     `json:"name" db:"name"`
	Path       string     `json:"path" db:"path"`
	IsPublic   bool       `json:"is_public" db:"is_public"`
	Color      string     `json:"color,omitempty" db:"color"`
	Icon       string     `json:"icon,omitempty" db:"icon"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at,omitempty" db:"modified_at"`
}

type Folders []*Folder

func (f *Folder) Rename(newName string, currentTime time.Time) {
	f.Name = newName
	f.ModifiedAt = &currentTime
}

func (f *Folder) IsRoot() bool {
	return f.ParentID == nil
}
