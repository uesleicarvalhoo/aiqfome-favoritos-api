package user

import (
	"strings"
	"time"

	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

func New(id uuid.ID, name, email, pswdHash string, role role.Role) (User, error) {
	c := User{
		ID:           id,
		Name:         strings.TrimSpace(name),
		Email:        strings.TrimSpace(email),
		PasswordHash: pswdHash,
		Active:       true,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := c.Validate(); err != nil {
		return User{}, err
	}

	return c, nil
}

type User struct {
	ID           uuid.ID   `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	Role         role.Role `json:"role"`
}

func (c User) Validate() error {
	v := validator.New()

	if c.ID.IsZero() {
		v.AddError("id", "campo obrigatório")
	}

	if c.Name == "" {
		v.AddError("name", "campo obrigatório")
	}

	if c.PasswordHash == "" {
		v.AddError("password", "campo obrigatório")
	}

	if !c.Role.IsValid() {
		v.AddError("role", "role invalida")
	}

	if c.Email == "" {
		v.AddError("email", "campo obrigatório")
	} else if !validator.IsEmailValid(c.Email) {
		v.AddError("email", "email invalido")
	}

	return v.Validate()
}
