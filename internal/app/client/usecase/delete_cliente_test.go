package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	userMocks "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestDeleteClientUseCase_Execute(t *testing.T) {
	t.Parallel()

	userID := uuid.NextID()

	paramsBuilder := fixture.AnyDeleteClientParams().
		WithClientID(userID)

	userBuilder := fixtureUser.AnyUser().
		WithID(userID)

	testCases := []struct {
		about       string
		params      dto.DeleteClientParams
		setupRepo   func(r *userMocks.Repository)
		expectedErr string
	}{
		{
			about:       "when params are invalid",
			params:      paramsBuilder.WithClientID(uuid.Nil).Build(),
			expectedErr: "[AQF002] clientId: campo obrigatório",
		},
		{
			about:  "when repo delete fails",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Find", mock.Anything, userID).
					Return(userBuilder.Build(), nil)
				r.On("Delete", mock.Anything, mock.AnythingOfType("user.User")).
					Return(errors.New("db error"))
			},
			expectedErr: "[AQF004] erro ao deletar cliente | cause: db error",
		},
		{
			about:  "when find returns an not found error",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Find", mock.Anything, userID).
					Return(user.User{}, user.ErrNotFound)
			},
			expectedErr: "[AQF003] cliente não encontrado",
		},
		{
			about:  "when find returns error an generic error",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Find", mock.Anything, userID).
					Return(user.User{}, errors.New("find error"))
			},
			expectedErr: "[AQF004] erro ao buscar cliente | cause: find error",
		},
		{
			about:  "when find returns error an generic error",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Find", mock.Anything, userID).
					Return(userBuilder.Build(), nil)
				r.On("Delete", mock.Anything, userBuilder.Build()).
					Return(errors.New("db error"))
			},
			expectedErr: "[AQF004] erro ao deletar cliente | cause: db error",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Find", mock.Anything, userID).
					Return(userBuilder.Build(), nil)
				r.On("Delete", mock.Anything, userBuilder.Build()).
					Return(nil)
			},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := userMocks.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			uc := usecase.NewDeleteClientUseCase(repo)

			// Action
			err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
