package mcardhttp

import (
	"net/url"

	"github.com/varavell/mcard/pkg/httputil"
)

// Client is an interface for interacting with the github API.
type Client interface {
	ReposClient
}

// Config is the information specific to communicating with the microservice for this agent.
type Config struct {
	BaseURL url.URL
}

// V1Client is a concrete Client that calls out to the github api.
type V1Client struct {
	httpu   httputil.HTTPUtilities
	baseURL url.URL
}

// NewV1Client creates a new client for interacting with github api.
func NewV1Client(httpu httputil.HTTPUtilities, config Config) *V1Client {
	return &V1Client{
		httpu:   httpu,
		baseURL: config.BaseURL,
	}
}
