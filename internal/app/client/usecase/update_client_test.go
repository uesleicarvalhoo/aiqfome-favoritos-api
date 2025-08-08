package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	fixtureDTO "github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	clientMocks "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestUpdateClientUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()

	paramsBuilder := fixtureDTO.AnyUpdateClientParams().
		WithClientID(clientID)

	userBuilder := fixtureUser.AnyUser().
		WithID(clientID).
		WithName("Client 1").
		WithEmail("client1@email.com").
		WithActive(true)

	clientBuilder := fixtureDTO.AnyClient().
		WithID(clientID).
		WithName("Client 1").
		WithEmail("client1@email.com").
		WithActive(true)

	// values for update
	newName := "New Name"
	newActive := false

	testCases := []struct {
		about        string
		params       dto.UpdateClientParams
		setupRepo    func(r *clientMocks.Repository)
		expectedErr  string
		expectedResp dto.Client
	}{
		{
			about:       "when params are invalid",
			params:      paramsBuilder.WithClientID(uuid.Nil).Build(),
			expectedErr: "[AQF002] clientId: campo inválido",
		},
		{
			about:  "when repo update fails",
			params: paramsBuilder.WithName(newName).WithActive(newActive).Build(),
			setupRepo: func(r *clientMocks.Repository) {
				r.On("Find", mock.Anything, clientID).
					Return(userBuilder.Build(), nil)

				r.On("Update", mock.Anything, mock.MatchedBy(func(c user.User) bool {
					return c.ID == clientID && c.Name == newName && c.Active == newActive
				})).
					Return(errors.New("db error"))
			},
			expectedErr: "[AQF004] failed to update client | cause: db error",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.WithName(newName).WithActive(newActive).Build(),
			setupRepo: func(r *clientMocks.Repository) {
				r.On("Find", mock.Anything, clientID).
					Return(userBuilder.Build(), nil)

				r.On("Update", mock.Anything, mock.MatchedBy(func(c user.User) bool {
					return c.ID == clientID && c.Name == newName && c.Active == newActive
				})).
					Return(nil)
			},
			expectedResp: clientBuilder.WithName(newName).WithActive(newActive).Build(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := clientMocks.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			uc := usecase.NewUpdateClientUseCase(repo)

			// Action
			res, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, dto.Client{}, res)
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp, res)
			}

			repo.AssertExpectations(t)
		})
	}
}
