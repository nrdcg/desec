package desec

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountClient_ObtainCaptcha(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /captcha/", fromFixtures("accounts_captcha.json", http.StatusOK))

	captcha, err := client.Account.ObtainCaptcha(context.Background())
	require.NoError(t, err)

	expected := &Captcha{
		ID:        "aaa",
		Challenge: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8zxD3EwAFuwJYQcBz5wAAAABJRU5ErkJggg==",
	}
	assert.Equal(t, expected, captcha)
}

func TestAccountClient_Register(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /auth/", func(rw http.ResponseWriter, _ *http.Request) {
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

	err := client.Account.Register(context.Background(), registration)
	require.NoError(t, err)
}

func TestAccountClient_Login(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /auth/login/", fromFixtures("accounts_login.json", http.StatusOK))

	token, err := client.Account.Login(context.Background(), "email@example.com", "secret")
	require.NoError(t, err)

	expected := &Token{
		ID:      "8f9cbae2-c862-48a4-b3f0-2cb1a80df168",
		Name:    "login",
		Value:   "f07Q0TRmEb-CRWPe4h64_iV2jbet",
		Created: mustParseTime("2018-09-06T09:07:43.762697Z"),
	}
	assert.Equal(t, expected, token)

	assert.Equal(t, expected.Value, client.token)
}

func TestAccountClient_Logout(t *testing.T) {
	client, mux := setupTest(t, "f07Q0TRmEb-CRWPe4h64_iV2jbet")

	mux.HandleFunc("POST /auth/logout/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Account.Logout(context.Background())
	require.NoError(t, err)

	assert.Empty(t, client.token)
}

func TestAccountClient_RetrieveInformation(t *testing.T) {
	client, mux := setupTest(t, "f07Q0TRmEb-CRWPe4h64_iV2jbet")

	mux.HandleFunc("POST /auth/account/", fromFixtures("accounts_retrieve.json", http.StatusOK))

	account, err := client.Account.RetrieveInformation(context.Background())
	require.NoError(t, err)

	expected := &Account{
		Email:        "youremailaddress@example.com",
		LimitDomains: 5,
		Created:      mustParseTime("2019-10-16T18:09:17.715702Z"),
	}
	assert.Equal(t, expected, account)
}

func TestAccountClient_PasswordReset(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	})

	captcha := Captcha{
		ID:       "00010203-0405-0607-0809-0a0b0c0d0e0f",
		Solution: "12H45",
	}

	err := client.Account.PasswordReset(context.Background(), "email@example.com", captcha)
	require.NoError(t, err)
}

func TestAccountClient_ChangeEmail(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	})

	err := client.Account.ChangeEmail(context.Background(), "email@example.com", "secret", "newemail@example.com")
	require.NoError(t, err)
}

func TestAccountClient_Delete(t *testing.T) {
	client, mux := setupTest(t, "")

	mux.HandleFunc("POST /auth/account/delete/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusAccepted)
	})

	err := client.Account.Delete(context.Background(), "email@example.com", "secret")
	require.NoError(t, err)
}
