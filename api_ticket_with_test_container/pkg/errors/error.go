package errors

import "fmt"

const (
	ErrFailedToList   = "failed to list %s"
	ErrFailedToGet    = "failed to get %s"
	ErrFailedToCreate = "failed to create %s"
	ErrFailedToUpdate = "failed to update %s"
	ErrFailedToDelete = "failed to delete %s"
	ErrNotFound       = "%s not found"
	ErrUnauthorized   = "unauthorized access to %s"
	ErrInvalidInput   = "invalid input for %s"
	ErrConflict       = "conflict detected in %s"
	ErrInternal       = "internal error with %s"
)

func FailedToList(entity string) string {
	return fmt.Sprintf(ErrFailedToList, entity)
}
func FailedToGet(entity string) string {
	return fmt.Sprintf(ErrFailedToGet, entity)
}

func FailedToCreate(entity string) string {
	return fmt.Sprintf(ErrFailedToCreate, entity)
}

func FailedToUpdate(entity string) string {
	return fmt.Sprintf(ErrFailedToUpdate, entity)
}

func FailedToDelete(entity string) string {
	return fmt.Sprintf(ErrFailedToDelete, entity)
}

func NotFound(entity string) string {
	return fmt.Sprintf(ErrNotFound, entity)
}

func Unauthorized(entity string) string {
	return fmt.Sprintf(ErrUnauthorized, entity)
}

func InvalidInput(entity string) string {
	return fmt.Sprintf(ErrInvalidInput, entity)
}

func Conflict(entity string) string {
	return fmt.Sprintf(ErrConflict, entity)
}

func Internal(entity string) string {
	return fmt.Sprintf(ErrInternal, entity)
}
