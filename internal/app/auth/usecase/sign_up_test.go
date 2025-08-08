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
	passwordMocks "github.com/uesleicarvalhoo/aiqfome/pkg/password/mocks"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	mocksUuid "github.com/uesleicarvalhoo/aiqfome/pkg/uuid/mocks"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	mocksUser "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestSignUpUseCase_Execute(t *testing.T) {
	t.Parallel()

	minLen := 8
	opts := usecase.SignUpOptions{MinPasswordLength: minLen}

	userID := uuid.NextID()
	paramsBuilder := fixtureAuth.AnySignUpParams().
		WithEmail("user@email.com").
		WithName("User LastName").
		WithPassword("secret-passwd")

	passwordToHash := fmt.Sprintf("%s:secret-passwd", userID)

	userBuilder := fixtureUser.AnyUser().
		WithEmail("user@email.com").
		WithName("User LastName").
		WithActive(true).
		WithID(userID)

	testCases := []struct {
		about        string
		params       dto.SignUpParams
		setupIDGen   func(g *mocksUuid.Generator)
		setupHasher  func(h *passwordMocks.Hasher)
		setupRepo    func(r *mocksUser.Repository)
		expectedErr  string
		expectedUser user.User
	}{
		{
			about: "when params are invalid",
			params: paramsBuilder.
				WithEmail("bad").
				WithName("").
				WithPassword("").
				Build(),
			expectedErr: "[AQF002] email: email inválido; nome: campo obrigatório; password: campo obrigatório",
		},
		{
			about:       "when password too short",
			params:      paramsBuilder.WithPassword("short").Build(),
			expectedErr: fmt.Sprintf("[AQF002] a senha deve ter pelo menos %d caracters", minLen),
		},
		{
			about:  "when email already exists",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, paramsBuilder.Build().Email).
					Return(userBuilder.Build(), nil)
			},
			expectedErr: "[USR001] já existe um usuário com este email",
		},
		{
			about:  "when hash generation fails",
			params: paramsBuilder.Build(),
			setupIDGen: func(g *mocksUuid.Generator) {
				g.On("NextID").Return(userID)
			},
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, paramsBuilder.Build().Email).
					Return(user.User{}, user.ErrNotFound)
			},
			setupHasher: func(h *passwordMocks.Hasher) {
				h.On("Hash", passwordToHash).
					Return("", errors.New("hash error"))
			},
			expectedErr: "[AQF002] erro ao gerar o hash da senha | cause: hash error",
		},
		{
			about:  "when repo.Create fails",
			params: paramsBuilder.Build(),

			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, paramsBuilder.Build().Email).
					Return(user.User{}, user.ErrNotFound)
				r.On("Create", mock.Anything, mock.AnythingOfType("user.User")).
					Return(errors.New("db error"))
			},
			setupHasher: func(h *passwordMocks.Hasher) {
				h.On("Hash", passwordToHash).
					Return("hashed", nil)
			},
			setupIDGen: func(g *mocksUuid.Generator) {
				g.On("NextID").Return(userID)
			},
			expectedErr: "[AQF004] erro ao criar usuário | cause: db error",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupRepo: func(r *mocksUser.Repository) {
				r.On("FindByEmail", mock.Anything, paramsBuilder.Build().Email).
					Return(user.User{}, user.ErrNotFound)
				r.On("Create", mock.Anything, mock.AnythingOfType("user.User")).
					Return(nil)
			},
			setupHasher: func(h *passwordMocks.Hasher) {
				h.On("Hash", passwordToHash).
					Return("hashed", nil)
			},
			setupIDGen: func(g *mocksUuid.Generator) {
				g.On("NextID").Return(userID)
			},
			expectedUser: userBuilder.WithPasswordHash("hashed").Build(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			idGen := mocksUuid.NewGenerator(t)
			if tc.setupIDGen != nil {
				tc.setupIDGen(idGen)
			}

			hasher := passwordMocks.NewHasher(t)
			if tc.setupHasher != nil {
				tc.setupHasher(hasher)
			}

			repo := mocksUser.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			uc := usecase.NewSignUpUseCase(idGen, hasher, repo, opts)

			// Action
			res, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, user.User{}, res)
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUser.ID, res.ID)
			assert.Equal(t, tc.expectedUser.Name, res.Name)
			assert.Equal(t, tc.expectedUser.Email, res.Email)
			assert.Equal(t, tc.expectedUser.Active, res.Active)
			assert.Equal(t, tc.expectedUser.PasswordHash, res.PasswordHash)
			assert.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)

			idGen.AssertExpectations(t)
			hasher.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
