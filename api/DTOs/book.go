package DTOs

import (
	"github.com/hbontempo-br/book-downloader-api/api/models"
)

type BookDTO struct {
	BookKey string        `json:"book_key"`
	Name    string        `json:"name"`
	Mask    string        `json:"mask"`
	Status  BookStatusDTO `json:"status"`
}

func NewBookDTO(model models.BookModel) BookDTO {
	return BookDTO{BookKey: model.BookKey.String(), Name: model.Name, Mask: model.Mask, Status: NewBookStatusDTO(model.Status)}
}
