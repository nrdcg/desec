package desec

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ApexZone apex zone name.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#accessing-the-zone-apex
const ApexZone = "@"

// IgnoreFilter is a specific value used to ignore a filter field.
const IgnoreFilter = "#IGNORE#"

// RRSet DNS Record Set.
type RRSet struct {
	Name    string     `json:"name,omitempty"`
	Domain  string     `json:"domain,omitempty"`
	SubName string     `json:"subname,omitempty"`
	Type    string     `json:"type,omitempty"`
	Records []string   `json:"records"`
	TTL     int        `json:"ttl,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	Touched *time.Time `json:"touched,omitempty"`
}

// RRSetFilter a RRSets filter.
type RRSetFilter struct {
	Type    string
	SubName string
}

// FilterRRSetOnlyOnType creates an RRSetFilter that ignore SubName.
func FilterRRSetOnlyOnType(t string) RRSetFilter {
	return RRSetFilter{
		Type:    t,
		SubName: IgnoreFilter,
	}
}

// FilterRRSetOnlyOnSubName creates an RRSetFilter that ignore Type.
func FilterRRSetOnlyOnSubName(n string) RRSetFilter {
	return RRSetFilter{
		Type:    IgnoreFilter,
		SubName: n,
	}
}

// RecordsService handles communication with the records related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/dns/rrsets.html
type RecordsService struct {
	client *Client
}

/*
	Domains
*/

// GetAll retrieving all RRSets in a zone.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#retrieving-all-rrsets-in-a-zone
func (s *RecordsService) GetAll(ctx context.Context, domainName string, filter *RRSetFilter) ([]RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	if filter != nil {
		query := endpoint.Query()

		if filter.Type != IgnoreFilter {
			query.Set("type", filter.Type)
		}

		if filter.SubName != IgnoreFilter {
			query.Set("subname", filter.SubName)
		}

		endpoint.RawQuery = query.Encode()
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

	var rrSets []RRSet
	err = handleResponse(resp, &rrSets)
	if err != nil {
		return nil, err
	}

	return rrSets, nil
}

// Create creates a new RRSet.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#creating-a-tlsa-rrset
func (s *RecordsService) Create(ctx context.Context, rrSet RRSet) (*RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", rrSet.Domain, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, endpoint, rrSet)
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

	var newRRSet RRSet
	err = handleResponse(resp, &newRRSet)
	if err != nil {
		return nil, err
	}

	return &newRRSet, nil
}

/*
	Domains + subname + type
*/

// Get gets a RRSet.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#retrieving-a-specific-rrset
func (s *RecordsService) Get(ctx context.Context, domainName, subName, recordType string) (*RRSet, error) {
	if subName == "" {
		subName = ApexZone
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
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

	var rrSet RRSet
	err = handleResponse(resp, &rrSet)
	if err != nil {
		return nil, err
	}

	return &rrSet, nil
}

// Update updates RRSet (PATCH).
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#modifying-an-rrset
func (s *RecordsService) Update(ctx context.Context, domainName, subName, recordType string, rrSet RRSet) (*RRSet, error) {
	if subName == "" {
		subName = ApexZone
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPatch, endpoint, rrSet)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// when a RRSet is deleted (empty records)
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	var updatedRRSet RRSet
	err = handleResponse(resp, &updatedRRSet)
	if err != nil {
		return nil, err
	}

	return &updatedRRSet, nil
}

// Replace replaces a RRSet (PUT).
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#modifying-an-rrset
func (s *RecordsService) Replace(ctx context.Context, domainName, subName, recordType string, rrSet RRSet) (*RRSet, error) {
	if subName == "" {
		subName = ApexZone
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPut, endpoint, rrSet)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	// when a RRSet is deleted (empty records)
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	var updatedRRSet RRSet
	err = handleResponse(resp, &updatedRRSet)
	if err != nil {
		return nil, err
	}

	return &updatedRRSet, nil
}

// Delete deletes a RRSet.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#deleting-an-rrset
func (s *RecordsService) Delete(ctx context.Context, domainName, subName, recordType string) error {
	if subName == "" {
		subName = ApexZone
	}

	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets", subName, recordType)
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

/*
	Bulk operations
*/

// UpdateMode the mode used to bulk update operations.
type UpdateMode string

const (
	// FullResource the full resource must be specified.
	FullResource UpdateMode = http.MethodPut
	// OnlyFields only fields you would like to modify need to be provided.
	OnlyFields UpdateMode = http.MethodPatch
)

// BulkCreate creates new RRSets in bulk.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#bulk-creation-of-rrsets
func (s *RecordsService) BulkCreate(ctx context.Context, domainName string, rrSets []RRSet) ([]RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, http.MethodPost, endpoint, rrSets)
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

	var newRRSets []RRSet
	err = handleResponse(resp, &newRRSets)
	if err != nil {
		return nil, err
	}

	return newRRSets, nil
}

// BulkUpdate updates RRSets in bulk.
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#bulk-modification-of-rrsets
func (s *RecordsService) BulkUpdate(ctx context.Context, mode UpdateMode, domainName string, rrSets []RRSet) ([]RRSet, error) {
	endpoint, err := s.client.createEndpoint("domains", domainName, "rrsets")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := s.client.newRequest(ctx, string(mode), endpoint, rrSets)
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

	var results []RRSet
	err = handleResponse(resp, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// BulkDelete deletes RRSets in bulk (uses FullResourceUpdateMode).
// https://desec.readthedocs.io/en/latest/dns/rrsets.html#bulk-deletion-of-rrsets
func (s *RecordsService) BulkDelete(ctx context.Context, domainName string, rrSets []RRSet) error {
	deleteRRSets := make([]RRSet, len(rrSets))
	for i, rrSet := range rrSets {
		rrSet.Records = []string{}
		deleteRRSets[i] = rrSet
	}

	_, err := s.BulkUpdate(ctx, FullResource, domainName, deleteRRSets)
	if err != nil {
		return err
	}

	return nil
}
