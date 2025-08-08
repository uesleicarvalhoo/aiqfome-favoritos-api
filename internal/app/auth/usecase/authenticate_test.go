package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/usecase"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	"github.com/uesleicarvalhoo/aiqfome/user/mocks"

	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	fixtureJwt "github.com/uesleicarvalhoo/aiqfome/pkg/jwt/fixture"
	mocksJwt "github.com/uesleicarvalhoo/aiqfome/pkg/jwt/mocks"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestAuthenticateUseCase_Execute(t *testing.T) {
	t.Parallel()

	token := "any-token"
	userID := uuid.NextID()

	userBuilder := fixtureUser.AnyUser().
		WithID(userID)

	claimsBuilder := fixtureJwt.AnyClaims().
		WithClientID(userID)

	testCases := []struct {
		about         string
		setupProvider func(provider *mocksJwt.Provider)
		setupRepo     func(repo *mocks.Repository)
		expectedErr   string
		expecteduser  user.User
	}{
		{
			about: "when token is invalid",
			setupProvider: func(provider *mocksJwt.Provider) {
				provider.On("Validate", mock.Anything, token).
					Return(jwt.Claims{}, errors.New("invalid token"))
			},
			expectedErr: "invalid token",
		},
		{
			about: "when user not found",
			setupProvider: func(provider *mocksJwt.Provider) {
				provider.On("Validate", mock.Anything, token).
					Return(claimsBuilder.Build(), nil)
			},
			setupRepo: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything, userID).
					Return(user.User{}, user.ErrNotFound)
			},
			expectedErr: "[AQF003] user not found",
		},
		{
			about: "when repo returns other error",
			setupProvider: func(provider *mocksJwt.Provider) {
				provider.On("Validate", mock.Anything, token).
					Return(claimsBuilder.Build(), nil)
			},
			setupRepo: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything, userID).
					Return(user.User{}, errors.New("i'm an repository error"))
			},
			expectedErr: "[AQF004] error while trying to find user",
		},
		{
			about: "when user is inactive",
			setupProvider: func(provider *mocksJwt.Provider) {
				provider.On("Validate", mock.Anything, token).
					Return(claimsBuilder.Build(), nil)
			},
			setupRepo: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything, userID).
					Return(userBuilder.WithActive(false).Build(), nil)
			},
			expectedErr: "[USR002] usuário bloqueado",
		},
		{
			about: "when all is valid",
			setupProvider: func(provider *mocksJwt.Provider) {
				provider.On("Validate", mock.Anything, token).
					Return(claimsBuilder.Build(), nil)
			},
			setupRepo: func(repo *mocks.Repository) {
				repo.On("Find", mock.Anything, userID).
					Return(userBuilder.Build(), nil)
			},
			expecteduser: userBuilder.Build(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := mocks.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			provider := mocksJwt.NewProvider(t)
			if tc.setupProvider != nil {
				tc.setupProvider(provider)
			}

			uc := usecase.NewAuthenticateUseCase(repo, provider)

			// Action
			res, err := uc.Execute(context.Background(), token)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, user.User{}, res)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expecteduser, res)

			repo.AssertExpectations(t)
			provider.AssertExpectations(t)
		})
	}
}
