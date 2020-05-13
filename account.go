package desec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Account an account representation.
type Account struct {
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	LimitDomains int        `json:"limit_domains,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
}

// Captcha a captcha representation.
type Captcha struct {
	ID        string `json:"id,omitempty"`
	Challenge string `json:"challenge,omitempty"`
	Solution  string `json:"solution,omitempty"`
}

// Registration a registration representation.
type Registration struct {
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	NewEmail string   `json:"new_email,omitempty"`
	Captcha  *Captcha `json:"captcha,omitempty"`
}

// AccountClient handles communication with the account related methods of the deSEC API.
//
// https://desec.readthedocs.io/en/latest/auth/account.html
type AccountClient struct {
	// HTTP client used to communicate with the API.
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL string
}

// NewAccountClient creates a new AccountClient.
func NewAccountClient() *AccountClient {
	return &AccountClient{
		HTTPClient: http.DefaultClient,
		BaseURL:    defaultBaseURL,
	}
}

// Login Log in.
// https://desec.readthedocs.io/en/latest/auth/account.html#log-in
func (s *AccountClient) Login(email, password string) (*Token, error) {
	endpoint, err := s.createEndpoint("auth", "login")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Account{Email: email, Password: password})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
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

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &token, nil
}

// Logout log out (= delete current token).
// https://desec.readthedocs.io/en/latest/auth/account.html#log-out
func (s *AccountClient) Logout(token string) error {
	endpoint, err := s.createEndpoint("auth", "logout")
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ObtainCaptcha Obtain a captcha.
// https://desec.readthedocs.io/en/latest/auth/account.html#obtain-a-captcha
func (s *AccountClient) ObtainCaptcha() (*Captcha, error) {
	endpoint, err := s.createEndpoint("captcha")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
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

	var captcha Captcha
	err = json.Unmarshal(body, &captcha)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &captcha, nil
}

// Register register account.
// https://desec.readthedocs.io/en/latest/auth/account.html#register-account
func (s *AccountClient) Register(registration Registration) error {
	endpoint, err := s.createEndpoint("auth")
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// RetrieveInformation retrieve account information.
// https://desec.readthedocs.io/en/latest/auth/account.html#retrieve-account-information
func (s *AccountClient) RetrieveInformation(token string) (*Account, error) {
	endpoint, err := s.createEndpoint("auth", "account")
	if err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := s.HTTPClient.Do(req)
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

	var account Account
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to umarshal response body: %w", err)
	}

	return &account, nil
}

// PasswordReset password reset and password change.
// https://desec.readthedocs.io/en/latest/auth/account.html#password-reset
// https://desec.readthedocs.io/en/latest/auth/account.html#password-change
func (s *AccountClient) PasswordReset(email string, captcha Captcha) error {
	endpoint, err := s.createEndpoint("auth", "account", "reset-password")
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Registration{Email: email, Captcha: &captcha})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ChangeEmail changes email address.
// https://desec.readthedocs.io/en/latest/auth/account.html#change-email-address
func (s *AccountClient) ChangeEmail(email, password string, newEmail string) error {
	endpoint, err := s.createEndpoint("auth", "account", "change-email")
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Registration{Email: email, Password: password, NewEmail: newEmail})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Delete deletes account.
// https://desec.readthedocs.io/en/latest/auth/account.html#delete-account
func (s *AccountClient) Delete(email, password string) error {
	endpoint, err := s.createEndpoint("auth", "account", "delete")
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %w", err)
	}

	raw, err := json.Marshal(Account{Email: email, Password: password})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("error: %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *AccountClient) createEndpoint(parts ...string) (*url.URL, error) {
	return createEndpoint(s.BaseURL, parts)
}
