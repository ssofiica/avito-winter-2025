package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		err      error
	}{
		{
			name:     "Success",
			password: "12345A",
			err:      nil,
		},
		{
			name:     "Fail",
			password: "111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			err:      errors.New("bcrypt: password length exceeds 72 bytes"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := Password("")
			err := pass.Hash(tt.password)

			assert.Equal(t, err, tt.err)
		})
	}
}
