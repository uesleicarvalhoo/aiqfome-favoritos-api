package usecase

import (
	"context"
	"fmt"

	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth"
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/logger"
	"github.com/uesleicarvalhoo/aiqfome/pkg/password"
	"github.com/uesleicarvalhoo/aiqfome/pkg/trace"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type SignUpOptions struct {
	MinPasswordLength int
}

type signUpUseCase struct {
	uuid   uuid.Generator
	hasher password.Hasher
	repo   user.Repository
	opts   SignUpOptions
}

func NewSignUpUseCase(idGen uuid.Generator, hasher password.Hasher, repo user.Repository, opts SignUpOptions) auth.SignUpUseCase {
	return &signUpUseCase{
		uuid:   idGen,
		hasher: hasher,
		repo:   repo,
		opts:   opts,
	}
}

func (u *signUpUseCase) Execute(ctx context.Context, p dto.SignUpParams) (user.User, error) {
	ctx, span := trace.NewSpan(ctx, "auth.signUp")
	defer span.End()

	if err := p.Validate(); err != nil {
		return user.User{}, err
	}

	if len(p.Password) < u.opts.MinPasswordLength {
		return user.User{}, domainerror.New(
			domainerror.InvalidParams,
			fmt.Sprintf("a senha deve ter pelo menos %d caracters", u.opts.MinPasswordLength),
			map[string]any{
				"password_length":     len(p.Password),
				"min_password_length": u.opts.MinPasswordLength,
			})
	}

	if _, err := u.repo.FindByEmail(ctx, p.Email); err == nil {
		return user.User{}, domainerror.New(domainerror.EmailAlreadyExists, "já existe um usuário com este email", map[string]any{
			"email": p.Email,
		})
	}

	id := u.uuid.NextID()

	hash, err := u.hasher.Hash(fmt.Sprintf("%s:%s", id.String(), p.Password))
	if err != nil {
		return user.User{}, domainerror.Wrap(err, domainerror.InvalidParams, "erro ao gerar o hash da senha", map[string]any{
			"error": err.Error(),
		})
	}

	c, err := user.New(
		id,
		p.Name,
		p.Email,
		hash,
		role.RoleClient,
	)
	if err != nil {
		return user.User{}, err
	}

	if err := u.repo.Create(ctx, c); err != nil {
		logger.ErrorF(ctx, "error while sign up user", logger.Fields{
			"user_email": p.Email,
			"error":      err.Error(),
		})

		return user.User{}, domainerror.Wrap(err, domainerror.DependecyError, "erro ao criar usuário", map[string]any{
			"email": p.Email,
			"name":  p.Name,
		})
	}

	return c, nil
}
