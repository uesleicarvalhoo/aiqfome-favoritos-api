package postgres_test

import (
	"context"
	"database/sql"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/uesleicarvalhoo/aiqfome/favorite"
	"github.com/uesleicarvalhoo/aiqfome/favorite/fixture"
	"github.com/uesleicarvalhoo/aiqfome/favorite/postgres"
	"github.com/uesleicarvalhoo/aiqfome/internal/infra/database"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/test"
	fixtureUser "github.com/uesleicarvalhoo/aiqfome/user/fixture"
	postgresUser "github.com/uesleicarvalhoo/aiqfome/user/postgres"
)

type TestSuitePostgresRepository struct {
	suite.Suite
	ctx       context.Context
	db        *sql.DB
	container *test.PostgresContainer
	repo      favorite.Repository
}

func TestFavoriteRepository(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(TestSuitePostgresRepository))
}

func (s *TestSuitePostgresRepository) SetupTest() {
	var err error

	s.ctx = context.Background()

	s.container, err = test.SetupPostgres(s.ctx)
	if err != nil {
		s.T().Fatalf("failed to setup postgres container: %s", err)
		return
	}

	s.T().Cleanup(func() {
		_ = s.container.Terminate(s.ctx)
	})

	db, err := database.NewPostgresWithMigration(
		database.Options{
			User:              s.container.Username,
			Password:          s.container.Password,
			Host:              s.container.Host,
			Port:              strconv.Itoa(s.container.Port),
			Name:              s.container.Database,
			PoolSize:          10,
			ConnMaxTTL:        0,
			TimeoutSeconds:    10,
			LockTimeoutMillis: 0,
		},
	)
	if err != nil {
		s.T().Fatalf("failed to connect to database: %s", err)
		return
	}

	s.db = db
	s.repo = postgres.NewRepository(s.db)
}

func (s *TestSuitePostgresRepository) TestCRUD() {
	productID := 1

	usr := fixtureUser.AnyUser().Build()
	require.NoError(s.T(), postgresUser.NewRepository(s.db).Create(s.ctx, usr), "failed to setup user")

	favoriteBuilder := fixture.AnyFavorite().
		WithClientID(usr.ID).
		WithProductID(productID)

	testCases := []struct {
		about       string
		setup       func()
		teardown    func()
		favorite    favorite.Favorite
		expectedErr string
	}{
		{
			about:       "when client doesn't exist in database",
			favorite:    favoriteBuilder.WithClientID(uuid.Nil).Build(),
			expectedErr: "SQLSTATE 23503",
		},
		{
			about:    "when product already is vinculated to client",
			favorite: favoriteBuilder.Build(),
			setup: func() {
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.Build()), "failed to create favorite")
			},
			expectedErr: "SQLSTATE 23505",
			teardown: func() {
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.Build()), "failed to remove favorite before create it")
			},
		},
		{
			about:    "when everything is fine",
			favorite: favoriteBuilder.Build(),
			teardown: func() {
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.Build()), "failed to remove favorite before create it")
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.about, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			if tc.teardown != nil {
				defer tc.teardown()
			}

			err := s.repo.Create(s.ctx, tc.favorite)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				found, err := s.repo.Find(s.ctx, tc.favorite.ClientID, tc.favorite.ProductID)
				require.NoError(s.T(), err, "failed to retrieve favorite")

				assert.Equal(s.T(), tc.favorite.ClientID, found.ClientID)
				assert.Equal(s.T(), tc.favorite.ProductID, found.ProductID)
			}
		})
	}
}

func (s *TestSuitePostgresRepository) TestPaginateByClientID() {
	usr := fixtureUser.AnyUser().WithEmail("user@email.com").Build()
	anotherUsr := fixtureUser.AnyUser().WithEmail("another_user@email.com").Build()
	favoriteBuilder := fixture.AnyFavorite().
		WithClientID(usr.ID)

	usrRepo := postgresUser.NewRepository(s.db)
	require.NoError(s.T(), usrRepo.Create(s.ctx, usr), "failed to setup user")
	require.NoError(s.T(), usrRepo.Create(s.ctx, anotherUsr), "failed to setup user")

	testCases := []struct {
		about             string
		setup             func()
		teardown          func()
		clientID          uuid.ID
		page              int
		pageSize          int
		expectedTotal     int
		expectedFavorites []favorite.Favorite
		expectedErr       string
	}{
		{
			about:             "when client doesn't exist in database",
			clientID:          uuid.NextID(),
			expectedTotal:     0,
			expectedFavorites: []favorite.Favorite{},
			page:              0,
			pageSize:          10,
		},
		{
			about:             "when there are no products vinculated to current client",
			clientID:          usr.ID,
			page:              0,
			pageSize:          10,
			expectedTotal:     0,
			expectedFavorites: []favorite.Favorite{},
			setup: func() {
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.WithClientID(anotherUsr.ID).Build()), "failed to create favorite")
			},
			teardown: func() {
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.WithClientID(anotherUsr.ID).Build()), "failed to remove favorite before create it")
			},
		},
		{
			about:         "when everything is fine",
			clientID:      usr.ID,
			page:          0,
			pageSize:      2,
			expectedTotal: 4,
			expectedFavorites: []favorite.Favorite{
				favoriteBuilder.WithProductID(1).Build(),
				favoriteBuilder.WithProductID(2).Build(),
			},
			setup: func() {
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.WithProductID(1).Build()), "failed to create favorite")
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.WithProductID(2).Build()), "failed to create favorite")
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.WithProductID(3).Build()), "failed to create favorite")
				require.NoError(s.T(), s.repo.Create(s.ctx, favoriteBuilder.WithProductID(4).Build()), "failed to create favorite")
			},
			teardown: func() {
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.WithProductID(1).Build()), "failed to remove favorite before create it")
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.WithProductID(2).Build()), "failed to remove favorite before create it")
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.WithProductID(3).Build()), "failed to remove favorite before create it")
				require.NoError(s.T(), s.repo.Remove(s.ctx, favoriteBuilder.WithProductID(4).Build()), "failed to remove favorite before create it")
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.about, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			if tc.teardown != nil {
				defer tc.teardown()
			}

			found, total, err := s.repo.PaginateByClientID(s.ctx, tc.clientID, tc.page, tc.pageSize)

			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.Empty(t, found)
				assert.Equal(t, total, 0)
			} else {
				assert.NoError(t, err)
				assert.Len(t, found, len(tc.expectedFavorites))

				for _, f := range tc.expectedFavorites {
					ok := slices.ContainsFunc(found, func(fv favorite.Favorite) bool {
						return fv.ClientID == f.ClientID && fv.ProductID == f.ProductID
					})

					assert.True(t, ok, "favorite '%+v' not found", f)
				}
				assert.Equal(t, tc.expectedTotal, total)
			}
		})
	}
}
