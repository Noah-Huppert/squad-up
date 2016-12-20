package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := APIError{"errorId", "errorMsg", 123}
	assert.Equal(t, err.Error(), "errorMsg (errorId: 123)")
}
