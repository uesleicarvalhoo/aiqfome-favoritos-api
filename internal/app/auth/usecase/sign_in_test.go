package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	fixtureAuth "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/usecase"
	mocksJwt "github.com/uesleicarvalhoo/aiqfome/pkg/jwt/mocks"
	mocksPassword "github.com/uesleicarvalhoo/aiqfome/pkg/password/mocks"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	mocksUser "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestSignInUseCase_Execute(t *testing.T) {
	t.Parallel()

	userID := uuid.NextID()
	email := "user@email.com"
	passwd := "secret"
	passwdWithSalt := fmt.Sprintf("%s:%s", userID.String(), passwd)
	passwdHash := "secret-hash"

	paramsBuilder := fixtureAuth.AnySignInParams().
		WithEmail(email).
		WithPassword(passwd)

	userBuilder := fixtureUser.AnyUser().
		WithID(userID).
		WithPasswordHash(passwdHash)

	opts := usecase.SignInOptions{
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Hour,
	}

	testCases := []struct {
		about            string
		params           dto.SignInParams
		setupRepo        func(repo *mocksUser.Repository)
		setupHasher      func(h *mocksPassword.Hasher)
		setupAccessProv  func(p *mocksJwt.Provider)
		setupRefreshProv func(p *mocksJwt.Provider)
		expectedErr      string
		expectedTokens   dto.AuthTokens
	}{
		{
			about:       "when params are invalid",
			params:      paramsBuilder.WithEmail("invalid").Build(),
			expectedErr: "[AQF002] email: email inválido",
		},
		{
			about:  "when user not found",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, email).
					Return(user.User{}, user.ErrNotFound)
			},
			expectedErr: "[AQF003] user not found",
		},
		{
			about:  "when repo returns other error",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, email).
					Return(user.User{}, errors.New("db error"))
			},
			expectedErr: "[AQF004] error while trying to find user | cause: db error",
		},
		{
			about:  "when password is incorrect",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, email).
					Return(userBuilder.Build(), nil)
			},
			setupHasher: func(h *mocksPassword.Hasher) {
				h.On("Compare", passwdHash, passwdWithSalt).
					Return(errors.New("no match"))
			},
			expectedErr: "[AUT001] senha invalida",
		},
		{
			about:  "when credentials are valid",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, email).
					Return(userBuilder.Build(), nil)
			},
			setupHasher: func(h *mocksPassword.Hasher) {
				h.On("Compare", passwdHash, passwdWithSalt).
					Return(nil)
			},
			setupAccessProv: func(p *mocksJwt.Provider) {
				p.On("Generate", mock.Anything, userID.String(), opts.AccessTokenDuration).
					Return("tokA", nil)
			},
			setupRefreshProv: func(p *mocksJwt.Provider) {
				p.On("Generate", mock.Anything, userID.String(), opts.RefreshTokenDuration).
					Return("tokR", nil)
			},
			expectedTokens: dto.AuthTokens{
				AccessToken:  "tokA",
				RefreshToken: "tokR",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := mocksUser.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			hasher := mocksPassword.NewHasher(t)
			if tc.setupHasher != nil {
				tc.setupHasher(hasher)
			}

			accessProv := mocksJwt.NewProvider(t)
			if tc.setupAccessProv != nil {
				tc.setupAccessProv(accessProv)
			}

			refreshProv := mocksJwt.NewProvider(t)
			if tc.setupRefreshProv != nil {
				tc.setupRefreshProv(refreshProv)
			}

			uc := usecase.NewSignUseCase(repo, hasher, opts, accessProv, refreshProv)

			// Action
			res, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, dto.AuthTokens{}, res)
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTokens, res)

			repo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			accessProv.AssertExpectations(t)
			refreshProv.AssertExpectations(t)
		})
	}
}
