package services

import (
	"errors"
	"fmt"

	"github.com/lapotkin/file-storage/internal/adapter/db"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserDeactivated    = errors.New("user is deactivated")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
	ErrInternal           = errors.New("internal server error")
)

func handleDBError(err error) error {
	if dbErr := pg.HandleError(err); dbErr != nil {
		var uniqueErr *db.UniqueViolationError
		if errors.As(dbErr, &uniqueErr) {
			return NewDuplicateUniqueValueError(uniqueErr)
		}
		var fkErr *db.ViolatedForeignKeyError
		if errors.As(dbErr, &fkErr) {
			return NewViolatedForeignKeyValueError(fkErr)
		}
		if errors.Is(dbErr, db.ErrNoRows) {
			return ErrNotFound
		}
	}

	return ErrInternal
}

type DuplicateUniqueValueError struct {
	Table   string `json:"table"`
	Details string `json:"details"`
}

func NewDuplicateUniqueValueError(err *db.UniqueViolationError) *DuplicateUniqueValueError {
	return &DuplicateUniqueValueError{
		Table:   err.TableName,
		Details: err.Details,
	}
}

type ViolatedForeignKeyValueError struct {
	Table   string `json:"table"`
	Details string `json:"details"`
}

func NewViolatedForeignKeyValueError(err *db.ViolatedForeignKeyError) *ViolatedForeignKeyValueError {
	return &ViolatedForeignKeyValueError{
		Table:   err.TableName,
		Details: err.Details,
	}
}

func (e *DuplicateUniqueValueError) Error() string {
	return fmt.Sprintf("Duplicate unique value error in table %s: %s", e.Table, e.Details)
}

func (e *ViolatedForeignKeyValueError) Error() string {
	return fmt.Sprintf("Foreign key violation value error in table %s: %s", e.Table, e.Details)
}
