// Package traefik_nova_plugin a plugin to use the Nova WAF in traefik
package traefik_nova_plugin

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

type Config struct {
	NovaContainerUrl string `json:"novaContainerUrl,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type Nova struct {
	next             http.Handler
	novaContainerUrl string
	name             string
	logger           *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.NovaContainerUrl) == 0 {
		return nil, fmt.Errorf("novaContainerUrl cannot be empty")
	}

	return &Nova{
		novaContainerUrl: config.NovaContainerUrl,
		next:             next,
		name:             name,
		logger:           log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

func (a *Nova) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Websocket scanning is not supported
	if isWebsocket(req) {
		a.next.ServeHTTP(rw, req)
		return
	}

	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		a.logger.Printf("Failed to read incoming request: %s", err.Error())
		http.Error(rw, "", http.StatusBadGateway)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s%s", a.novaContainerUrl, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

	// Overwrite Host header with original header for Nova logging
	// which uses the Host header
	proxyReq.Host = req.Host

	if err != nil {
		a.logger.Printf("Failed to prepare forwarded request: %s", err.Error())
		http.Error(rw, "", http.StatusBadGateway)
		return
	}

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		a.logger.Printf("Failed to send HTTP request to Nova (is it running?): %s", err.Error())
		http.Error(rw, "", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		forwardResponse(resp, rw)
		return
	}

	a.next.ServeHTTP(rw, req)
}

func isWebsocket(req *http.Request) bool {
	for _, header := range req.Header["Upgrade"] {
		if header == "websocket" {
			return true
		}
	}
	return false
}

func forwardResponse(resp *http.Response, rw http.ResponseWriter) {
	for k, vv := range resp.Header {
		for _, v := range vv {
			rw.Header().Set(k, v)
		}
	}

	rw.WriteHeader(resp.StatusCode)

	io.Copy(rw, resp.Body)
}
