package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// OAuthClient represents a Zendesk OAuth client
type OAuthClient struct {
	ID          int64    `json:"id,omitempty"`
	URL         string   `json:"url,omitempty"`
	UserID      int64    `json:"user_id,omitempty"`
	Name        string   `json:"name"`
	Identifier  string   `json:"identifier,omitempty"`
	Kind        string   `json:"kind,omitempty"`
	Company     string   `json:"company,omitempty"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	RedirectURI []string `json:"redirect_uri,omitempty"`
	LogoURL     string   `json:"logo_url,omitempty"`
	Global      bool     `json:"global,omitempty"`
	// Pointer so that an explicit false is sent rather than dropped; Zendesk defaults this to true.
	EnforceTokenExpiration *bool `json:"enforce_token_expiration,omitempty"`
	// Returned in full only by Create. Every later read answers with a truncated stub.
	Secret    string `json:"secret,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// OAuthClientAPI interface for OAuth client operations
type OAuthClientAPI interface {
	GetOAuthClient(ctx context.Context, id int64) (OAuthClient, error)
	CreateOAuthClient(ctx context.Context, oauthClient OAuthClient) (OAuthClient, error)
	UpdateOAuthClient(ctx context.Context, id int64, oauthClient OAuthClient) (OAuthClient, error)
	DeleteOAuthClient(ctx context.Context, id int64) error
}

// GetOAuthClient returns a specific OAuth client
// ref: https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/#show-client
func (z *Client) GetOAuthClient(ctx context.Context, id int64) (OAuthClient, error) {
	var result struct {
		Client OAuthClient `json:"client"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/oauth/clients/%d.json", id))
	if err != nil {
		return OAuthClient{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OAuthClient{}, err
	}

	return result.Client, nil
}

// CreateOAuthClient creates a new OAuth client. This is the only response carrying the full secret.
// ref: https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/#create-client
func (z *Client) CreateOAuthClient(ctx context.Context, oauthClient OAuthClient) (OAuthClient, error) {
	var data, result struct {
		Client OAuthClient `json:"client"`
	}
	data.Client = oauthClient

	body, err := z.Post(ctx, "/oauth/clients.json", data)
	if err != nil {
		return OAuthClient{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OAuthClient{}, err
	}

	return result.Client, nil
}

// UpdateOAuthClient updates an OAuth client. Zendesk rejects a change to name, identifier or kind.
// ref: https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/#update-client
func (z *Client) UpdateOAuthClient(ctx context.Context, id int64, oauthClient OAuthClient) (OAuthClient, error) {
	var data, result struct {
		Client OAuthClient `json:"client"`
	}
	data.Client = oauthClient

	body, err := z.Put(ctx, fmt.Sprintf("/oauth/clients/%d.json", id), data)
	if err != nil {
		return OAuthClient{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OAuthClient{}, err
	}

	return result.Client, nil
}

// DeleteOAuthClient deletes an OAuth client
// ref: https://developer.zendesk.com/api-reference/ticketing/oauth/oauth_clients/#delete-client
func (z *Client) DeleteOAuthClient(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/oauth/clients/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
