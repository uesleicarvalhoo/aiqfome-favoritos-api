package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/uesleicarvalhoo/aiqfome/pkg/domainerror"
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
)

type Options struct {
	Issuer    string
	Audiencer string
	Secret    string
}

type Provider interface {
	Generate(ctx context.Context, sub string, d time.Duration) (string, error)
	Validate(ctx context.Context, tokenStr string) (Claims, error)
}

type provider struct {
	opts Options
}

func NewProvider(opts Options) Provider {
	return &provider{
		opts: opts,
	}
}

func (p *provider) Generate(_ context.Context, sub string, d time.Duration) (string, error) {
	exp := time.Now().Add(d)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": p.opts.Issuer,
		"sub": sub,
		"aud": p.opts.Audiencer,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": exp.Unix(),
	})

	return claims.SignedString([]byte(p.opts.Secret))
}

func (p *provider) Validate(_ context.Context, tokenStr string) (Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domainerror.New(domainerror.AutenticationInvalid, "método de assinatura inválido", nil)
		}

		return []byte(p.opts.Secret), nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, domainerror.New(domainerror.AutenticationInvalid, "token inválido", nil)
	}

	if iss, _ := claims["iss"].(string); iss != p.opts.Issuer {
		return Claims{}, domainerror.New(domainerror.AutenticationInvalid, "issuer inválido", nil)
	}

	if aud, _ := claims["aud"].(string); aud != p.opts.Audiencer {
		return Claims{}, domainerror.New(domainerror.AutenticationInvalid, "audience inválido", nil)
	}

	sub, okSub := claims["sub"].(string)
	if !okSub {
		return Claims{}, domainerror.New(domainerror.AutenticationInvalid, "claims faltando sub", nil)
	}

	clientID, err := uuid.ParseID(sub)
	if err != nil {
		return Claims{}, domainerror.New(domainerror.AutenticationInvalid, "ID do usuário inválido no token", nil)
	}

	return Claims{
		UserID: clientID,
	}, nil
}
