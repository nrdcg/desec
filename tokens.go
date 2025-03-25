package desec

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/* from https://desec.readthedocs.io/en/latest/auth/tokens.html#token-field-reference
{
		"id": "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
		"created": "2018-09-06T09:08:43.762697Z",
		"last_used": null,
		"owner": "youremailaddress@example.com"",
		"user_override": null,
		"max_age": "365 00:00:00",
		"max_unused_period": null,
		"name": "my new token",
		"perm_create_domain": false,
		"perm_delete_domain": false,
		"perm_manage_tokens": false,
		"allowed_subnets": [
				"0.0.0.0/0",
				"::/0"
		],
		"auto_policy": false,
		"token": "4pnk7u-NHvrEkFzrhFDRTjGFyX_S"
}
*/

// Token a token representation.
type Token struct {
	ID               string     `json:"id,omitempty"`
	Created          *time.Time `json:"created,omitempty"`
	LastUsed         *time.Time `json:"last_used,omitempty"`
	Owner            string     `json:"owner,omitempty"`
	UserOverride     string     `json:"user_override,omitempty"`
	Name             string     `json:"name,omitempty"`
	PermCreateDomain bool       `json:"perm_create_domain"`
	PermDeleteDomain bool       `json:"perm_delete_domain"`
	PermManageTokens bool       `json:"perm_manage_tokens"`
	IsValid          bool       `json:"is_valid,omitempty"`
	AllowedSubnets   []string   `json:"allowed_subnets,omitempty"`
	AutoPolicy       bool       `json:"auto_policy"`
	Value            string     `json:"token,omitempty"`
	// Not currently implemented
	// MaxAge           *time.Duration `json:"name,omitempty"`
	// MaxUnusedPeriod  *time.Duration `json:"name,omitempty"`
}

// TokensService handles communication with the tokens related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/auth/tokens.html
type TokensService struct {
	client *Client
}

// GetAll retrieving all current tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#retrieving-all-current-tokens
func (s *TokensService) GetAll(ctx context.Context) ([]Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	var tokens []Token

	err = handleResponse(resp, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Create creates additional tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#create-additional-tokens
func (s *TokensService) Create(ctx context.Context, name string) (*Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, endpoint, Token{Name: name})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, handleError(resp)
	}

	var token Token

	err = handleResponse(resp, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Delete deletes tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#delete-tokens
func (s *TokensService) Delete(ctx context.Context, tokenID string) error {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID)
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		return handleError(resp)
	}

	return nil
}
