package mock_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26keychain/credentials/mock"
)

func TestKeychainCredentials_Credentials(t *testing.T) {
	t.Parallel()

	p := mock.Mock(func(p *mock.KeychainCredentials) {
		p.On("Username").Return("username")
		p.On("Password").Return("password")
	})(t)

	expectedUsername := "username"
	expectedPassword := "password"

	assert.Equal(t, expectedUsername, p.Username())
	assert.Equal(t, expectedPassword, p.Password())
}

func TestKeychainCredentials_Update(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mock          mock.Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mock: mock.Mock(func(p *mock.KeychainCredentials) {
				p.On("Update", "username", "password").
					Return(errors.New("update error"))
			}),
			expectedError: "update error",
		},
		{
			scenario: "success",
			mock: mock.Mock(func(p *mock.KeychainCredentials) {
				p.On("Update", "username", "password").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := tc.mock(t)
			err := p.Update("username", "password")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestKeychainCredentials_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mock          mock.Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mock: mock.Mock(func(p *mock.KeychainCredentials) {
				p.On("Delete").
					Return(errors.New("delete error"))
			}),
			expectedError: "delete error",
		},
		{
			scenario: "success",
			mock: mock.Mock(func(p *mock.KeychainCredentials) {
				p.On("Delete").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := tc.mock(t)
			err := p.Delete()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestKeychainCredentials_KeychainCredentials(t *testing.T) {
	t.Parallel()

	p := mock.NoMock(t)

	assert.Equal(t, p, p.KeychainCredentials())
}
