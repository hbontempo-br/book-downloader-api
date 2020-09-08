package models

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBookModel_TableName(t *testing.T) {
	t.Run("TableName method", func(t *testing.T) {
		expectedTableName := "Book"
		bookModel := BookModel{}
		assert.Equal(t, expectedTableName, bookModel.TableName(), "Unexpected TableName method return")
	})
}

func TestBook_String(t *testing.T) {
	t.Run("String method", func(t *testing.T) {

		testUUID, err := uuid.NewV4()
		if err != nil {
			t.Error("Test failed on setup, unable to create random UUID")
		}
		testTime := time.Now()
		bookModel := BookModel{
			ID: 1,
			BookKey: testUUID,
			Name: "name",
			Mask: "mask",
			StatusID: 1,
			Status: BookStatusModel{1, "status"},
			CreatedAt: &testTime,
			DeletedAt: &testTime,
		}

		expectedTableName := fmt.Sprintf("{ID:1, BookKey:%v, Name:name, Mask:mask, StatusID:1, Status:{ID:1, Enumerator:status}, CreatedAt:%v, DeletedAt:%v}", testUUID, testTime, testTime)

		assert.Equal(t, expectedTableName, bookModel.String(), "Unexpected String method return")
	})
}
