package util

import (
	"crypto/tls"
	"net/http"
)

func NewHTTPClient() http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return http.Client{Transport: transport}
}
