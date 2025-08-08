package fixture

import (
	"github.com/uesleicarvalhoo/aiqfome/internal/app/auth/dto"
	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/user"
	userFixture "github.com/uesleicarvalhoo/aiqfome/user/fixture"
)

type AuthorizeParamsBuilder struct {
	user     user.User
	resource role.Resource
	action   role.Action
}

func AnyAuthorizeParams() AuthorizeParamsBuilder {
	return AuthorizeParamsBuilder{
		user:     userFixture.AnyUser().Build(),
		resource: "default_resource",
		action:   role.ActionRead,
	}
}

func (b AuthorizeParamsBuilder) WithUser(u user.User) AuthorizeParamsBuilder {
	b.user = u
	return b
}

func (b AuthorizeParamsBuilder) WithResource(r role.Resource) AuthorizeParamsBuilder {
	b.resource = r
	return b
}

func (b AuthorizeParamsBuilder) WithAction(a role.Action) AuthorizeParamsBuilder {
	b.action = a
	return b
}

func (b AuthorizeParamsBuilder) Build() dto.AuthorizeParams {
	return dto.AuthorizeParams{
		User:     b.user,
		Resource: b.resource,
		Action:   b.action,
	}
}
