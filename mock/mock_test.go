package mock_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26keychain/mock"
)

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockStorage    mock.StorageMocker
		expectedResult string
		expectedError  string
	}{
		{
			scenario: "error is not nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "foo").Return("", errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "error is nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Get", "foo").Return("bar", nil)
			}),
			expectedResult: "bar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			data, err := tc.mockStorage(t).Get("foo")

			assert.Equal(t, tc.expectedResult, data)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockStorage    mock.StorageMocker
		expectedResult string
		expectedError  string
	}{
		{
			scenario: "error is not nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Set", "foo", "bar").Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "error is nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Set", "foo", "bar").Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockStorage(t).Set("foo", "bar")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockStorage    mock.StorageMocker
		expectedResult string
		expectedError  string
	}{
		{
			scenario: "error is not nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", "foo").Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "error is nil",
			mockStorage: mock.MockStorage(func(s *mock.Storage) {
				s.On("Delete", "foo").Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			err := tc.mockStorage(t).Delete("foo")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
