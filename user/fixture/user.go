package fixture

import (
	"time"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/role"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type UserBuilder struct {
	id           uuid.ID
	name         string
	email        string
	passwordHash string
	active       bool
	createdAt    time.Time
	role         role.Role
}

func AnyUser() UserBuilder {
	return UserBuilder{
		id:           uuid.NextID(),
		name:         "Ueslei Carvalho",
		email:        "user@email.com",
		passwordHash: "hashed_password",
		active:       true,
		role:         role.RoleAdmin,
		createdAt:    time.Now(),
	}
}

func (b UserBuilder) WithID(id uuid.ID) UserBuilder {
	b.id = id
	return b
}

func (b UserBuilder) WithName(name string) UserBuilder {
	b.name = name
	return b
}

func (b UserBuilder) WithEmail(email string) UserBuilder {
	b.email = email
	return b
}

func (b UserBuilder) WithPasswordHash(hash string) UserBuilder {
	b.passwordHash = hash
	return b
}

func (b UserBuilder) WithActive(active bool) UserBuilder {
	b.active = active
	return b
}

func (b UserBuilder) WithCreatedAt(t time.Time) UserBuilder {
	b.createdAt = t
	return b
}

func (b UserBuilder) WithRole(r role.Role) UserBuilder {
	b.role = r
	return b
}

func (b UserBuilder) Build() user.User {
	return user.User{
		ID:           b.id,
		Name:         b.name,
		Email:        b.email,
		PasswordHash: b.passwordHash,
		Active:       b.active,
		CreatedAt:    b.createdAt,
		Role:         b.role,
	}
}
