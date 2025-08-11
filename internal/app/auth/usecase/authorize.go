package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	usecase "github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/cache"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

type authorizeUseCase struct {
	repo         role.Repository
	cache        cache.Cache
	cacheExpTime time.Duration
}

func NewAuthorizeUseCase(repo role.Repository, cache cache.Cache, cacheExpTime time.Duration) usecase.AuthorizeUseCase {
	return &authorizeUseCase{
		repo:         repo,
		cache:        cache,
		cacheExpTime: cacheExpTime,
	}
}

func (u *authorizeUseCase) Execute(ctx context.Context, params dto.AuthorizeParams) error {
	ctx, span := trace.NewSpan(ctx, "auth.authorize")
	defer span.End()

	pp, err := u.getRolePermissions(ctx, params.User.Role)
	if err != nil {
		return err
	}

	if params.User.Role == role.RoleAdmin {
		logger.InfoF(ctx, "access grant by role admin", logger.Fields{
			"client_id":   params.User.ID,
			"client_role": params.User.Role,
		})

		return nil
	}

	for _, p := range pp {
		if p.Resource != params.Resource {
			continue
		}

		switch {
		case p.Action == role.ActionManage:
			logger.InfoF(ctx, "authorization granted by manage", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})

			return nil

		case p.Action == params.Action:
			logger.InfoF(ctx, "authorization granded", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})

			return nil

		case p.Action == role.ActionWrite && params.Action == role.ActionRead:
			logger.InfoF(ctx, "authorization granded by write > read", logger.Fields{
				"client_id": params.User.ID,
				"role":      params.User.Role,
				"resource":  params.Resource,
				"action":    params.Action,
			})
			return nil
		}
	}

	logger.ErrorF(ctx, "permission denied", logger.Fields{
		"client_id":   params.User.ID,
		"role":        params.User.Role,
		"resource":    params.Resource,
		"action":      params.Action,
		"permissions": pp,
	})

	return domainerror.New(
		domainerror.OperationNotAllowed,
		"permissão negada",
		map[string]any{
			"client_id":   params.User.ID,
			"role":        params.User.Role,
			"resource":    params.Resource,
			"action":      params.Action,
			"permissions": pp,
		},
	)
}

func (u *authorizeUseCase) getRolePermissions(ctx context.Context, rl role.Role) ([]role.Permission, error) {
	key := fmt.Sprintf("role-permissions:%s", rl)
	data, err := u.cache.Get(ctx, key)
	if err == nil && data != nil {
		var pp []role.Permission
		if err := json.Unmarshal(data, &pp); err == nil {
			logger.DebugF(ctx, "read permissions from cache", logger.Fields{
				"role":        rl,
				"permissions": pp,
			})
			return pp, nil
		} else {
			logger.ErrorF(ctx, "failed to unmarshal roles from cache", logger.Fields{
				"role":  rl,
				"error": err.Error(),
			})
		}
	}

	if err != nil {
		logger.ErrorF(ctx, "falied to read role permissions from cache", logger.Fields{
			"role":  rl,
			"error": err.Error(),
		})
	}

	pp, err := u.repo.FindPermissions(ctx, rl)
	if err != nil {
		if _, ok := err.(*role.ErrNotFound); ok {
			logger.ErrorF(ctx, "role not found", logger.Fields{
				"role": rl,
			})

			return nil, domainerror.New(domainerror.ResourceNotFound, "role não encontrada", map[string]any{"role": rl})
		}

		return nil, domainerror.Wrap(err, domainerror.DependecyError, "ocorreu um erro ao obter as permissões", map[string]any{"role": rl, "error": err.Error()})
	}

	go func() {
		v, err := json.Marshal(pp)
		if err != nil {
			logger.ErrorF(ctx, "failed to marshal role permissions", logger.Fields{
				"role":        rl,
				"permissions": pp,
			})
			return
		}

		if err := u.cache.Set(ctx, key, v, u.cacheExpTime); err != nil {
			logger.ErrorF(ctx, "failed to save role permissions on cache", logger.Fields{
				"role":        rl,
				"permissions": pp,
			})
		}
	}()

	return pp, nil
}
