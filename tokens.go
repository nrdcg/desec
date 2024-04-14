package desec

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Token a token representation.
type Token struct {
	ID      string     `json:"id,omitempty"`
	Name    string     `json:"name,omitempty"`
	Value   string     `json:"token,omitempty"`
	Created *time.Time `json:"created,omitempty"`
}

// TokenPolicy represents a policy applied to a token.
type TokenPolicy struct {
	ID              string  `json:"id,omitempty"`
	Domain          *string `json:"domain"`
	SubName         *string `json:"subname"`
	Type            *string `json:"type"`
	WritePermission bool    `json:"perm_write,omitempty"` // Go `encoding/json` default boolean value is false
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

// GetPolicies retrieves token rrset's policies.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokensService) GetPolicies(ctx context.Context, tokenID string) ([]TokenPolicy, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID, "policies", "rrsets")
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

	var policies []TokenPolicy
	err = handleResponse(resp, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

// CreatePolicy creates token policy.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#create-additional-tokens
func (s *TokensService) CreatePolicy(ctx context.Context, tokenID string, policy TokenPolicy) (*TokenPolicy, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID, "policies", "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, endpoint, policy)
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

	var tokenPolicy TokenPolicy
	err = handleResponse(resp, &tokenPolicy)
	if err != nil {
		return nil, err
	}

	return &tokenPolicy, nil
}

// DeletePolicy deletes a token rrset's policy.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokensService) DeletePolicy(ctx context.Context, tokenID, policyID string) error {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID, "policies", "rrsets", policyID)
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
