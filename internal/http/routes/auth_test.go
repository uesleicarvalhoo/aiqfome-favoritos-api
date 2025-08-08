package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
	authMocks "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/mocks"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/test"
	"github.com/uesleicarvalhoo/aiqfome/user"
	userFixture "github.com/uesleicarvalhoo/aiqfome/user/fixture"
)

func Test_signIn(t *testing.T) {
	t.Parallel()

	paramsBuilder := fixture.AnySignInParams().
		WithEmail("user@email.com").
		WithPassword("i'm secret")

	tokensBuilder := fixture.AnyAuthTokens()

	testCases := []struct {
		about           string
		params          dto.SignInParams
		setupUC         func(uc *authMocks.SignInUseCase)
		expectedStatus  int
		expectedTokens  *dto.AuthTokens
		expectedErrCode string
	}{
		{
			about:  "when params are invalid",
			params: paramsBuilder.WithEmail("invalid").WithPassword("").Build(),
			setupUC: func(uc *authMocks.SignInUseCase) {
				err := domainerror.New(domainerror.InvalidParams, "email: email inválido; password: campo obrigatório", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.WithEmail("invalid").WithPassword("").Build()).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedErrCode: string(domainerror.InvalidParams),
		},
		{
			about:  "when user not found",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignInUseCase) {
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(dto.AuthTokens{}, domainerror.New(domainerror.ResourceNotFound, "user not found", nil))
			},
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: string(domainerror.ResourceNotFound),
		},
		{
			about:  "when password is invalid",
			params: paramsBuilder.WithPassword("wrong").Build(),
			setupUC: func(uc *authMocks.SignInUseCase) {
				err := domainerror.New(domainerror.InvalidPassword, "senha invalida", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.WithPassword("wrong").Build()).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedErrCode: string(domainerror.InvalidPassword),
		},
		{
			about:  "when dependency error happens",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignInUseCase) {
				err := domainerror.Wrap(errors.New("db error"), domainerror.DependecyError, "error while trying to find user", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: string(domainerror.DependecyError),
		},
		{
			about:  "when ok",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignInUseCase) {
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(tokensBuilder.Build(), nil)
			},
			expectedStatus: http.StatusOK,
			expectedTokens: test.Ptr(tokensBuilder.Build()),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := authMocks.NewSignInUseCase(t)

			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Post("/sign-in", signIn(uc))

			raw, err := json.Marshal(tc.params)
			require.NoError(t, err)

			// Action
			req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader(raw))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedTokens != nil {
				var got dto.AuthTokens
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, *tc.expectedTokens, got)
			}
			if tc.expectedErrCode != "" {
				var apiErr utils.APIError
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&apiErr))
				assert.Equal(t, tc.expectedErrCode, apiErr.Code)
			}

			uc.AssertExpectations(t)
		})
	}
}

func Test_signUp(t *testing.T) {
	t.Parallel()

	paramsBuilder := fixture.AnySignUpParams().
		WithEmail("user@email.com").
		WithName("User").
		WithPassword("secret")

	userBuilder := userFixture.AnyUser().
		WithEmail("user@email.com").
		WithName("User")

	testCases := []struct {
		about           string
		params          dto.SignUpParams
		setupUC         func(uc *authMocks.SignUpUseCase)
		expectedStatus  int
		expectedClient  *user.User
		expectedErrCode string
	}{
		{
			about:  "when params are invalid",
			params: paramsBuilder.WithEmail("invalid").WithName("").WithPassword("").Build(),
			setupUC: func(uc *authMocks.SignUpUseCase) {
				params := paramsBuilder.WithEmail("invalid").WithName("").WithPassword("").Build()
				err := domainerror.New(domainerror.InvalidParams, "email: email inválido; nome: campo obrigatório; password: campo obrigatório", nil)
				uc.On("Execute", mock.Anything, params).Return(user.User{}, err)
			},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedErrCode: string(domainerror.InvalidParams),
		},
		{
			about:  "when email already exists",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignUpUseCase) {
				err := domainerror.New(domainerror.EmailAlreadyExists, "já existe um usuário com este email", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(user.User{}, err)
			},
			expectedStatus:  http.StatusConflict,
			expectedErrCode: string(domainerror.EmailAlreadyExists),
		},
		{
			about:  "when dependency error happens",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignUpUseCase) {
				err := domainerror.Wrap(errors.New("db error"), domainerror.DependecyError, "erro ao criar usuário", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(user.User{}, err)
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: string(domainerror.DependecyError),
		},
		{
			about:  "when ok",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.SignUpUseCase) {
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).
					Return(userBuilder.Build(), nil)
			},
			expectedStatus: http.StatusCreated,
			// Caso o JSON não exponha PasswordHash, manter vazio evita diferenças
			expectedClient: test.Ptr(userBuilder.WithPasswordHash("").Build()),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := authMocks.NewSignUpUseCase(t)
			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Post("/sign-up", signUp(uc))

			raw, err := json.Marshal(tc.params)
			require.NoError(t, err)

			// Action
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewReader(raw))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedClient != nil {
				var got user.User
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, tc.expectedClient.ID, got.ID)
				assert.Equal(t, tc.expectedClient.Name, got.Name)
				assert.Equal(t, tc.expectedClient.Email, got.Email)
				assert.Equal(t, tc.expectedClient.Active, got.Active)
				assert.Equal(t, tc.expectedClient.Role, got.Role)
			}
			if tc.expectedErrCode != "" {
				var apiErr utils.APIError
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&apiErr))
				assert.Equal(t, tc.expectedErrCode, apiErr.Code)
			}

			uc.AssertExpectations(t)
		})
	}
}

func Test_tokenRefresh(t *testing.T) {
	t.Parallel()

	paramsBuilder := fixture.AnyRefreshTokenParams().
		WithRefreshToken("refresh-123")

	tokensBuilder := fixture.AnyAuthTokens().
		WithAccessToken("newA").
		WithRefreshToken("newR")

	testCases := []struct {
		about           string
		params          dto.RefreshTokenParams
		setupUC         func(uc *authMocks.RefreshTokenUseCase)
		expectedStatus  int
		expectedTokens  *dto.AuthTokens
		expectedErrCode string
	}{
		{
			about:  "when params are invalid",
			params: paramsBuilder.WithRefreshToken("").Build(),
			setupUC: func(uc *authMocks.RefreshTokenUseCase) {
				params := paramsBuilder.WithRefreshToken("").Build()
				err := domainerror.New(domainerror.InvalidParams, "refreshToken: campo obrigatório", nil)
				uc.On("Execute", mock.Anything, params).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedErrCode: string(domainerror.InvalidParams),
		},
		{
			about:  "when refresh token is invalid",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.RefreshTokenUseCase) {
				err := domainerror.New(domainerror.AutenticationInvalid, "invalid refresh token", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedErrCode: string(domainerror.AutenticationInvalid),
		},
		{
			about:  "when dependency error happens",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.RefreshTokenUseCase) {
				err := domainerror.Wrap(errors.New("sign error"), domainerror.DependecyError, "erro ao gerar tokens", nil)
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).Return(dto.AuthTokens{}, err)
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: string(domainerror.DependecyError),
		},
		{
			about:  "when ok",
			params: paramsBuilder.Build(),
			setupUC: func(uc *authMocks.RefreshTokenUseCase) {
				uc.On("Execute", mock.Anything, paramsBuilder.Build()).
					Return(tokensBuilder.Build(), nil)
			},
			expectedStatus: http.StatusOK,
			expectedTokens: test.Ptr(tokensBuilder.Build()),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := authMocks.NewRefreshTokenUseCase(t)
			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Post("/token/refresh", tokenRefresh(uc))

			raw, err := json.Marshal(tc.params)
			require.NoError(t, err)

			// Action
			req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewReader(raw))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedTokens != nil {
				var got dto.AuthTokens
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, *tc.expectedTokens, got)
			}
			if tc.expectedErrCode != "" {
				var apiErr utils.APIError
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&apiErr))
				assert.Equal(t, tc.expectedErrCode, apiErr.Code)
			}

			uc.AssertExpectations(t)
		})
	}
}
