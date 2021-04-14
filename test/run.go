package test

import (
	"errors"
	"sync"
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/nhatthm/n26keychain"
)

var mu sync.Mutex

// RunExpect is an expect to run Run.
type RunExpect func(t *testing.T, s n26keychain.Storage)

func runTest(
	t *testing.T,
	service, key string,
	expect RunExpect,
	test func(t *testing.T),
) {
	t.Helper()

	mu.Lock()

	s := n26keychain.NewStorage(service)

	// Before scenario.
	err := s.Delete(key)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := s.Delete(key)

		defer mu.Unlock()

		if err != nil && !errors.Is(err, keyring.ErrNotFound) {
			t.Fatal(err)
		}
	})

	if expect != nil {
		expect(t, s)
	}

	if test != nil {
		test(t)
	}
}
