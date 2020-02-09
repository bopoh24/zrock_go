package apiserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_CreateJWT(t *testing.T) {
	token, err := CreateJWT(10)
	assert.NoError(t, err)
	assert.NotEqual(t, "", token)
}

func TestAuth_ParseJWT(t *testing.T) {
	token, _ := CreateJWT(10)
	testCases := []struct {
		name    string
		token   string
		isValid bool
	}{
		{
			name:    "empty token",
			token:   "",
			isValid: false,
		},
		{
			name:    "incorrect data",
			token:   "someToken",
			isValid: false,
		},
		{
			name:    "valid",
			token:   token,
			isValid: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := ParseJWT(tc.token)
			if tc.isValid {
				assert.NoError(t, err)
				assert.Equal(t, 10, userID)
			} else {
				assert.Error(t, err)
				assert.Equal(t, 0, userID)
			}
		})
	}
}
