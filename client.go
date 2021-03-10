package x_http_client

import (
	"fmt"
	"net/http"
)

type (
	// Client OSS client
	Client struct {
		Config     *Config      // OSS client configuration
		Conn       *Conn        // Send HTTP request
		HTTPClient *http.Client // http.Client to use - if nil will make its own
	}

	// ClientOption client option such as UseCname, Timeout, SecurityToken.
	ClientOption func(*Client)
)

func New(endpoint, accessAppID, accessAppSecret string, options ...ClientOption) (*Client, error) {
	// Configuration
	config := getDefaultConfig()
	config.Endpoint = endpoint
	config.AccessAppID = accessAppID
	config.AccessAppSecret = accessAppSecret

	// URL parse
	url := &urlMaker{}
	err := url.Init(config.Endpoint)
	if err != nil {
		return nil, err
	}

	// HTTP connect
	conn := &Conn{config: config, url: url}

	// OSS client
	client := &Client{
		Config: config,
		Conn:   conn,
	}

	// Client options parse
	for _, option := range options {
		option(client)
	}

	if config.AuthVersion != AuthV1 && config.AuthVersion != AuthV2 {
		return nil, fmt.Errorf("Init client Error, invalid Auth version: %v", config.AuthVersion)
	}

	// Create HTTP connection
	err = conn.init(config, url, client.HTTPClient)

	return client, err
}
