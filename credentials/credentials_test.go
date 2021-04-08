package credentials

import (
	"errors"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/nhatthm/n26keychain"
	"github.com/nhatthm/n26keychain/mock"
	"github.com/nhatthm/n26keychain/test"
)

func TestCredentials_Load(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario         string
		mockStorage      mock.StorageMocker
		expectedUsername string
		expectedPassword string
		expectedError    string
	}{
		{
			scenario: "missing credentials",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "default").Return("", keyring.ErrNotFound)
			}),
		},
		{
			scenario: "could not get credentials",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "default").Return("", errors.New("get error"))
			}),
			expectedError: "error: could not get credentials {\"error\":{}}\n",
		},
		{
			scenario: "credentials is in wrong format",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "default").Return("{", nil)
			}),
			expectedError: "error: could not unmarshal credentials {\"error\":{\"Offset\":1}}\n",
		},
		{
			scenario: "success",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "default").Return(`{"username":"foo","password":"bar"}`, nil)
			}),
			expectedUsername: "foo",
			expectedPassword: "bar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockStorage(t)
			l := &ctxd.LoggerMock{}

			c := New(
				WithStorage(s),
				WithLogger(l),
			)

			assert.Equal(t, tc.expectedUsername, c.Username())
			assert.Equal(t, tc.expectedPassword, c.Password())
			assert.Equal(t, tc.expectedError, l.String())
		})
	}
}

func TestCredentials_LoadOnce(t *testing.T) {
	expectedUsername := "foo"
	expectedPassword := "bar"

	storage := mock.MockStorage(func(s *mock.Storage) {
		s.On("Get", "default").
			Return(`{"username":"foo","password":"bar"}`, nil).
			Once()
	})(t)

	c := New(WithStorage(storage))

	// 1st run calls storage.
	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())

	// 2nd run does not call storage.
	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())
}

func TestCredentials_LoadKeyring(t *testing.T) {
	deviceID := uuid.New()

	expect := func(t *testing.T, s n26keychain.Storage) { // nolint: thelper
		err := s.Set(deviceID.String(), `{"username":"foo","password":"bar"}`)
		require.NoError(t, err)
	}

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), expect, func(t *testing.T) { // nolint: thelper
		c := New(WithDeviceID(deviceID))

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())
	})
}

func TestCredentials_Update(t *testing.T) {
	t.Parallel()

	username := "foo"
	password := "bar"

	testCases := []struct {
		scenario         string
		mockStorage      mock.StorageMocker
		expectedUsername string
		expectedPassword string
		expectedError    string
	}{
		{
			scenario: "could not update",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Set", "default", `{"username":"foo","password":"bar"}`).
					Return(errors.New("update error"))
			}),
			expectedError: "update error",
		},
		{
			scenario: "success",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Set", "default", `{"username":"foo","password":"bar"}`).
					Return(nil)
			}),
			expectedUsername: "foo",
			expectedPassword: "bar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockStorage(t)
			c := New(WithStorage(s))

			err := c.Update(username, password)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUsername, c.Username())
				assert.Equal(t, tc.expectedPassword, c.Password())
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCredentials_UpdateOnce(t *testing.T) {
	storage := mock.MockStorage(func(s *mock.Storage) {
		s.On("Get", "default").
			Return(`{"username":"foo","password":"bar"}`, nil).
			Once()

		s.On("Set", "default", `{"username":"john","password":"doe"}`).
			Return(nil)
	})(t)

	c := New(WithStorage(storage))

	// 1st run calls storage.
	expectedUsername := "foo"
	expectedPassword := "bar"

	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())

	// Update.
	err := c.Update("john", "doe")
	require.NoError(t, err)

	// 2nd run does not call storage.
	expectedUsername = "john"
	expectedPassword = "doe"

	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())
}

func TestCredentials_UpdateKeyring(t *testing.T) {
	deviceID := uuid.New()

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), nil, func(t *testing.T) { // nolint: thelper
		c := New(WithDeviceID(deviceID))

		err := c.Update("foo", "bar")
		assert.NoError(t, err)

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())

		// Get from keychain.
		data, err := keyring.Get(credentialsService, deviceID.String())
		expectedData := `{"username":"foo","password":"bar"}`

		assert.Equal(t, expectedData, data)
		assert.NoError(t, err)
	})
}

func TestCredentials_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockStorage   mock.StorageMocker
		expectedError string
	}{
		{
			scenario: "error not found",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", "default").Return(keyring.ErrNotFound)
			}),
		},
		{
			scenario: "could not delete",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", "default").Return(errors.New("delete error"))
			}),
			expectedError: "delete error",
		},
		{
			scenario: "success",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", "default").Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockStorage(t)
			c := New(WithStorage(s))

			err := c.Delete()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCredentials_DeleteOnce(t *testing.T) {
	storage := mock.MockStorage(func(s *mock.Storage) {
		s.On("Get", "default").
			Return(`{"username":"foo","password":"bar"}`, nil).
			Once()

		s.On("Delete", "default").Return(nil).Once()

		s.On("Get", "default").
			Return("", keyring.ErrNotFound).
			Once()
	})(t)

	c := New(WithStorage(storage))

	// 1st run calls storage.
	expectedUsername := "foo"
	expectedPassword := "bar"

	assert.Equal(t, expectedUsername, c.Username())
	assert.Equal(t, expectedPassword, c.Password())

	// Delete.
	err := c.Delete()
	require.NoError(t, err)

	// 2nd run calls storage again.
	assert.Empty(t, c.Username())
	assert.Empty(t, c.Password())
}

func TestCredentials_DeleteKeyring(t *testing.T) {
	deviceID := uuid.New()

	expect := func(t *testing.T, s n26keychain.Storage) { // nolint: thelper
		err := s.Set(deviceID.String(), `{"username":"foo","password":"bar"}`)
		require.NoError(t, err)
	}

	expectedUsername := "foo"
	expectedPassword := "bar"

	test.Run(t, credentialsService, deviceID.String(), expect, func(t *testing.T) { // nolint: thelper
		c := New(WithDeviceID(deviceID))

		assert.Equal(t, expectedUsername, c.Username())
		assert.Equal(t, expectedPassword, c.Password())

		// Delete.
		err := c.Delete()
		require.NoError(t, err)

		// The key should not be found anymore.
		_, err = keyring.Get(credentialsService, deviceID.String())
		assert.Equal(t, keyring.ErrNotFound, err)
	})
}
