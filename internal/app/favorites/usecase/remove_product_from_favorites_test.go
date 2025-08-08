package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	mocksFavorite "github.com/uesleicarvalhoo/aiqfome/favorite/mocks"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	fixtureDto "github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

func TestRemoveProductFromFavoritesUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()
	productID := 1

	paramsBuilder := fixtureDto.AnyRemoveProductFromFavoritesParams().
		WithClientID(clientID).
		WithProductID(productID)

	testCases := []struct {
		about       string
		params      dto.RemoveProductFromFavoritesParams
		setupRepo   func(m *mocksFavorite.Repository)
		expectedErr string
	}{
		{
			about:       "when params are invalid",
			params:      dto.RemoveProductFromFavoritesParams{},
			setupRepo:   nil,
			expectedErr: "clientId: campo obrigatório; productId: campo obrigatório",
		},
		{
			about:  "when favorite not found",
			params: paramsBuilder.Build(),
			setupRepo: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{}, &favorite.ErrFavoriteNotFound{ClientID: clientID, ProductID: productID})
			},
			expectedErr: "[AQF003] favorito não encontrado",
		},
		{
			about:  "when find returns other error",
			params: paramsBuilder.Build(),
			setupRepo: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{}, errors.New("db find error"))
			},
			expectedErr: "[AQF004] error while to trying find favorite",
		},
		{
			about:  "when remove fails",
			params: paramsBuilder.Build(),
			setupRepo: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{ClientID: clientID, ProductID: productID}, nil)
				m.On("Remove", mock.Anything, favorite.Favorite{ClientID: clientID, ProductID: productID}).
					Return(errors.New("db remove error"))
			},
			expectedErr: "[AQF004] error while trying to remove favorite",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupRepo: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{ClientID: clientID, ProductID: productID}, nil)
				m.On("Remove", mock.Anything, favorite.Favorite{ClientID: clientID, ProductID: productID}).
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
			repo := mocksFavorite.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			uc := usecase.NewRemoveProductFromFavoritesUseCase(repo)

			// Action
			err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)

				return
			}

			assert.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}
}
