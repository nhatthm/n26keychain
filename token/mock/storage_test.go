package mock_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nhatthm/n26api/pkg/auth"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/n26keychain/token/mock"
)

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	token := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	testCases := []struct {
		scenario       string
		mock           mock.Mocker
		expectedResult auth.OAuthToken
		expectedError  string
	}{
		{
			scenario: "error",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Get", context.Background(), "key").
					Return(auth.OAuthToken{}, errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "success",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Get", context.Background(), "key").
					Return(token, nil)
			}),
			expectedResult: token,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mock(t)
			result, err := s.Get(context.Background(), "key")

			assert.Equal(t, tc.expectedResult, result)

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

	token := auth.OAuthToken{
		AccessToken:      "access",
		RefreshToken:     "refresh",
		ExpiresAt:        time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		RefreshExpiresAt: time.Date(2020, 1, 2, 4, 4, 5, 0, time.UTC),
	}

	testCases := []struct {
		scenario      string
		mock          mock.Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Set", context.Background(), "key", token).
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "success",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Set", context.Background(), "key", token).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mock(t)
			err := s.Set(context.Background(), "key", token)

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
		scenario      string
		mock          mock.Mocker
		expectedError string
	}{
		{
			scenario: "error",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Delete", context.Background(), "key").
					Return(errors.New("error"))
			}),
			expectedError: "error",
		},
		{
			scenario: "success",
			mock: mock.Mock(func(s *mock.Storage) {
				s.On("Delete", context.Background(), "key").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mock(t)
			err := s.Delete(context.Background(), "key")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
