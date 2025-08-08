package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/favorite"
	fixtureFavorite "github.com/uesleicarvalhoo/aiqfome/favorite/fixture"
	favMocks "github.com/uesleicarvalhoo/aiqfome/favorite/mocks"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto"
	fixtureDto "github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto/fixture"
	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/product"
	fixtureProduct "github.com/uesleicarvalhoo/aiqfome/product/fixture"
	prodMocks "github.com/uesleicarvalhoo/aiqfome/product/mocks"
)

func TestGetClientFavoritesUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()
	paramsBuilder := fixtureDto.AnyGetClientFavoritesParams().
		WithClientID(clientID)

	favoriteBuilder := fixtureFavorite.AnyFavorite().
		WithClientID(clientID)

	productBuilder := fixtureProduct.AnyProduct()

	testCases := []struct {
		about          string
		params         dto.GetClientFavoritesParams
		setupFavorites func(m *favMocks.Repository)
		setupProducts  func(m *prodMocks.Repository)
		expectedErr    string
		expectedResult dto.ClientFavorites
	}{
		{
			about:       "when params invalid",
			params:      dto.GetClientFavoritesParams{},
			expectedErr: "[AQF002] clientId: campo obrigatório",
		},
		{
			about:  "when favorites paginate fails",
			params: paramsBuilder.Build(),
			setupFavorites: func(m *favMocks.Repository) {
				m.On("PaginateByClientID", mock.Anything, clientID, 1, 20).
					Return([]favorite.Favorite{}, 0, errors.New("db error"))
			},
			expectedErr: "erro ao paginar favoritos",
		},
		{
			about:  "when getProducts returns not found",
			params: paramsBuilder.Build(),
			setupFavorites: func(m *favMocks.Repository) {
				m.On("PaginateByClientID", mock.Anything, clientID, 1, 20).
					Return([]favorite.Favorite{favoriteBuilder.Build()}, 1, nil)
			},
			setupProducts: func(m *prodMocks.Repository) {
				m.On("FindMultiple", mock.Anything, []int{1}).
					Return([]product.Product{}, &product.ErrProductsNotFound{IDs: []int{1}})
			},
			expectedErr: "produtos não encontrados",
		},
		{
			about:  "when getProducts returns other error",
			params: paramsBuilder.Build(),
			setupFavorites: func(m *favMocks.Repository) {
				m.On("PaginateByClientID", mock.Anything, clientID, 1, 20).
					Return([]favorite.Favorite{favoriteBuilder.Build()}, 1, nil)
			},
			setupProducts: func(m *prodMocks.Repository) {
				m.On("FindMultiple", mock.Anything, []int{1}).
					Return([]product.Product{}, errors.New("service down"))
			},
			expectedErr: "erro ao buscar produtos",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupFavorites: func(m *favMocks.Repository) {
				m.On("PaginateByClientID", mock.Anything, clientID, 1, 20).
					Return([]favorite.Favorite{
						favoriteBuilder.WithProductID(1).Build(),
						favoriteBuilder.WithProductID(2).Build(),
					}, 2, nil)
			},
			setupProducts: func(m *prodMocks.Repository) {
				m.On("FindMultiple", mock.Anything, []int{1, 2}).
					Return([]product.Product{
						productBuilder.WithID(1).Build(),
						productBuilder.WithID(2).Build(),
					}, nil)
			},
			expectedResult: dto.ClientFavorites{
				ClientID: clientID,
				Products: []product.Product{
					productBuilder.WithID(1).Build(),
					productBuilder.WithID(2).Build(),
				},
				Total: 2,
				Pages: 1,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// arrange
			favRepo := favMocks.NewRepository(t)
			if tc.setupFavorites != nil {
				tc.setupFavorites(favRepo)
			}

			prodRepo := prodMocks.NewRepository(t)
			if tc.setupProducts != nil {
				tc.setupProducts(prodRepo)
			}

			uc := usecase.NewGetClientFavoritesUseCase(favRepo, prodRepo)

			// act
			res, err := uc.Execute(context.Background(), tc.params)

			// assert
			if tc.expectedErr != "" {
				assert.Equal(t, dto.ClientFavorites{}, res)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			}

			favRepo.AssertExpectations(t)
			prodRepo.AssertExpectations(t)
		})
	}
}
