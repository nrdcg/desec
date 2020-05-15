package desec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// TokensService handles communication with the tokens related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/auth/tokens.html
type TokensService struct {
	client *Client

	token string
}

// GetAll retrieving all current tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#retrieving-all-current-tokens
func (s *TokensService) GetAll() ([]Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.token))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var tokens []Token
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return tokens, nil
}

// Create creates additional tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#create-additional-tokens
func (s *TokensService) Create(name string) (*Token, error) {
	endpoint, err := s.client.createEndpoint("auth", "tokens")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Token{Name: name})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.token))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &token, nil
}

// Delete deletes tokens.
// https://desec.readthedocs.io/en/latest/auth/tokens.html#delete-tokens
func (s *TokensService) Delete(tokenID string) error {
	endpoint, err := s.client.createEndpoint("auth", "tokens", tokenID)
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.token))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
