package n26keychain

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bool64/ctxd"
	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/zalando/go-keyring"
)

const tokenStorageService = "n26api.token" // nolint: gosec

var _ auth.TokenStorage = (*tokenStorage)(nil)

// tokenStorage provides token from keychain.
type tokenStorage struct{}

// Get gets token from keychain.
func (t *tokenStorage) Get(ctx context.Context, key string) (auth.OAuthToken, error) {
	data, err := keyring.Get(tokenStorageService, key)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return auth.OAuthToken{}, nil
		}

		return auth.OAuthToken{}, err
	}

	var token auth.OAuthToken

	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return auth.OAuthToken{}, ctxd.WrapError(ctx, err, "could not unmarshal token")
	}

	return token, nil
}

// Set persists token to keychain.
func (t *tokenStorage) Set(ctx context.Context, key string, token auth.OAuthToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return ctxd.WrapError(ctx, err, "could not marshal token")
	}

	return keyring.Set(tokenStorageService, key, string(data))
}

// TokenStorage returns keychain as a token storage.
func TokenStorage() auth.TokenStorage {
	return &tokenStorage{}
}
