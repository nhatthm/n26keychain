package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/bool64/ctxd"
	"github.com/google/uuid"
	"github.com/nhatthm/n26api"
	"github.com/zalando/go-keyring"

	"github.com/nhatthm/n26keychain"
)

var _ KeychainCredentials = (*Credentials)(nil)

// KeychainCredentials manages credentials in keychain.
type KeychainCredentials interface {
	n26api.CredentialsProvider

	// Update persists new credentials to keychain.
	Update(username, password string) error
	// Delete deletes the credentials in keychain.
	Delete() error
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Option configures Credentials.
type Option func(p *Credentials)

// Credentials provides credentials from keychain.
type Credentials struct {
	storage n26keychain.Storage
	logger  ctxd.Logger

	mu sync.Mutex

	key      string
	loaded   bool
	username string
	password string
}

func (c *Credentials) load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.loaded = true
	c.username = ""
	c.password = ""

	data, err := c.storage.Get(c.key)
	if err != nil {
		if !errors.Is(err, keyring.ErrNotFound) {
			c.logger.Error(context.Background(), "could not get credentials", "error", err)
		}

		return
	}

	var t credentials

	if err := json.Unmarshal([]byte(data), &t); err != nil {
		c.logger.Error(context.Background(), "could not unmarshal credentials", "error", err)

		return
	}

	c.loaded = true
	c.username = t.Username
	c.password = t.Password
}

// Username returns the username from keychain.
func (c *Credentials) Username() string {
	if c.loaded {
		return c.username
	}

	c.load()

	return c.username
}

// Password returns the password from keychain.
func (c *Credentials) Password() string {
	if c.loaded {
		return c.password
	}

	c.load()

	return c.password
}

// Update persists new credentials to keychain.
func (c *Credentials) Update(username, password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(credentials{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	if err := c.storage.Set(c.key, string(data)); err != nil {
		return err
	}

	c.loaded = true
	c.username = username
	c.password = password

	return nil
}

// Delete deletes the credentials in keychain.
func (c *Credentials) Delete() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.storage.Delete(c.key); err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}

	c.loaded = false
	c.username = ""
	c.password = ""

	return nil
}

// New initiates a new Credentials.
func New(deviceID uuid.UUID, options ...Option) *Credentials {
	c := &Credentials{
		storage: n26keychain.NewStorage(credentialsService),
		logger:  ctxd.NoOpLogger{},

		key: deviceID.String(),
	}

	for _, o := range options {
		o(c)
	}

	return c
}

// WithStorage sets storage for Credentials.
func WithStorage(storage n26keychain.Storage) Option {
	return func(p *Credentials) {
		p.storage = storage
	}
}

// WithLogger sets logger for Credentials.
func WithLogger(logger ctxd.Logger) Option {
	return func(p *Credentials) {
		p.logger = logger
	}
}

// WithCredentialsProvider sets keychain as a credential provider.
func WithCredentialsProvider(options ...Option) n26api.Option {
	return func(c *n26api.Client) {
		n26api.WithCredentialsProvider(New(c.DeviceID(), options...))(c)
	}
}
