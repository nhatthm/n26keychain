package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26keychain"
)

// StorageMocker is Storage mocker.
type StorageMocker func(tb testing.TB) *Storage

// NoMockStorage is no mock Storage.
var NoMockStorage = MockStorage()

var _ n26keychain.Storage = (*Storage)(nil)

// Storage is a n26keychain.Storage.
type Storage struct {
	mock.Mock
}

// Set satisfies n26keychain.Storage.
func (s *Storage) Set(user, password string) error {
	return s.Called(user, password).Error(0)
}

// Get satisfies n26keychain.Storage.
func (s *Storage) Get(user string) (string, error) {
	ret := s.Called(user)

	return ret.String(0), ret.Error(1)
}

// Delete satisfies n26keychain.Storage.
func (s *Storage) Delete(user string) error {
	return s.Called(user).Error(0)
}

// mockStorage mocks n26keychain.Storage interface.
func mockStorage(mocks ...func(s *Storage)) *Storage {
	s := &Storage{}

	for _, m := range mocks {
		m(s)
	}

	return s
}

// MockStorage creates Storage mock with cleanup to ensure all the expectations are met.
func MockStorage(mocks ...func(s *Storage)) StorageMocker { // nolint: revive
	return func(tb testing.TB) *Storage {
		tb.Helper()

		s := mockStorage(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, s.Mock.AssertExpectations(tb))
		})

		return s
	}
}
