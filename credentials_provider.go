package n26keychain

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/nhatthm/n26api"
	"github.com/zalando/go-keyring"
)

const (
	credentialsService = "n26api.credentials"
	credentialsKey     = "default"
)

var mu sync.Mutex

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials returns credentials from keychain.
func Credentials() (n26api.CredentialsProvider, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := keyring.Get(credentialsService, credentialsKey)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return noCredentials(), nil
		}

		return nil, err
	}

	var cred credentials

	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}

	return n26api.Credentials(cred.Username, cred.Password), nil
}

// PersistCredentials persists the credentials to keychain.
func PersistCredentials(username, password string) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.Marshal(credentials{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}

	return keyring.Set(credentialsService, credentialsKey, string(data))
}

// DeleteCredentials deletes the credentials in keychain.
func DeleteCredentials() error {
	mu.Lock()
	defer mu.Unlock()

	return keyring.Delete(credentialsService, credentialsKey)
}

func noCredentials() n26api.CredentialsProvider {
	return n26api.Credentials("", "")
}
