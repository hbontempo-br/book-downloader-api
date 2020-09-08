package DTOs

import (
	"github.com/gofrs/uuid"
	"github.com/hbontempo-br/book-downloader-api/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBookDTO(t *testing.T) {
	t.Run("BookDTO generation", func(t *testing.T) {

		testUUID, err := uuid.NewV4()
		if err != nil {
			t.Error("Test failed on setup, unable to create random UUID")
		}
		testTime := time.Now()
		bookModel := models.BookModel{
			ID:        1,
			BookKey:   testUUID,
			Name:      "name",
			Mask:      "mask",
			StatusID:  1,
			Status:    models.BookStatusModel{1, "status"},
			CreatedAt: &testTime,
			DeletedAt: &testTime,
		}
		bookDTO := NewBookDTO(bookModel)
		expectedBookDTO := BookDTO{
			Name:    "name",
			Mask:    "mask",
			BookKey: testUUID.String(),
			Status:  "status",
		}

		assert.Equal(t, expectedBookDTO, bookDTO, "Unexpected BookDTO")

	})
}
