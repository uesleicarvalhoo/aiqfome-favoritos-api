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
	fixtureFavorites "github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/dto/fixture"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/favorites/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/product"
	fixtureProd "github.com/uesleicarvalhoo/aiqfome/product/fixture"
	mocksProduct "github.com/uesleicarvalhoo/aiqfome/product/mocks"
)

func TestAddProductToFavoritesUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()
	productID := 1
	paramsBuilder := fixtureFavorites.AnyAddProductToFavoritesParams().
		WithClientID(clientID)

	productBuilder := fixtureProd.AnyProduct().WithID(productID)

	testCases := []struct {
		about          string
		params         dto.AddProductToFavoritesParams
		setupProducts  func(m *mocksProduct.Reader)
		setupFavorites func(m *mocksFavorite.Repository)
		expectedErr    string
		expectedResult dto.ProductFavorite
	}{
		{
			about:       "when params are invalid",
			params:      dto.AddProductToFavoritesParams{},
			expectedErr: "[AQF002] clientId: campo obrigatório; productId: campo obrigatório",
		},
		{
			about:  "when product not found",
			params: paramsBuilder.Build(),
			setupProducts: func(m *mocksProduct.Reader) {
				m.On("Find", mock.Anything, productID).
					Return(product.Product{}, &product.ErrNotFound{ID: productID})
			},
			expectedErr: "[AQF003] produto não encontrado",
		},
		{
			about:  "when product reader returns other error",
			params: paramsBuilder.Build(),
			setupProducts: func(m *mocksProduct.Reader) {
				m.On("Find", mock.Anything, productID).
					Return(product.Product{}, errors.New("service down"))
			},
			expectedErr: "[AQF004] erro ao obter dados do produto",
		},
		{
			about:  "when already favorite",
			params: paramsBuilder.Build(),
			setupProducts: func(m *mocksProduct.Reader) {
				m.On("Find", mock.Anything, productID).
					Return(productBuilder.Build(), nil)
			},
			setupFavorites: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{}, nil)
			},
			expectedErr: "[FAV001] o produto já está nos favoritos",
		},
		{
			about:  "when favorites repository Create fails",
			params: paramsBuilder.Build(),
			setupProducts: func(m *mocksProduct.Reader) {
				m.On("Find", mock.Anything, productID).
					Return(productBuilder.Build(), nil)
			},
			setupFavorites: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{}, errors.New("db error"))
				m.On("Create", mock.Anything, mock.AnythingOfType("favorite.Favorite")).
					Return(errors.New("db error"))
			},
			expectedErr: "[AQF004] erro ao adicionar o produto aos favoritos",
		},
		{
			about:  "when all is valid",
			params: paramsBuilder.Build(),
			setupProducts: func(m *mocksProduct.Reader) {
				m.On("Find", mock.Anything, productID).
					Return(productBuilder.Build(), nil)
			},
			setupFavorites: func(m *mocksFavorite.Repository) {
				m.On("Find", mock.Anything, clientID, productID).
					Return(favorite.Favorite{}, errors.New("not found"))
				m.On("Create", mock.Anything, mock.MatchedBy(func(f favorite.Favorite) bool {
					return f.ClientID == clientID && f.ProductID == productID
				})).Return(nil)
			},
			expectedResult: dto.ProductFavorite{ClientID: clientID, Product: productBuilder.Build()},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			prodReader := mocksProduct.NewReader(t)
			if tc.setupProducts != nil {
				tc.setupProducts(prodReader)
			}

			favRepo := mocksFavorite.NewRepository(t)
			if tc.setupFavorites != nil {
				tc.setupFavorites(favRepo)
			}

			uc := usecase.NewAddProductToFavoritesUseCase(prodReader, favRepo)

			// Action
			res, err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedErr != "" {
				assert.Equal(t, dto.ProductFavorite{}, res)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			}

			prodReader.AssertExpectations(t)
			favRepo.AssertExpectations(t)
		})
	}
}
