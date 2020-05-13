package desec

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountClient_ObtainCaptcha(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/captcha/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/accounts_captcha.json")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = file.Close() }()

		_, err = io.Copy(rw, file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	captcha, err := client.ObtainCaptcha()
	require.NoError(t, err)

	expected := &Captcha{
		ID:        "aaa",
		Challenge: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8zxD3EwAFuwJYQcBz5wAAAABJRU5ErkJggg==",
	}
	assert.Equal(t, expected, captcha)
}

func TestAccountClient_Register(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	})

	registration := Registration{
		Email:    "email@example.com",
		Password: "secret",
		Captcha: &Captcha{
			ID:       "00010203-0405-0607-0809-0a0b0c0d0e0f",
			Solution: "12H45",
		},
	}

	err := client.Register(registration)
	require.NoError(t, err)
}

func TestAccountClient_Login(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/login/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/accounts_login.json")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = file.Close() }()

		_, err = io.Copy(rw, file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	token, err := client.Login("email@example.com", "secret")
	require.NoError(t, err)

	expected := &Token{
		ID:      "8f9cbae2-c862-48a4-b3f0-2cb1a80df168",
		Name:    "login",
		Value:   "f07Q0TRmEb-CRWPe4h64_iV2jbet",
		Created: mustParseTime("2018-09-06T09:07:43.762697Z"),
	}
	assert.Equal(t, expected, token)
}

func TestAccountClient_Logout(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/logout/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Logout("f07Q0TRmEb-CRWPe4h64_iV2jbet")
	require.NoError(t, err)
}

func TestAccountClient_RetrieveInformation(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/account/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/accounts_retrieve.json")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = file.Close() }()

		_, err = io.Copy(rw, file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	account, err := client.RetrieveInformation("f07Q0TRmEb-CRWPe4h64_iV2jbet")
	require.NoError(t, err)

	expected := &Account{
		Email:        "youremailaddress@example.com",
		LimitDomains: 5,
		Created:      mustParseTime("2019-10-16T18:09:17.715702Z"),
	}
	assert.Equal(t, expected, account)
}

func TestAccountClient_PasswordReset(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	})

	captcha := Captcha{
		ID:       "00010203-0405-0607-0809-0a0b0c0d0e0f",
		Solution: "12H45",
	}

	err := client.PasswordReset("email@example.com", captcha)
	require.NoError(t, err)
}

func TestAccountClient_ChangeEmail(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	})

	err := client.ChangeEmail("email@example.com", "secret", "newemail@example.com")
	require.NoError(t, err)
}

func TestAccountClient_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewAccountClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/account/delete/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
	})

	err := client.Delete("email@example.com", "secret")
	require.NoError(t, err)
}
