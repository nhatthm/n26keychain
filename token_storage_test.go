package n26keychain

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

var tokenStorageUser = "john@example.com"

func setupTokenStorageScenario(t *testing.T, prepare func(t *testing.T)) {
	t.Helper()

	testMU.Lock()

	keyring.MockInit()

	// Before scenario.
	err := keyring.Delete(credentialsService, tokenStorageUser)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := keyring.Delete(credentialsService, tokenStorageUser)

		defer testMU.Unlock()

		if err != nil && !errors.Is(err, keyring.ErrNotFound) {
			t.Fatal(err)
		}
	})

	if prepare != nil {
		prepare(t)
	}
}

func TestTokenStorage_Get(t *testing.T) {
	testCases := []struct {
		scenario      string
		prepare       func(t *testing.T)
		expectedToken auth.OAuthToken
		expectedError string
	}{
		{
			scenario: "token not found",
		},
		{
			scenario: "invalid token",
			prepare: func(t *testing.T) { // nolint: thelper
				err := keyring.Set(tokenStorageService, tokenStorageUser, "{")
				require.NoError(t, err)
			},
			expectedError: `could not unmarshal token: unexpected end of JSON input`,
		},
		{
			scenario: "success",
			prepare: func(t *testing.T) { // nolint: thelper
				err := keyring.Set(tokenStorageService, tokenStorageUser, `{"access_token":"access","refresh_token":"refresh","expires_at":"2020-01-02T03:04:05.000Z","refresh_expires_at":"2020-01-02T04:04:05.000Z"}`)
				require.NoError(t, err)
			},
			expectedToken: auth.OAuthToken{
				AccessToken:      "access",
				RefreshToken:     "refresh",
				ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
				RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			setupTokenStorageScenario(t, tc.prepare)

			p := TokenStorage()
			token, err := p.Get(context.Background(), tokenStorageUser)

			assert.Equal(t, tc.expectedToken, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTokenStorage_Set(t *testing.T) {
	setupTokenStorageScenario(t, nil)

	expected := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	p := TokenStorage()

	err := p.Set(context.Background(), tokenStorageUser, expected)
	require.NoError(t, err)

	token, err := p.Get(context.Background(), tokenStorageUser)

	assert.Equal(t, expected, token)
	assert.NoError(t, err)
}
