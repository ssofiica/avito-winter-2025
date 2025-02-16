package token

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJWT(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		duration string
		err      error
		want     JWT
	}{
		{
			name:     "Valid duration",
			secret:   "secret",
			duration: "2h",
			err:      nil,
			want: JWT{
				Secret:  []byte("secret"),
				ExpTime: time.Duration(2 * time.Hour),
			},
		},
		{
			name:     "Invalid duration",
			secret:   "secret",
			duration: "",
			err:      errors.New("time: invalid duration \"\""),
			want: JWT{
				Secret:  []byte("secret"),
				ExpTime: time.Duration(48 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtInstance, err := NewJWT(tt.secret, tt.duration)

			assert.Equal(t, err, tt.err)
			assert.Equal(t, jwtInstance, tt.want)
		})
	}
}
