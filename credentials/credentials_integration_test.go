// +build integration

package credentials

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keyring "github.com/zalando/go-keyring"

	"github.com/nhatthm/n26keychain"
	"github.com/nhatthm/n26keychain/test"
)

func TestIntegrationCredentials_LoadKeyringNotFound(t *testing.T) {
	deviceID := uuid.New()

	test.Run(t, credentialsService, deviceID.String(), nil, func(t *testing.T) { // nolint: thelper
		c := New(deviceID)

		assert.Empty(t, c.Username())
		assert.Empty(t, c.Password())
	})
}

func TestIntegrationCredentials_LoadKeyring(t *testing.T) {
	deviceID := uuid.New()

	expect := func(t *testing.T, s n26keychain.Storage) { // nolint: thelper
		err := s.Set(deviceID.String(), `{"username":"foo","password":"bar"}`)
		require.NoError(t, err)
	}

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), expect, func(t *testing.T) { // nolint: thelper
		c := New(deviceID)

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())
	})
}

func TestIntegrationCredentials_UpdateKeyring(t *testing.T) {
	deviceID := uuid.New()

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), nil, func(t *testing.T) { // nolint: thelper
		c := New(deviceID)

		_, err := keyring.Get(credentialsService, deviceID.String())
		require.Equal(t, keyring.ErrNotFound, err)

		err = c.Update("foo", "bar")
		assert.NoError(t, err)

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())

		// Get from keychain.
		data, err := keyring.Get(credentialsService, deviceID.String())
		expectedData := `{"username":"foo","password":"bar"}`

		assert.Equal(t, expectedData, data)
		assert.NoError(t, err)
	})
}

func TestIntegrationCredentials_DeleteKeyring(t *testing.T) {
	deviceID := uuid.New()

	expect := func(t *testing.T, s n26keychain.Storage) { // nolint: thelper
		err := s.Set(deviceID.String(), `{"username":"foo","password":"bar"}`)
		require.NoError(t, err)
	}

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), expect, func(t *testing.T) { // nolint: thelper
		c := New(deviceID)

		_, err := keyring.Get(credentialsService, deviceID.String())
		require.NoError(t, err)

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())

		// Delete.
		err = c.Delete()
		require.NoError(t, err)

		// The key should not be found anymore.
		_, err = keyring.Get(credentialsService, deviceID.String())
		assert.Equal(t, keyring.ErrNotFound, err)
	})
}

func TestIntegrationCredentials_DeleteKeyringNotFound(t *testing.T) {
	deviceID := uuid.New()

	test.Run(t, credentialsService, deviceID.String(), nil, func(t *testing.T) { // nolint: thelper
		c := New(deviceID)

		_, err := keyring.Get(credentialsService, deviceID.String())
		require.Equal(t, keyring.ErrNotFound, err)

		// Delete.
		err = c.Delete()
		require.NoError(t, err)
	})
}
