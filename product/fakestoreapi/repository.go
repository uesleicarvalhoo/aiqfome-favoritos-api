package fakestoreapi

import (
	"context"
	"encoding/json"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/uesleicarvalhoo/aiqfome/internal/infra/requester"
	"github.com/uesleicarvalhoo/aiqfome/product"
)

type Options struct {
	BaseUrl         string
	GetAllEndpoint  string
	GetByIdEndpoint string
}

type repository struct {
	baseUrl         string
	requester       requester.Requester
	getAllEndpoint  string
	getByIdEndpoint string
}

func NewRepository(rq requester.Requester, opts Options) product.Repository {
	return &repository{
		baseUrl:         opts.BaseUrl,
		getAllEndpoint:  opts.GetAllEndpoint,
		getByIdEndpoint: opts.GetByIdEndpoint,
		requester:       rq,
	}
}

func (r *repository) Find(ctx context.Context, id int) (product.Product, error) {
	endpoint, err := url.JoinPath(r.baseUrl, strings.ReplaceAll(r.getByIdEndpoint, "{id}", strconv.Itoa(id)))
	if err != nil {
		return product.Product{}, err
	}

	res, _, err := r.requester.Get(ctx, endpoint, nil)
	if err != nil {
		return product.Product{}, err
	}

	if len(res) == 0 {
		return product.Product{}, &product.ErrNotFound{ID: id}
	}

	var p product.Product
	if err := json.Unmarshal(res, &p); err != nil {
		return product.Product{}, err
	}

	return p, nil
}

func (r *repository) FindMultiple(ctx context.Context, ids []int) ([]product.Product, error) {
	endpoint, err := url.JoinPath(r.baseUrl, r.getAllEndpoint)
	if err != nil {
		return []product.Product{}, err
	}

	res, _, err := r.requester.Get(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}

	found := make([]product.Product, 0)
	if err := json.Unmarshal(res, &found); err != nil {
		return nil, err
	}

	notFound := []int{}

	pp := make([]product.Product, 0, len(ids))
	for _, id := range ids {
		idx := slices.IndexFunc(found, func(p product.Product) bool {
			return p.ID == id
		})
		if idx < 0 {
			notFound = append(notFound, id)
			continue
		}

		pp = append(pp, found[idx])
	}

	if len(notFound) > 0 {
		return []product.Product{}, &product.ErrProductsNotFound{
			IDs: notFound,
		}
	}

	return pp, nil
}
