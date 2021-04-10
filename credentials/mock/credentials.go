package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"github.com/nhatthm/n26keychain/credentials"
)

// Mocker is KeychainCredentials mocker.
type Mocker func(tb testing.TB) *KeychainCredentials

// NoMock is no mock KeychainCredentials.
var NoMock = Mock()

var (
	_ credentials.KeychainCredentials         = (*KeychainCredentials)(nil)
	_ credentials.KeychainCredentialsProvider = (*KeychainCredentials)(nil)
)

// KeychainCredentials is a credentials.KeychainCredentials.
type KeychainCredentials struct {
	testifyMock.Mock
}

// Username satisfies credentials.KeychainCredentials.
func (p *KeychainCredentials) Username() string {
	return p.Called().String(0)
}

// Password satisfies credentials.KeychainCredentials.
func (p *KeychainCredentials) Password() string {
	return p.Called().String(0)
}

// Update satisfies credentials.KeychainCredentials.
func (p *KeychainCredentials) Update(username, password string) error {
	return p.Called(username, password).Error(0)
}

// Delete satisfies credentials.KeychainCredentials.
func (p *KeychainCredentials) Delete() error {
	return p.Called().Error(0)
}

// KeychainCredentials satisfies credentials.KeychainCredentialsProvider.
func (p *KeychainCredentials) KeychainCredentials() credentials.KeychainCredentials {
	return p
}

// mock mocks credentials.KeychainCredentialsProvider interface.
func mock(mocks ...func(p *KeychainCredentials)) *KeychainCredentials {
	p := &KeychainCredentials{}

	for _, m := range mocks {
		m(p)
	}

	return p
}

// Mock creates KeychainCredentials mock with cleanup to ensure all the expectations are met.
func Mock(mocks ...func(p *KeychainCredentials)) Mocker {
	return func(tb testing.TB) *KeychainCredentials {
		tb.Helper()

		p := mock(mocks...)

		tb.Cleanup(func() {
			assert.True(tb, p.Mock.AssertExpectations(tb))
		})

		return p
	}
}
