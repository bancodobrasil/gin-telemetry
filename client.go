package telemetry

import (
	"crypto/tls"
	"net/http"

	"github.com/spf13/viper"
)

// HTTPClient ...
var HTTPClient *http.Client = newHTTPClient()

func newHTTPClient() *http.Client {
	enableTLS := !viper.GetBool("TELEMETRY_HTTP_CLIENT_TLS")
	var transCfg http.RoundTripper
	if enableTLS {
		transCfg = http.DefaultTransport
	} else {
		transCfg = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &http.Client{
		Transport: &propagatorSpanContext{
			core: transCfg,
		},
	}
}

type propagatorSpanContext struct {
	core http.RoundTripper
}

func (p propagatorSpanContext) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	Inject(ctx, r.Header)
	return p.core.RoundTrip(r)
}
