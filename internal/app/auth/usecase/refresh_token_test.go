package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	jwtFixture "github.com/uesleicarvalhoo/aiqfome/pkg/jwt/fixture"
	jwtMocks "github.com/uesleicarvalhoo/aiqfome/pkg/jwt/mocks"
)

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	t.Parallel()

	refreshToken := "refresh-token"
	accessToken := "new-access-token"
	newRefresh := "new-refresh-token"
	clientID := jwtFixture.AnyClaims().Build().UserID
	opts := usecase.RefreshTokenOptions{
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Hour,
	}

	claimsBuilder := jwtFixture.AnyClaims().WithClientID(clientID)

	testCases := []struct {
		about             string
		setupRefreshProv  func(p *jwtMocks.Provider)
		setupAccessProv   func(p *jwtMocks.Provider)
		params            dto.RefreshTokenParams
		expectedErrSubstr string
		expectedTokens    dto.AuthTokens
	}{
		{
			about: "when refresh token is invalid",
			setupRefreshProv: func(p *jwtMocks.Provider) {
				p.On("Validate", mock.Anything, refreshToken).
					Return(jwt.Claims{}, errors.New("invalid refresh"))
			},
			params:            dto.RefreshTokenParams{RefreshToken: refreshToken},
			expectedErrSubstr: "invalid refresh",
		},
		{
			about: "when refresh is valid and token generation succeeds",
			setupRefreshProv: func(p *jwtMocks.Provider) {
				p.On("Validate", mock.Anything, refreshToken).
					Return(claimsBuilder.Build(), nil)
				p.On("Generate", mock.Anything, clientID.String(), opts.RefreshTokenDuration).
					Return(newRefresh, nil)
			},
			setupAccessProv: func(p *jwtMocks.Provider) {
				p.On("Generate", mock.Anything, clientID.String(), opts.AccessTokenDuration).
					Return(accessToken, nil)
			},
			params:         dto.RefreshTokenParams{RefreshToken: refreshToken},
			expectedTokens: dto.AuthTokens{AccessToken: accessToken, RefreshToken: newRefresh},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			refreshProv := jwtMocks.NewProvider(t)
			tc.setupRefreshProv(refreshProv)

			accessProv := jwtMocks.NewProvider(t)
			if tc.setupAccessProv != nil {
				tc.setupAccessProv(accessProv)
			}

			uc := usecase.NewRefreshTokenUseCase(opts, accessProv, refreshProv)

			// Action
			tokens, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErrSubstr != "" {
				assert.Equal(t, dto.AuthTokens{}, tokens)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrSubstr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTokens, tokens)
			refreshProv.AssertExpectations(t)
			accessProv.AssertExpectations(t)
		})
	}
}
