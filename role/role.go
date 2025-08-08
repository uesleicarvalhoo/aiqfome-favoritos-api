package role

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleClient Role = "client"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleClient:
		return true
	}

	return false
}
