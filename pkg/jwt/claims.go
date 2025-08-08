package jwt

import "github.com/uesleicarvalhoo/aiqfome/pkg/uuid"

type Claims struct {
	UserID uuid.ID `json:"userId"`
}
