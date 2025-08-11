package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	clientMocks "github.com/uesleicarvalhoo/aiqfome/user/mocks"
)

func TestFiendClientUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()

	userBuilder := fixtureUser.AnyUser().
		WithID(clientID).
		WithName("Client 1").
		WithEmail("client1@email.com").
		WithActive(true)

	testCases := []struct {
		about        string
		clientID     uuid.ID
		setupRepo    func(r *clientMocks.Repository)
		expectedErr  string
		expectedResp dto.Client
	}{
		{
			about:    "when repo fails",
			clientID: clientID,
			setupRepo: func(r *clientMocks.Repository) {
				r.On("Find", mock.Anything, clientID).
					Return(userBuilder.Build(), errors.New("db error"))
			},
			expectedErr: "[AQF004] error while to trying find client | cause: db error",
		},
		{
			about:    "when all is valid",
			clientID: clientID,
			setupRepo: func(r *clientMocks.Repository) {
				r.On("Find", mock.Anything, clientID).
					Return(userBuilder.Build(), nil)
			},
			expectedResp: dto.FromDomain(userBuilder.Build()),
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

			uc := usecase.NewFindClientUseCase(repo)

			// Action
			res, err := uc.Execute(context.Background(), tc.clientID)

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
