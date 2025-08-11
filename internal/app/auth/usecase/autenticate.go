package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/pkg/cache"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type authenticateUseCase struct {
	repo         user.Repository
	access       jwt.Provider
	cache        cache.Cache
	cacheExpTime time.Duration
}

func NewAuthenticateUseCase(repo user.Repository, access jwt.Provider, cache cache.Cache, cacheExpTime time.Duration) auth.AuthenticateUseCase {
	return &authenticateUseCase{
		repo:         repo,
		access:       access,
		cache:        cache,
		cacheExpTime: cacheExpTime,
	}
}

func (u *authenticateUseCase) Execute(ctx context.Context, token string) (user.User, error) {
	ctx, span := trace.NewSpan(ctx, "auth.authenticate")
	defer span.End()

	cl, err := u.access.Validate(ctx, token)
	if err != nil {
		logger.InfoF(ctx, "invalid access token", logger.Fields{
			"error": err.Error(),
		})

		return user.User{}, err
	}

	usr, err := u.getUser(ctx, cl.UserID)
	if err != nil {
		return user.User{}, err
	}

	if !usr.Active {
		logger.WarnF(ctx, "usuário bloqueado", logger.Fields{
			"client_id": usr.ID,
		})
		return user.User{}, domainerror.New(domainerror.UserNotActive, "usuário bloqueado", map[string]any{
			"client_id": usr.ID,
		})
	}

	return usr, nil
}

func (u *authenticateUseCase) getUser(ctx context.Context, id uuid.ID) (user.User, error) {
	key := fmt.Sprintf("user:%s", id.String())

	data, err := u.cache.Get(ctx, key)
	if err == nil {
		if data != nil {
			var usr user.User
			err := json.Unmarshal(data, &usr)
			if err == nil {
				logger.DebugF(ctx, "read user from cache", logger.Fields{
					"user_id": id,
				})
				return usr, nil
			}
			logger.ErrorF(ctx, "failed to unmarshal user from cache", logger.Fields{
				"user_id": id,
				"error":   err.Error(),
			})

		}
	} else {
		logger.ErrorF(ctx, "falied to read user from cache", logger.Fields{
			"user_id": id,
			"error":   err.Error(),
		})
	}

	usr, err := u.repo.Find(ctx, id)
	if err != nil {
		logger.InfoF(ctx, "error while trying to find user", logger.Fields{
			"user_id": id,
			"error":   err.Error(),
		})

		if errors.Is(err, user.ErrNotFound) {
			return user.User{}, domainerror.Wrap(err, domainerror.ResourceNotFound, "user not found", logger.Fields{
				"user_id": id,
			})
		}

		return user.User{}, domainerror.Wrap(err, domainerror.DependecyError, "error while trying to find user", logger.Fields{
			"user_id": id,
			"error":   err.Error(),
		})
	}

	go func(usr user.User) {
		v, err := json.Marshal(usr)
		if err != nil {
			logger.ErrorF(ctx, "failed to unmarshal user", logger.Fields{
				"user_id": id,
				"error":   err.Error(),
			})
			return
		}

		if err := u.cache.Set(ctx, key, v, u.cacheExpTime); err != nil {
			logger.ErrorF(ctx, "failed to save user on cache", logger.Fields{
				"user_id": id,
				"error":   err.Error(),
			})
		}
	}(usr)

	return usr, nil
}
