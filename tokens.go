package desec

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Token a token representation.
//
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-field-reference
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

// Get retrieves a specific token.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#retrieving-a-specific-token
// NOTE: This method used to retrieve all policies for a token, that is now done by GetAll.
func (s *TokensService) Get(ctx context.Context, id string) (*Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", id)
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	token := &Token{}

	err = handleResponse(resp, token)
	if err != nil {
		return nil, err
	}

	return token, nil
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

// Update a token.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#modifying-a-token
func (s *TokensService) Update(ctx context.Context, id string, token *Token) (*Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", id)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	// Copy values, including only fields that can be modified
	req, err := s.client.newRequest(ctx, http.MethodPatch, endpoint, Token{
		Owner:            token.Owner,
		UserOverride:     token.UserOverride,
		Name:             token.Name,
		PermCreateDomain: token.PermCreateDomain,
		PermDeleteDomain: token.PermDeleteDomain,
		PermManageTokens: token.PermManageTokens,
		AllowedSubnets:   token.AllowedSubnets,
		AutoPolicy:       token.AutoPolicy,
	})
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

	result := &Token{}

	err = handleResponse(resp, result)
	if err != nil {
		return nil, err
	}

	return result, nil
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
