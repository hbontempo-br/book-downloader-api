package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ResponseWriterMock struct {
	pureJSONCalls int
	resp          interface{}
	status        int
}

func (c *ResponseWriterMock) PureJSON(status int, response interface{}) {
	c.pureJSONCalls += 1
	c.resp = response
	c.status = status
}

func TestDefaultErrorMessage(t *testing.T) {

	t.Run("Existent status (200: OK)", func(t *testing.T) {
		httpStatus := 200
		details := map[string]interface{}{
			"in": map[string]string{
				"ter": "face",
			},
		}
		responseWriter := ResponseWriterMock{pureJSONCalls: 0}

		expectedResponse := map[string]interface{}{
			"code":    httpStatus,
			"message": "OK",
			"details": details,
		}

		err := DefaultErrorMessage(&responseWriter, httpStatus, details)
		assert.Nil(t, err, "Unexpected error")
		assert.Equal(t, 1, responseWriter.pureJSONCalls, "Only one call to responseWriter.PureJson was expected")
		assert.Equal(t, httpStatus, responseWriter.status, "Unexpected status received on responseWriter.PureJson")
		assert.Equal(t, expectedResponse, responseWriter.resp, "Unexpected message received on responseWriter.PureJson")
	})

	t.Run("Nonexistent status", func(t *testing.T) {
		httpStatus := 5000 // Nonexistent status
		responseWriter := ResponseWriterMock{pureJSONCalls: 0}
		err := DefaultErrorMessage(&responseWriter, httpStatus, "")
		assert.EqualError(t, err, InvalidHTTPStatus.Error(), "Expected an InvalidHTTPStatus error")
		assert.Equal(t, 0, responseWriter.pureJSONCalls, "No call to responseWriter.PureJson was expected")
	})

}
