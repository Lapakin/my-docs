package models

import (
	"time"
)

type Document struct {
	ID          uint64     `json:"id" db:"id"`
	UserID      uint64     `json:"user_id" db:"user_id"`
	FolderID    *uint64    `json:"folder_id,omitempty" db:"folder_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	FilePath    string     `json:"file_path" db:"file_path"`
	FileSize    int64      `json:"file_size" db:"file_size"`
	MimeType    string     `json:"mime_type" db:"mime_type"`
	IsPublic    bool       `json:"is_public" db:"is_public"`
	Tags        []string   `json:"tags,omitempty" db:"-"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty" db:"modified_at"`
}

type Documents []*Document

func (d *Document) Rename(newName string, currentTime time.Time) {
	d.Name = newName
	d.ModifiedAt = &currentTime
}

func (d *Document) MakePublic(currentTime time.Time) {
	d.IsPublic = true
	d.ModifiedAt = &currentTime
}

func (d *Document) MakePrivate(currentTime time.Time) {
	d.IsPublic = false
	d.ModifiedAt = &currentTime
}
