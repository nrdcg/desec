package desec

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokensService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := New("token", NewDefaultClientOptions())
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/tokens/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusCreated)

		file, err := os.Open("./fixtures/tokens_create.json")
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
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := New("token", NewDefaultClientOptions())
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/tokens/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/tokens_getall.json")
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
	}
	assert.Equal(t, expected, tokens)
}

func TestTokensService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := New("token", NewDefaultClientOptions())
	client.BaseURL = server.URL

	mux.HandleFunc("/auth/tokens/aaa/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Tokens.Delete(context.Background(), "aaa")
	require.NoError(t, err)
}
