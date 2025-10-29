package desec

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokensService_Create(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("POST /auth/tokens/", fromFixtures("tokens_create.json", http.StatusCreated))

	newToken, err := client.Tokens.Create(context.Background(), "my new token")
	require.NoError(t, err)

	expected := &Token{
		ID:      "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
		Name:    "my new token",
		Value:   "4pnk7u-NHvrEkFzrhFDRTjGFyX_S",
		Created: mustParseTime("2018-09-06T09:08:43.762697Z"),
	}
	assert.Equal(t, expected, newToken)
}

func TestTokensService_GetAll(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /auth/tokens/", fromFixtures("tokens_getall.json", http.StatusOK))

	tokens, err := client.Tokens.GetAll(context.Background())
	require.NoError(t, err)

	expected := []Token{
		{
			ID:      "3159e485-5499-46c0-ae2b-aeb84d627a8e",
			Name:    "login",
			Created: mustParseTime("2018-09-06T07:05:54.080564Z"),
		},
		{
			ID:      "76d6e39d-65bc-4ab2-a1b7-6e94eee0a534",
			Name:    "sample",
			Created: mustParseTime("2018-09-06T08:53:26.428396Z"),
		},
		{
			ID:               "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
			Created:          mustParseTime("2018-09-06T09:08:43.762697Z"),
			Owner:            "youremailaddress@example.com",
			Name:             "my new token",
			PermCreateDomain: false,
			PermDeleteDomain: false,
			PermManageTokens: false,
			AllowedSubnets: []string{
				"0.0.0.0/0",
				"::/0",
			},
			AutoPolicy: false,
			Value:      "4pnk7u-NHvrEkFzrhFDRTjGFyX_S",
		},
	}
	assert.Equal(t, expected, tokens)
}

func TestTokensService_Get(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /auth/tokens/3a6b94b5-d20e-40bd-a7cc-521f5c79fab3/",
		fromFixtures("tokens_get.json", http.StatusOK))

	token, err := client.Tokens.Get(context.Background(), "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3")
	require.NoError(t, err)

	expected := &Token{
		ID:               "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
		Created:          mustParseTime("2018-09-06T09:08:43.762697Z"),
		Owner:            "youremailaddress@example.com",
		Name:             "my new token",
		PermCreateDomain: false,
		PermDeleteDomain: false,
		PermManageTokens: false,
		AllowedSubnets: []string{
			"0.0.0.0/0",
			"::/0",
		},
		AutoPolicy: false,
		Value:      "4pnk7u-NHvrEkFzrhFDRTjGFyX_S",
	}
	assert.Equal(t, expected, token)
}

func TestTokensService_Update(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("PATCH /auth/tokens/3a6b94b5-d20e-40bd-a7cc-521f5c79fab3/",
		fromFixtures("tokens_update.json", http.StatusOK))

	token, err := client.Tokens.Update(context.Background(), "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3", &Token{
		PermCreateDomain: true,
		PermDeleteDomain: false,
		PermManageTokens: true,
	})

	require.NoError(t, err)

	expected := &Token{
		ID:               "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
		Created:          mustParseTime("2018-09-06T09:08:43.762697Z"),
		Owner:            "youremailaddress@example.com",
		Name:             "my new token",
		PermCreateDomain: true,
		PermDeleteDomain: false,
		PermManageTokens: true,
		AllowedSubnets: []string{
			"0.0.0.0/0",
			"::/0",
		},
		AutoPolicy: false,
		Value:      "4pnk7u-NHvrEkFzrhFDRTjGFyX_S",
	}
	assert.Equal(t, expected, token)
}

func TestTokensService_Delete(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("DELETE /auth/tokens/aaa/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Tokens.Delete(context.Background(), "aaa")
	require.NoError(t, err)
}
