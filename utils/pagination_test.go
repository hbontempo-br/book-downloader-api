package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatPaginatedResponse(t *testing.T) {

	testData := []interface{}{
		map[string]string{
			"test": "case",
		},
		map[string]string{
			"test2": "case2",
		},
	}

	t.Run("Invalid pageSize (pageSize = 0)", func(t *testing.T) {
		resp, err := FormatPaginatedResponse(testData, 0, 1, 99)
		assert.Equal(t, InvalidPageSizeError, err, "Expect InvalidPageSizeError error")
		assert.Nil(t, resp, "Result should be nil on error")
	})

	t.Run("Invalid page (pageSize = 0)", func(t *testing.T) {
		resp, err := FormatPaginatedResponse(testData, 1, 0, 99)
		assert.Equal(t, InvalidPageError, err, "Expect InvalidPageError error")
		assert.Nil(t, resp, "Result should be nil on error")
	})

	t.Run("Invalid totalCount (totalCount = -1)", func(t *testing.T) {
		resp, err := FormatPaginatedResponse(testData, 1, 1, -1)
		assert.Equal(t, InvalidTotalCountError, err, "Expect InvalidTotalCountError error")
		assert.Nil(t, resp, "Result should be nil on error")
	})

	t.Run("PageSize < len(data)", func(t *testing.T) {
		resp, err := FormatPaginatedResponse(testData, 1, 1, 1)
		assert.Equal(t, PageSizeDataMismatchError, err, "Expect PageSizeDataMismatchError error")
		assert.Nil(t, resp, "Result should be nil on error")
	})

	t.Run("Valid case (No nextPage or previousPage", func(t *testing.T) {
		pageSize := 1
		page := 1
		totalCount := 0
		expectedResp := PaginatedResponse{
			Data: []interface{}{},
			Pagination: Pagination{
				CurrentPage:  page,
				NextPage:     nil,
				PreviousPage: nil,
				MaxPage:      0,
				RowsPerPage:  pageSize,
				TotalRows:    totalCount,
			},
		}
		resp, err := FormatPaginatedResponse([]interface{}{}, pageSize, page, totalCount)
		assert.Nil(t, err, "No error expected")
		assert.Equal(t, &expectedResp, resp, "Expected response does not match the expected one")
	})

	t.Run("Valid case (Full response)", func(t *testing.T) {
		pageSize := 10
		page := 3
		totalCount := 44
		data := []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		}

		nextPage := 4
		previousPage := 2
		expectedResp := PaginatedResponse{
			Data: data,
			Pagination: Pagination{
				CurrentPage:  page,
				NextPage:     &nextPage,
				PreviousPage: &previousPage,
				MaxPage:      5,
				RowsPerPage:  pageSize,
				TotalRows:    totalCount,
			},
		}
		resp, err := FormatPaginatedResponse(data, pageSize, page, totalCount)
		assert.Nil(t, err, "No error expected")
		assert.Equal(t, &expectedResp, resp, "Expected response does not match the expected one")
	})
}
