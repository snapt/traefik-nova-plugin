// Package traefik_nova_plugin a plugin to use the Nova WAF in traefik
package traefik_nova_plugin

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNova_ServeHTTP(t *testing.T) {

	req, err := http.NewRequest(http.MethodGet, "http://proxy.com/test", bytes.NewBuffer([]byte("Request")))

	if err != nil {
		log.Fatal(err)
	}

	type response struct {
		Body       string
		StatusCode int
	}

	serviceResponse := response{
		StatusCode: 200,
		Body:       "Response from service",
	}

	tests := []struct {
		name            string
		request         http.Request
		wafResponse     response
		serviceResponse response
		expectBody      string
		expectStatus    int
	}{
		{
			name:    "Forward request when WAF found no threats",
			request: *req,
			wafResponse: response{
				StatusCode: 200,
				Body:       "Response from waf",
			},
			serviceResponse: serviceResponse,
			expectBody:      "Response from service",
			expectStatus:    200,
		},
		{
			name:    "Intercepts request when WAF found threats",
			request: *req,
			wafResponse: response{
				StatusCode: 403,
				Body:       "Response from waf",
			},
			serviceResponse: serviceResponse,
			expectBody:      "Response from waf",
			expectStatus:    403,
		},
		{
			name: "Does not forward Websockets",
			request: http.Request{
				Body: http.NoBody,
				Header: http.Header{
					"Upgrade": []string{"Websocket"},
				},
			},
			wafResponse: response{
				StatusCode: 200,
				Body:       "Response from waf",
			},
			serviceResponse: serviceResponse,
			expectBody:      "Response from service",
			expectStatus:    200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			novaMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := http.Response{
					Body:       io.NopCloser(bytes.NewReader([]byte(tt.wafResponse.Body))),
					StatusCode: tt.wafResponse.StatusCode,
					Header:     http.Header{},
				}
				forwardResponse(&resp, w)
			}))

			httpServiceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := http.Response{
					Body:       io.NopCloser(bytes.NewReader([]byte(tt.serviceResponse.Body))),
					StatusCode: tt.serviceResponse.StatusCode,
					Header:     http.Header{},
				}
				forwardResponse(&resp, w)
			})

			middleware := &Nova{
				next:             httpServiceHandler,
				novaContainerUrl: novaMockServer.URL,
				name:             "nova-middleware",
			}

			rw := httptest.NewRecorder()

			middleware.ServeHTTP(rw, &tt.request)

			resp := rw.Result()
			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectBody, string(body))
			assert.Equal(t, tt.expectStatus, resp.StatusCode)
		})
	}
}
