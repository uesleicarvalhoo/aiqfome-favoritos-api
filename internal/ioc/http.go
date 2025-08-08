package ioc

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var (
	httpClient     *http.Client
	httpClientOnce sync.Once
)

func HttpClient() *http.Client {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Timeout: config.GetDuration("HTTP_CLIENT_TIMEOUT"),
			Transport: otelhttp.NewTransport(http.DefaultTransport,
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					return fmt.Sprintf("HTTP %s - %s", r.Method, r.Host)
				}),
				otelhttp.WithTracerProvider(otel.GetTracerProvider()),
			),
		}
	})

	return httpClient
}
