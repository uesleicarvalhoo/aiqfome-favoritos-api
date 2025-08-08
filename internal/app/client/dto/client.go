package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/uuid"
	"github.com/uesleicarvalhoo/aiqfome/user"
)

type Client struct {
	ID     uuid.ID `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Active bool    `json:"active"`
}

func FromDomain(u user.User) Client {
	return Client{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Active: u.Active,
	}
}
