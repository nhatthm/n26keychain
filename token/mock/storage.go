package mock

import (
	"context"
	"testing"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26keychain/token"
)

// Mocker is Storage mocker.
type Mocker func(tb testing.TB) *Storage

// NoMock is no mock Storage.
var NoMock = Mock()

var _ token.KeychainStorage = (*Storage)(nil)

// Storage is a token.KeychainStorage.
type Storage struct {
	testifyMock.Mock
}

// Get satisfies token.KeychainStorage.
func (s *Storage) Get(ctx context.Context, key string) (auth.OAuthToken, error) {
	ret := s.Called(ctx, key)

	return ret.Get(0).(auth.OAuthToken), ret.Error(1)
}

// Set satisfies token.KeychainStorage.
func (s *Storage) Set(ctx context.Context, key string, token auth.OAuthToken) error {
	return s.Called(ctx, key, token).Error(0)
}

// Delete satisfies token.KeychainStorage.
func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.Called(ctx, key).Error(0)
}

// mock mocks token.Storage interface.
func mock(mocks ...func(s *Storage)) *Storage {
	s := &Storage{}

	for _, m := range mocks {
		m(s)
	}

	return s
}

// Mock creates Storage mock with cleanup to ensure all the expectations are met.
func Mock(mocks ...func(s *Storage)) Mocker {
	return func(tb testing.TB) *Storage {
		tb.Helper()

		s := mock(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, s.Mock.AssertExpectations(tb))
		})

		return s
	}
}
