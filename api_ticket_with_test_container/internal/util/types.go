package util

import (
	"database/sql"
	"time"
)

func NullInt32ToPtr(n sql.NullInt32) *int32 {
	if n.Valid {
		return &n.Int32
	}
	return nil
}

func NullStringToPtr(n sql.NullString) *string {
	if n.Valid {
		return &n.String
	}
	return nil
}

func NullDateToPtr(n sql.NullTime) *time.Time {
	if n.Valid {
		return &n.Time
	}
	return nil
}
