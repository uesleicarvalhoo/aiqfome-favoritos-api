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

	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	clientDTO "github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/client/dto/fixture"
	clientMocks "github.com/uesleicarvalhoo/aiqfome/internal/app/client/mocks"
	"github.com/uesleicarvalhoo/aiqfome/internal/http/utils"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/test"
)

func Test_listClients(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		about           string
		query           string
		setupUC         func(uc *clientMocks.ListClientsUseCase)
		expectedStatus  int
		expectedBody    *clientDTO.PaginatedClients
		expectedErrCode string
	}{
		{
			about: "when ok and pageSize default to 10",
			query: "?page=0&pageSize=0",
			setupUC: func(uc *clientMocks.ListClientsUseCase) {
				resp := clientDTO.PaginatedClients{
					Clients: []dto.Client{
						fixture.AnyClient().Build(),
						fixture.AnyClient().Build(),
					},
					Total: 2,
					Pages: 1,
				}
				uc.
					On("Execute", mock.Anything, clientDTO.ListClientsParams{Page: 0, PageSize: 10}).
					Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &clientDTO.PaginatedClients{
				Clients: []dto.Client{
					fixture.AnyClient().Build(),
					fixture.AnyClient().Build(),
				},
				Total: 2,
				Pages: 1,
			},
		},
		{
			about: "when usecase returns dependency error",
			query: "?page=0&pageSize=10",
			setupUC: func(uc *clientMocks.ListClientsUseCase) {
				err := domainerror.Wrap(errors.New("db error"), domainerror.DependecyError, "erro ao paginar clientes", nil)
				uc.
					On("Execute", mock.Anything, clientDTO.ListClientsParams{Page: 0, PageSize: 10}).
					Return(clientDTO.PaginatedClients{}, err)
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: string(domainerror.DependecyError),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := clientMocks.NewListClientsUseCase(t)
			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Get("/", listClients(uc))

			// Action
			req := httptest.NewRequest(http.MethodGet, "/"+tc.query, nil)
			resp, err := app.Test(req)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedBody != nil {
				var got clientDTO.PaginatedClients
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
				assert.Equal(t, tc.expectedBody.Total, got.Total)
				assert.Equal(t, tc.expectedBody.Pages, got.Pages)
				assert.Len(t, got.Clients, len(tc.expectedBody.Clients))
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

func Test_updateClient(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()
	newName := "New name"

	clientBuilder := fixture.AnyClient().
		WithID(clientID).
		WithName("old name")

	testCases := []struct {
		about           string
		id              string
		body            any
		setupUC         func(uc *clientMocks.UpdateClientUseCase)
		expectedStatus  int
		expectedClient  *dto.Client
		expectedErrCode string
	}{
		{
			about:           "when id is invalid uuid",
			id:              "not-a-uuid",
			body:            map[string]any{"name": newName},
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedErrCode: string(domainerror.InvalidParams),
		},
		{
			about: "when usecase returns not found",
			id:    clientID.String(),
			body:  map[string]any{"name": newName},
			setupUC: func(uc *clientMocks.UpdateClientUseCase) {
				err := domainerror.New(domainerror.ResourceNotFound, "cliente não encontrado", nil)
				uc.
					On("Execute", mock.Anything, mock.MatchedBy(func(p clientDTO.UpdateClientParams) bool {
						return p.ClientID == clientID && p.Name != nil && *p.Name == newName
					})).
					Return(dto.Client{}, err)
			},
			expectedStatus:  http.StatusNotFound,
			expectedErrCode: string(domainerror.ResourceNotFound),
		},
		{
			about: "when ok",
			id:    clientID.String(),
			body:  map[string]any{"name": newName},
			setupUC: func(uc *clientMocks.UpdateClientUseCase) {
				uc.
					On("Execute", mock.Anything, mock.MatchedBy(func(p clientDTO.UpdateClientParams) bool {
						return p.ClientID == clientID && p.Name != nil && *p.Name == newName
					})).
					Return(clientBuilder.WithName(newName).Build(), nil)
			},
			expectedStatus: http.StatusOK,
			expectedClient: test.Ptr(clientBuilder.WithName(newName).Build()),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := clientMocks.NewUpdateClientUseCase(t)
			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Post("/:id", updateClient(uc))

			var body []byte
			if tc.body != nil {
				body, _ = json.Marshal(tc.body)
			}

			// Action
			req := httptest.NewRequest(http.MethodPost, "/"+tc.id, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Assert
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedClient != nil {
				var c dto.Client
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&c))
				assert.Equal(t, tc.expectedClient.ID, c.ID)
				assert.Equal(t, tc.expectedClient.Name, c.Name)
				assert.Equal(t, tc.expectedClient.Email, c.Email)
				assert.Equal(t, tc.expectedClient.Active, c.Active)
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

func Test_deleteClient(t *testing.T) {
	t.Parallel()

	clientID := uuid.NextID()

	testCases := []struct {
		about           string
		id              string
		setupUC         func(uc *clientMocks.DeleteClientUseCase)
		expectedStatus  int
		expectedErrCode string
	}{
		{
			about:           "when id is invalid uuid",
			id:              "invalid",
			expectedStatus:  http.StatusUnprocessableEntity,
			expectedErrCode: string(domainerror.InvalidParams),
		},
		{
			about: "when usecase returns dependency error",
			id:    clientID.String(),
			setupUC: func(uc *clientMocks.DeleteClientUseCase) {
				err := domainerror.Wrap(errors.New("db error"), domainerror.DependecyError, "erro ao deletar cliente", nil)
				uc.
					On("Execute", mock.Anything, clientDTO.DeleteClientParams{ClientID: clientID}).
					Return(err)
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedErrCode: string(domainerror.DependecyError),
		},
		{
			about: "when ok",
			id:    clientID.String(),
			setupUC: func(uc *clientMocks.DeleteClientUseCase) {
				uc.
					On("Execute", mock.Anything, clientDTO.DeleteClientParams{ClientID: clientID}).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.about, func(t *testing.T) {
			t.Parallel()

			// Arrange
			uc := clientMocks.NewDeleteClientUseCase(t)
			if tc.setupUC != nil {
				tc.setupUC(uc)
			}

			app := fiber.New()
			app.Delete("/:id", deleteClient(uc))

			// Action
			req := httptest.NewRequest(http.MethodDelete, "/"+tc.id, nil)

			resp, err := app.Test(req)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedErrCode != "" {
				var apiErr utils.APIError
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&apiErr))
				assert.Equal(t, tc.expectedErrCode, apiErr.Code)
			}

			uc.AssertExpectations(t)
		})
	}
}
