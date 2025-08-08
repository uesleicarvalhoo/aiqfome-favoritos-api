package role

type Resource string

const (
	ResourceMe        Resource = "me"
	ResourceClient    Resource = "client"
	ResourceFavorites Resource = "favorite"
)

func (r Resource) IsValid() bool {
	switch r {
	case ResourceClient, ResourceFavorites, ResourceMe:
		return true
	default:
		return false
	}
}
