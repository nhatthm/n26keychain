package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentials_KeychainCredentials(t *testing.T) {
	t.Parallel()

	c := New()

	assert.Equal(t, c, c.KeychainCredentials())
}
