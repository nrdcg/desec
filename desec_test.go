package desec

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTest(t *testing.T, token string) (*Client, *http.ServeMux) {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	options := NewDefaultClientOptions()
	options.HTTPClient = server.Client()

	client := New(token, options)
	client.BaseURL = server.URL

	return client, mux
}

func fromFixtures(filename string, statusCode int) http.HandlerFunc {
	return func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		file, err := os.Open(filepath.Clean(filepath.Join("fixtures", filename)))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}

		defer func() { _ = file.Close() }()

		rw.WriteHeader(statusCode)

		_, err = io.Copy(rw, file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}

func mustParseTime(value string) *time.Time {
	date, _ := time.Parse(time.RFC3339, value)
	return &date
}
