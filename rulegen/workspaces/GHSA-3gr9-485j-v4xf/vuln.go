package main

import (
	"github.com/enchant97/note-mark/backend/db"
	"github.com/google/uuid"
)

type BooksService struct{}
}

func (s BooksService) DeleteBookByID(currentUserID uuid.UUID, bookID uuid.UUID) error {
	result := db.DB.
		Where("id = ? AND owner_id = ?", bookID, currentUserID).
		Delete(&db.Book{})
	if err := result.Error; err != nil {
		return dbErrorToServiceError(err)
	}
	if result.RowsAffected == 0 {
		return NotFoundError
	}
	return nil
}
