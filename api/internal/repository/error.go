package repository

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrUnknown          = errors.New("unknown error")
	ErrDuplicateEntry   = errors.New("duplicate entry")
	ErrForeignKey       = errors.New("foreign key violation")
	ErrNotNull          = errors.New("not null violation")
	ErrConnectionFailed = errors.New("database connection failed")
)

func MapDBError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return fmt.Errorf("failed to insert record: %w", ErrDuplicateEntry)
		case "23503":
			return fmt.Errorf("failed to insert record: %w", ErrForeignKey)
		case "23502":
			return fmt.Errorf("failed to insert record: %w", ErrNotNull)
		case "08000", "08003", "08006":
			return fmt.Errorf("database operation failed: %w", ErrConnectionFailed)
		default:
			return fmt.Errorf("database operation failed: %w", ErrUnknown)
		}
	}
	return fmt.Errorf("database operation failed: %w", ErrUnknown)
}
