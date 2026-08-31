package main

import (
	"github.com/enchant97/note-mark/backend/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BooksService struct{}
}

func (s BooksService) DeleteBookByID(currentUserID uuid.UUID, bookID uuid.UUID) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&db.Book{}).
			Where("id = ? AND owner_id = ?", bookID, currentUserID).
			Update("is_public", false).
			Error; err != nil {
			return dbErrorToServiceError(err)
		}
		result := tx.
			Where("id = ? AND owner_id = ?", bookID, currentUserID).
			Delete(&db.Book{})
		if err := result.Error; err != nil {
			return dbErrorToServiceError(err)
		}
		if result.RowsAffected == 0 {
			return NotFoundError
		}
		return nil
	})
}
