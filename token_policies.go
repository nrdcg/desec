package desec

import (
	"context"
	"fmt"
	"net/http"
)

// TokenPolicy represents a policy applied to a token.
type TokenPolicy struct {
	ID              string  `json:"id,omitempty"`
	Domain          *string `json:"domain"`
	SubName         *string `json:"subname"`
	Type            *string `json:"type"`
	WritePermission bool    `json:"perm_write,omitempty"`
}

// TokenPoliciesService handles communication with the token policy related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/auth/tokens.html
type TokenPoliciesService struct {
	client *Client
}

// Deprecated: use [TokenPoliciesService.GetAll] instead.
func (s *TokenPoliciesService) Get(ctx context.Context, tokenID string) ([]TokenPolicy, error) {
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

// GetOne retrieves a specific token rrset policy.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokenPoliciesService) GetOne(ctx context.Context, tokenID, policyID string) (*TokenPolicy, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID, "policies", "rrsets", policyID)
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

	policy := &TokenPolicy{}

	err = handleResponse(resp, policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// GetAll retrieves all rrset policies for a token.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokenPoliciesService) GetAll(ctx context.Context, tokenID string) ([]TokenPolicy, error) {
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

// Create creates token policy.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#create-additional-tokens
func (s *TokenPoliciesService) Create(ctx context.Context, tokenID string, policy TokenPolicy) (*TokenPolicy, error) {
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

// Update a token policy
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokenPoliciesService) Update(ctx context.Context, tokenID, policyID string, policy TokenPolicy) (*TokenPolicy, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID, "policies", "rrsets", policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	// Copy values, including only fields that can be modified
	req, err := s.client.newRequest(ctx, http.MethodPatch, endpoint, TokenPolicy{
		Domain:          policy.Domain,
		SubName:         policy.SubName,
		Type:            policy.Type,
		WritePermission: policy.WritePermission,
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

	result := &TokenPolicy{}

	err = handleResponse(resp, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes a token rrset's policy.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#token-policy-management
func (s *TokenPoliciesService) Delete(ctx context.Context, tokenID, policyID string) error {
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
