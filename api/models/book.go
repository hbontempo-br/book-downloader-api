package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type BookModel struct {
	ID        uint            `gorm:"primary_key"`
	BookKey   uuid.UUID       `gorm:"column:book_key;type:uuid;primary_key;"`
	Name      string          `gorm:"column:name"`
	Mask      string          `gorm:"column:mask"`
	StatusID  uint            `gorm:"column:status_id"`
	Status    BookStatusModel `gorm:"association_autoupdate:false"`
	CreatedAt *time.Time
	DeletedAt *time.Time
}

func (BookModel) TableName() string {
	return "Book"
}

func (bm BookModel) String() string {
	return fmt.Sprintf("{ID:%v, BookKey:%v, Name:%v, Mask:%v, StatusID:%v, Status:%v, CreatedAt:%v, DeletedAt:%v}", bm.ID, bm.BookKey, bm.Name, bm.Mask, bm.StatusID, bm.Status.String(), bm.CreatedAt, bm.DeletedAt)
}
