package credentials

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCredentials_KeychainCredentials(t *testing.T) {
	t.Parallel()

	c := New(uuid.New())

	assert.Equal(t, c, c.KeychainCredentials())
}
