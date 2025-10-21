package models

import (
	"time"
)

type Share struct {
	ID          uint64     `json:"id" db:"id"`
	DocumentID  uint64     `json:"document_id" db:"document_id"`
	OwnerID     uint64     `json:"owner_id" db:"owner_id"`
	SharedWith  *uint64    `json:"shared_with,omitempty" db:"shared_with"`
	ShareLink   string     `json:"share_link" db:"share_link"`
	Permission  Permission `json:"permission" db:"permission"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	AccessCount int        `json:"access_count" db:"access_count"`
	MaxAccess   int        `json:"max_access" db:"max_access"`
	Password    string     `json:"-" db:"password"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty" db:"modified_at"`
}

type Shares []*Share
type Permission string

const (
	PermissionView     Permission = "view"
	PermissionDownload Permission = "download"
	PermissionEdit     Permission = "edit"
)

func (s *Share) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

func (s *Share) IsAccessLimitReached() bool {
	if s.MaxAccess < 0 {
		return false
	}
	return s.AccessCount >= s.MaxAccess
}

func (s *Share) IsValid() bool {
	return !s.IsExpired() && !s.IsAccessLimitReached()
}

func (s *Share) CanView() bool {
	return s.Permission == PermissionView || s.Permission == PermissionDownload || s.Permission == PermissionEdit
}

func (s *Share) CanDownload() bool {
	return s.Permission == PermissionDownload || s.Permission == PermissionEdit
}

func (s *Share) CanEdit() bool {
	return s.Permission == PermissionEdit
}
