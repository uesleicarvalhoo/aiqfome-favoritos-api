package role

import "github.com/uesleicarvalhoo/aiqfome/pkg/validator"

type Permission struct {
	Resource Resource `json:"resource"`
	Action   Action   `json:"action"`
}

func (p Permission) String() string {
	return string(p.Resource) + ":" + string(p.Action)
}

func (p Permission) Validate() error {
	v := validator.New()

	if !p.Resource.IsValid() {
		v.AddError("resource", "resource inválido")
	}

	if !p.Action.IsValid() {
		v.AddError("action", "action inválido")
	}

	return v.Validate()
}

func NewPermission(r Resource, a Action) (Permission, error) {
	p := Permission{
		Resource: r,
		Action:   a,
	}

	if err := p.Validate(); err != nil {
		return Permission{}, err
	}

	return p, nil
}
