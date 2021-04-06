package n26keychain

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

var testMU sync.Mutex

func setupCredentialsScenario(t *testing.T, prepare func(t *testing.T)) {
	t.Helper()

	testMU.Lock()

	keyring.MockInit()

	// Before scenario.
	err := keyring.Delete(credentialsService, credentialsKey)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := keyring.Delete(credentialsService, credentialsKey)

		defer testMU.Unlock()

		if err != nil && !errors.Is(err, keyring.ErrNotFound) {
			t.Fatal(err)
		}
	})

	if prepare != nil {
		prepare(t)
	}
}

func TestCredentials(t *testing.T) {
	testCases := []struct {
		scenario         string
		prepare          func(t *testing.T)
		expectedUsername string
		expectedPassword string
		expectedError    string
	}{
		{
			scenario: "missing credentials",
		},
		{
			scenario: "credentials is in wrong format",
			prepare: func(t *testing.T) { // nolint: thelper
				err := keyring.Set(credentialsService, credentialsKey, "{")
				require.NoError(t, err)
			},
			expectedError: "unexpected end of JSON input",
		},
		{
			scenario: "success",
			prepare: func(t *testing.T) { // nolint: thelper
				err := keyring.Set(credentialsService, credentialsKey, `{"username":"foo","password":"bar"}`)
				require.NoError(t, err)
			},
			expectedUsername: "foo",
			expectedPassword: "bar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			setupCredentialsScenario(t, tc.prepare)

			c, err := Credentials()

			if tc.expectedError == "" {
				assert.NotNil(t, c)
				assert.Equal(t, tc.expectedUsername, c.Username())
				assert.Equal(t, tc.expectedPassword, c.Password())
				assert.NoError(t, err)
			} else {
				assert.Nil(t, c)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestPersistCredentials(t *testing.T) {
	setupCredentialsScenario(t, nil)

	err := PersistCredentials("foo", "bar")
	require.NoError(t, err)

	c, err := Credentials()
	require.NoError(t, err)

	expectedUsername := "foo"
	expectedPassword := "bar"

	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())
}

func TestDeleteCredentials(t *testing.T) {
	setupCredentialsScenario(t, nil)

	err := PersistCredentials("foo", "bar")
	require.NoError(t, err)

	err = DeleteCredentials()
	require.NoError(t, err)

	_, err = keyring.Get(credentialsService, credentialsKey)

	assert.Equal(t, keyring.ErrNotFound, err)
}
