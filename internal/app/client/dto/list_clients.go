package dto

import (
	"github.com/uesleicarvalhoo/aiqfome/pkg/validator"
)

type ListClientsParams struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p ListClientsParams) Validate() error {
	v := validator.New()

	if p.PageSize < 1 {
		v.AddError("pageSize", "deve ser maior do que 1")
	}

	if p.Page < 0 {
		v.AddError("page", "não pode ser negativo")
	}

	return v.Validate()
}

type PaginatedClients struct {
	Clients []Client `json:"clients"`
	Total   int      `json:"total"`
	Pages   int      `json:"pages"`
}
