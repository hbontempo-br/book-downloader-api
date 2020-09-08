package DTOs

import (
	"github.com/hbontempo-br/book-downloader-api/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBookStatusDTO(t *testing.T) {
	t.Run("BookStatusDTO generation", func(t *testing.T) {
		testStatus := "status"
		bookStatusModel := models.BookStatusModel{ID: 1, Enumerator: testStatus}
		bookStatusDTO := NewBookStatusDTO(bookStatusModel)
		expectedBookStatusDTO := BookStatusDTO(testStatus)

		assert.Equal(t, expectedBookStatusDTO, bookStatusDTO, "Unexpected BookStatusDTO")
	})
}
