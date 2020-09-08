package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookStatusModel_TableName(t *testing.T) {
	t.Run("TableName method", func(t *testing.T) {
		expectedTableName := "BookStatus"
		bookStatusModel := BookStatusModel{}
		assert.Equal(t, expectedTableName, bookStatusModel.TableName(), "Unexpected TableName method return")
	})
}

func TestBookModel_String(t *testing.T) {
	t.Run("String method", func(t *testing.T) {
		expectedTableName := "{ID:1, Enumerator:a}"
		bookStatusModel := BookStatusModel{ID: 1, Enumerator: "a"}
		assert.Equal(t, expectedTableName, bookStatusModel.String(), "Unexpected String method return")
	})
}
