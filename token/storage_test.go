//go:build !integration
// +build !integration

package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/nhatthm/n26keychain"
	"github.com/nhatthm/n26keychain/mock"
	"github.com/nhatthm/n26keychain/test"
)

var tokenStorageKey = "john@example.com:54252481-ff0e-4903-9a9c-1886d16eab73"

func TestTokenStorage_Get(t *testing.T) {
	testCases := []struct {
		scenario      string
		mockStorage   mock.StorageMocker
		expectedToken auth.OAuthToken
		expectedError string
	}{
		{
			scenario: "token not found",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", tokenStorageKey).
					Return("", keyring.ErrNotFound)
			}),
		},
		{
			scenario: "could not get token",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", tokenStorageKey).
					Return("", errors.New("get error"))
			}),
			expectedError: "get error",
		},
		{
			scenario: "invalid token",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", tokenStorageKey).
					Return("{", nil)
			}),
			expectedError: `could not unmarshal token: unexpected end of JSON input`,
		},
		{
			scenario: "success",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", tokenStorageKey).
					Return(`{"access_token":"access","refresh_token":"refresh","expires_at":"2020-01-02T03:04:05.000Z","refresh_expires_at":"2020-01-02T04:04:05.000Z"}`, nil)
			}),
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
			p := NewStorage(WithKeyring(tc.mockStorage(t)))
			token, err := p.Get(context.Background(), tokenStorageKey)

			assert.Equal(t, tc.expectedToken, token)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTokenStorage_GetKeyring(t *testing.T) {
	expect := func(t *testing.T, s n26keychain.Storage) { // nolint: thelper
		err := s.Set(tokenStorageKey, `{"access_token":"access","refresh_token":"refresh","expires_at":"2020-01-02T03:04:05.000Z","refresh_expires_at":"2020-01-02T04:04:05.000Z"}`)
		require.NoError(t, err)
	}

	expectedToken := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	test.Run(t, tokenStorageService, tokenStorageKey, expect, func(t *testing.T) { // nolint: thelper
		p := NewStorage()

		token, err := p.Get(context.Background(), tokenStorageKey)

		assert.Equal(t, expectedToken, token)
		assert.NoError(t, err)
	})
}

func TestTokenStorage_SetAndDeleteKeyring(t *testing.T) {
	expectedToken := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	test.Run(t, tokenStorageService, tokenStorageKey, nil, func(t *testing.T) { // nolint: thelper
		p := NewStorage()

		err := p.Set(context.Background(), tokenStorageKey, expectedToken)
		assert.NoError(t, err)

		// Get from keychain.
		data, err := keyring.Get(tokenStorageService, tokenStorageKey)
		expectedData := `{"access_token":"access","refresh_token":"refresh","expires_at":"2020-01-02T03:04:05Z","refresh_expires_at":"2020-01-02T04:04:05Z"}`

		assert.Equal(t, expectedData, data)
		assert.NoError(t, err)
	})
}

func TestTokenStorage_Delete(t *testing.T) {
	testCases := []struct {
		scenario      string
		mockStorage   mock.StorageMocker
		expectedError string
	}{
		{
			scenario: "token not found",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", tokenStorageKey).
					Return(keyring.ErrNotFound)
			}),
		},
		{
			scenario: "could not delete token",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", tokenStorageKey).
					Return(errors.New("delete error"))
			}),
			expectedError: "delete error",
		},
		{
			scenario: "success",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", tokenStorageKey).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			p := NewStorage(WithKeyring(tc.mockStorage(t)))
			err := p.Delete(context.Background(), tokenStorageKey)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTokenStorage_DeleteKeyring(t *testing.T) {
	token := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	test.Run(t, tokenStorageService, tokenStorageKey, nil, func(t *testing.T) { // nolint: thelper
		p := NewStorage()

		// Prepare data.
		err := p.Set(context.Background(), tokenStorageKey, token)
		assert.NoError(t, err)

		// Verify data.
		_, err = keyring.Get(tokenStorageService, tokenStorageKey)
		assert.NoError(t, err)

		// Test.
		err = p.Delete(context.Background(), tokenStorageKey)
		assert.NoError(t, err)

		// Verify.
		_, err = keyring.Get(tokenStorageService, tokenStorageKey)
		assert.Equal(t, keyring.ErrNotFound, err)
	})
}
