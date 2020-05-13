package desec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// RRSet DNS Record Set.
type RRSet struct {
	Name    string     `json:"name,omitempty"`
	Domain  string     `json:"domain,omitempty"`
	SubName string     `json:"subname,omitempty"`
	Type    string     `json:"type,omitempty"`
	Records []string   `json:"records"`
	TTL     int        `json:"ttl,omitempty"`
	Created *time.Time `json:"created,omitempty"`
}

// RRSetFilter a RRsets filter.
type RRSetFilter struct {
	Type    string
	SubName string
}

// RecordsService handles communication with the records related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/dns/rrsets.html
type RecordsService struct {
	client *Client

	token string
}

// Get gets a RRSet.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#retrieving-a-specific-rrset
func (s *RecordsService) Get(domainName, subName string, recordType string) (*RRSet, error) {
	if subName == "" {
		subName = "@"
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		var notFound NotFound
		err = json.Unmarshal(body, &notFound)
		if err != nil {
			return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
		}

		return nil, &notFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var rrSet RRSet
	err = json.Unmarshal(body, &rrSet)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &rrSet, nil
}

// GetAll retrieving all RRsets in a zone.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#retrieving-all-rrsets-in-a-zone
func (s *RecordsService) GetAll(domainName string, filter *RRSetFilter) ([]RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	if filter != nil {
		query := endpoint.Query()
		query.Set("type", filter.Type)
		query.Set("subname", filter.SubName)
		endpoint.RawQuery = query.Encode()
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var rrSets []RRSet
	err = json.Unmarshal(body, &rrSets)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return rrSets, nil
}

// Create creates a new RRSet.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#creating-a-tlsa-rrset
func (s *RecordsService) Create(rrSet RRSet) (*RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", rrSet.Domain, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(rrSet)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var newRRSet RRSet
	err = json.Unmarshal(body, &newRRSet)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &newRRSet, nil
}

// Update updates RRSet records.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#modifying-an-rrset
func (s *RecordsService) Update(domainName string, subName string, recordType string, records []string) (*RRSet, error) {
	if subName == "" {
		subName = "@"
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(RRSet{Records: records})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.token))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// when a RRSet is deleted (empty records)
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	var updatedRRSet RRSet
	err = json.Unmarshal(body, &updatedRRSet)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &updatedRRSet, nil
}

// Delete deletes a RRset.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#deleting-an-rrset
func (s *RecordsService) Delete(domainName string, subName string, recordType string) error {
	if subName == "" {
		subName = "@"
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
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

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
