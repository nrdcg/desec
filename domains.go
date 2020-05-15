package desec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Domain a domain representation.
type Domain struct {
	Name       string      `json:"name,omitempty"`
	MinimumTTL int         `json:"minimum_ttl,omitempty"`
	Keys       []DomainKey `json:"keys,omitempty"`
	Created    *time.Time  `json:"created,omitempty"`
	Published  *time.Time  `json:"published,omitempty"`
}

// DomainKey a domain key representation.
type DomainKey struct {
	DNSKey  string   `json:"dnskey,omitempty"`
	DS      []string `json:"ds,omitempty"`
	Flags   int      `json:"flags,omitempty"`
	KeyType string   `json:"keytype,omitempty"`
}

// DomainsService handles communication with the domain related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/dns/domains.html
type DomainsService struct {
	client *Client

	token string
}

// Create creating a domain.
// https://desec.readthedocs.io/en/latest/dns/domains.html#creating-a-domain
func (s *DomainsService) Create(domainName string) (*Domain, error) {
	endpoint, err := s.client.createEndpoint("domains")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Domain{Name: domainName})
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

	var domain Domain
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &domain, nil
}

// GetAll listing domains.
// https://desec.readthedocs.io/en/latest/dns/domains.html#listing-domains
func (s *DomainsService) GetAll() ([]Domain, error) {
	endpoint, err := s.client.createEndpoint("domains")
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

	var domains []Domain
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return domains, nil
}

// Get retrieving a specific domain.
// https://desec.readthedocs.io/en/latest/dns/domains.html#retrieving-a-specific-domain
func (s *DomainsService) Get(domainName string) (*Domain, error) {
	endpoint, err := s.client.createEndpoint("domains", domainName)
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

	var domains Domain
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &domains, nil
}

// Delete deleting a domain.
// https://desec.readthedocs.io/en/latest/dns/domains.html#deleting-a-domain
func (s *DomainsService) Delete(domainName string) error {
	endpoint, err := s.client.createEndpoint("domains", domainName)
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
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
