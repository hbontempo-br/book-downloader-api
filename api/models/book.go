package models

import (
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type BookModel struct {
	ID        uint      `gorm:"primary_key"`
	BookKey   uuid.UUID `gorm:"column:book_key;type:uuid;primary_key;"`
	Name      string    `gorm:"column:name"`
	Mask      string    `gorm:"column:mask"`
	StatusId  uint      `gorm:"column:status_id"`
	Status    BookStatusModel
	CreatedAt *time.Time
	DeletedAt *time.Time
}

func (BookModel) TableName() string {
	return "Book"
}

func (bm BookModel) String() string {
	return fmt.Sprintf("{ID:%v, BookKey:%v, Name:%v, Mask:%v, StatusId:%v, Status:%v, CreatedAt:%v, DeletedAt:%v}", bm.ID, bm.BookKey, bm.Name, bm.Mask, bm.StatusId, bm.Status.String(), bm.CreatedAt, bm.DeletedAt)
}
