package requester

import (
	"context"
	"io"
	"net/http"
)

type Requester interface {
	Get(ctx context.Context, url string, params map[string]string) ([]byte, int, error)
}

type requester struct {
	client *http.Client
}

func New(client *http.Client) Requester {
	return &requester{
		client: client,
	}
}

func (r *requester) Get(ctx context.Context, url string, params map[string]string) ([]byte, int, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	if params != nil {
		q := rq.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}

		rq.URL.RawQuery = q.Encode()
	}

	res, err := r.client.Do(rq)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return b, res.StatusCode, nil
}
