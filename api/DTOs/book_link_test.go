package DTOs

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewBookLinkDTO(t *testing.T) {
	t.Run("BookLinkDTO generation", func(t *testing.T) {
		testRawUrl := "https://test.com"
		testUrl, err := url.Parse(testRawUrl)
		if err != nil {
			t.Error("Test failed on setup, unable to create test URL")
		}
		expectedLinkDTO := BookLinkDTO{DownloadLink: testRawUrl}
		bookLinkDTO := NewBookLinkDTO(*testUrl)
		assert.Equal(t, expectedLinkDTO, bookLinkDTO, "Unexpected BookLinkDTO")
	})
}
