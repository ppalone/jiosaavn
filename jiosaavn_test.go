package jiosaavn_test

import (
	"testing"

	"github.com/ppalone/jiosaavn"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := jiosaavn.NewClient(nil)
	assert.NotNil(t, c)
}
