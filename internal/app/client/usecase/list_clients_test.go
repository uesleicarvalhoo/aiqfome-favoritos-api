package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	fixtureDTO "github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	userMocks "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestListClientsUseCase_Execute(t *testing.T) {
	t.Parallel()

	paramsBuilder := fixtureDTO.AnyListClientsParams().
		WithPage(0).
		WithPageSize(2)

	userID := uuid.NextID()

	userBuilder := fixtureUser.AnyUser().
		WithID(userID).
		WithName("client 1").
		WithEmail("client1@email.com").
		WithActive(true)

	clientBuilder := fixture.AnyClient().
		WithID(userID).
		WithName("client 1").
		WithEmail("client1@email.com").
		WithActive(true)

	testCases := []struct {
		about        string
		params       dto.ListClientsParams
		setupRepo    func(r *userMocks.Repository)
		expectedErr  string
		expectedResp dto.PaginatedClients
	}{
		{
			about:       "when params are invalid",
			params:      paramsBuilder.WithPage(-1).Build(),
			expectedErr: "[AQF002] page: não pode ser negativo",
		},
		{
			about:  "when repo paginate fails",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Paginate", mock.Anything, 0, 2).
					Return([]user.User{}, 0, errors.New("db error"))
			},
			expectedErr: "[AQF004] erro ao paginar clientes | cause: db error",
		},
		{
			about:  "when pageSize is zero it defaults to 10",
			params: fixtureDTO.AnyListClientsParams().WithPage(0).WithPageSize(0).Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Paginate", mock.Anything, 0, 10).
					Return([]user.User{userBuilder.Build()}, 11, nil)
			},
			expectedResp: dto.PaginatedClients{
				Clients: []dto.Client{clientBuilder.Build()},
				Total:   11,
				Pages:   2,
			},
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupRepo: func(r *userMocks.Repository) {
				r.On("Paginate", mock.Anything, 0, 2).
					Return([]user.User{userBuilder.Build(), userBuilder.Build()}, 5, nil)
			},
			expectedResp: dto.PaginatedClients{
				Clients: []dto.Client{clientBuilder.Build(), clientBuilder.Build()},
				Total:   5,
				Pages:   3,
			},
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

			uc := usecase.NewListClientsUseCase(repo)

			// Action
			res, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, dto.PaginatedClients{}, res)
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp, res)
			}

			repo.AssertExpectations(t)
		})
	}
}
