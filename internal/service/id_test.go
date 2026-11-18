package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thecodest22/YP_url-shortener/internal/testutil"
)

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "zero length", length: 0},
		{name: "single character", length: 1},
		{name: "usual length", length: 8},
		{name: "long ID", length: 64},
	}

	for i, test := range tests {
		t.Run(testutil.GetSubtestTitle(i+1, test.name), func(t *testing.T) {
			id := generateID(test.length)

			assert.Len(t, id, test.length)
			assert.Regexp(t, `^[A-Za-z0-9]*$`, id)
		})
	}
}
