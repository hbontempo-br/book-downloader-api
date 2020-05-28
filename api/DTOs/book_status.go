package DTOs

import "github.com/hbontempo-br/book-downloader-api/api/models"

type BookStatusDTO string

func NewBookStatusDTO(model models.BookStatusModel) BookStatusDTO {
	return BookStatusDTO(model.Enumerator)
}
