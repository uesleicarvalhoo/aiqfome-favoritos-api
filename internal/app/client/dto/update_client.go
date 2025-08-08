package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
	"github.com/uesleicarvalhoo/aiqfome/role"
)

type UpdateClientParams struct {
	ClientID uuid.ID    `json:"-"`
	Name     *string    `json:"name,omitempty"`
	Active   *bool      `json:"active,omitempty"`
	Role     *role.Role `json:"role,omitempty"`
}

func (p UpdateClientParams) Validate() error {
	v := validator.New()

	if p.ClientID.IsZero() {
		v.AddError("clientId", "campo inválido")
	}

	if p.Name != nil && *p.Name == "" {
		v.AddError("name", "campo obrigatório")
	}

	if p.Role != nil && !p.Role.IsValid() {
		v.AddError("role", "role invalida")
	}

	return v.Validate()
}
