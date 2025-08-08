package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type AuthorizeParams struct {
	User     user.User
	Resource role.Resource
	Action   role.Action
}
