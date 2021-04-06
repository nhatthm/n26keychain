package n26keychain

import "github.com/nhatthm/n26api"

// WithTokenStorage sets keychain as a token storage for n26 client.
func WithTokenStorage() n26api.Option {
	return n26api.WithTokenStorage(TokenStorage())
}
