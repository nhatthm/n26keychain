package n26keychain

import (
	"github.com/zalando/go-keyring"
)

var _ Storage = (*storage)(nil)

// Storage is a keychain storage.
type Storage interface {
	// Set sets password in keychain for user.
	Set(user, password string) error
	// Get gets password from keychain.
	Get(user string) (string, error)
	// Delete deletes secret from keychain.
	Delete(user string) error
}

type storage struct {
	service string
}

// Set sets password in keychain for user.
func (s *storage) Set(user, password string) error {
	return keyring.Set(s.service, user, password)
}

// Get gets password from keychain.
func (s *storage) Get(user string) (string, error) {
	return keyring.Get(s.service, user)
}

// Delete deletes secret from keychain.
func (s *storage) Delete(user string) error {
	return keyring.Delete(s.service, user)
}

// NewStorage creates a keychain storage.
func NewStorage(service string) Storage {
	return &storage{
		service: service,
	}
}
