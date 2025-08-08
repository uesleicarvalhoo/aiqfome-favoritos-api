package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto/fixture"
	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/auth/usecase"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/role/mocks"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
)

func TestAuthorizeUseCase_Execute(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()
	clientRole := role.RoleClient
	resource := role.ResourceClient

	// Parâmetros básicos
	clientBuilder := fixtureUser.AnyUser().
		WithRole(role.Role(clientRole)).
		WithID(clientID)

	paramsBuilder := fixture.AnyAuthorizeParams().
		WithUser(clientBuilder.Build()).
		WithAction(role.ActionRead).
		WithResource(role.ResourceClient)

	testCases := []struct {
		about         string
		setupRepo     func(repo *mocks.Repository)
		params        dto.AuthorizeParams
		expectedError string
	}{
		{
			about:  "when role not found in repo",
			params: paramsBuilder.Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return(nil, &role.ErrNotFound{Name: string(clientRole)})
			},
			expectedError: "[AQF003] role não encontrada",
		},
		{
			about:  "when repo returns generic error",
			params: paramsBuilder.Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return(nil, errors.New("i'm a repository error"))
			},
			expectedError: "[AQF004] ocorreu um erro ao obter as permissões | cause: i'm a repository error",
		},
		{
			about: "when role is admin, should by pass permissions",
			params: paramsBuilder.WithUser(
				clientBuilder.
					WithRole(role.RoleAdmin).
					Build(),
			).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, role.RoleAdmin).
					Return([]role.Permission{}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when has manage permission, and require manage",
			params: paramsBuilder.WithAction(role.ActionManage).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: resource, Action: role.ActionManage},
					}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when has manage permission, and require read",
			params: paramsBuilder.WithAction(role.ActionRead).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: resource, Action: role.ActionManage},
					}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when has manage permission, and require write",
			params: paramsBuilder.WithAction(role.ActionWrite).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: resource, Action: role.ActionManage},
					}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when has exact permission",
			params: paramsBuilder.WithAction(role.ActionRead).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: resource, Action: role.ActionRead},
					}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when has permission to write and require read",
			params: paramsBuilder.WithAction(role.ActionRead).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: resource, Action: role.ActionWrite},
					}, nil)
			},
			expectedError: "",
		},
		{
			about:  "when no permission matches",
			params: paramsBuilder.WithUser(clientBuilder.Build()).WithResource(role.ResourceFavorites).Build(),
			setupRepo: func(repo *mocks.Repository) {
				repo.On("FindPermissions", mock.Anything, clientRole).
					Return([]role.Permission{
						{Resource: role.ResourceClient, Action: role.ActionManage},
					}, nil)
			},
			expectedError: "[AQF005] permissão negada",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo := mocks.NewRepository(t)
			if tc.setupRepo != nil {
				tc.setupRepo(repo)
			}

			uc := usecase.NewAuthorizeUseCase(repo)

			// Action
			err := uc.Execute(context.Background(), tc.params)

			// Assert
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
