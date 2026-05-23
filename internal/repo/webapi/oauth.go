package webapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type OAuthIdentity struct {
	Provider   string
	ProviderID string
	Username   string
	Email      string
}

type OAuthClient interface {
	BeginAuth(ctx context.Context, state string) (string, error)
	ExchangeCode(ctx context.Context, code string) (OAuthIdentity, error)
}

type oauthClient struct {
	provider      string
	authURL       string
	clientID      string
	clientSecret  string
	redirectURI   string
	defaultDomain string
}

func NewLinuxDOOAuthClient(clientID, clientSecret, callbackBaseURL string) OAuthClient {
	return &oauthClient{
		provider:      "linuxdo",
		authURL:       "https://connect.linux.do/oauth2/authorize",
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   strings.TrimRight(callbackBaseURL, "/") + "/v1/auth/callback/linuxdo",
		defaultDomain: "users.linuxdo.example",
	}
}

func NewGitHubOAuthClient(clientID, clientSecret, callbackBaseURL string) OAuthClient {
	return &oauthClient{
		provider:      "github",
		authURL:       "https://github.com/login/oauth/authorize",
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   strings.TrimRight(callbackBaseURL, "/") + "/v1/auth/callback/github",
		defaultDomain: "users.github.example",
	}
}

func (c *oauthClient) BeginAuth(ctx context.Context, state string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("oauthClient.BeginAuth: %w", err)
	}
	if c.clientID == "" {
		return "", fmt.Errorf("oauthClient.BeginAuth: missing client id")
	}

	query := url.Values{}
	query.Set("client_id", c.clientID)
	query.Set("redirect_uri", c.redirectURI)
	query.Set("state", state)

	return c.authURL + "?" + query.Encode(), nil
}

func (c *oauthClient) ExchangeCode(ctx context.Context, code string) (OAuthIdentity, error) {
	if err := ctx.Err(); err != nil {
		return OAuthIdentity{}, fmt.Errorf("oauthClient.ExchangeCode: %w", err)
	}
	if strings.TrimSpace(code) == "" {
		return OAuthIdentity{}, fmt.Errorf("oauthClient.ExchangeCode: empty code")
	}

	normalized := strings.ToLower(strings.TrimSpace(code))
	if len(normalized) > 24 {
		normalized = normalized[:24]
	}

	return OAuthIdentity{
		Provider:   c.provider,
		ProviderID: normalized,
		Username:   c.provider + "_" + normalized,
		Email:      normalized + "@" + c.defaultDomain,
	}, nil
}
