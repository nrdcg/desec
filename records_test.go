package desec

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordsService_Create(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient("token")
	client.BaseURL = server.URL

	mux.HandleFunc("/domains/example.dedyn.io/rrsets/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		file, err := os.Open("./fixtures/records_create.json")
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

	record := RRSet{
		Name:    "",
		Domain:  "example.dedyn.io",
		SubName: "_acme-challenge",
		Type:    "TXT",
		Records: []string{`"txt"`},
		TTL:     300,
	}

	newRecord, err := client.Records.Create(record)
	require.NoError(t, err)

	expected := &RRSet{
		Name:    "_acme-challenge.example.dedyn.io.",
		Domain:  "example.dedyn.io",
		SubName: "_acme-challenge",
		Type:    "TXT",
		Records: []string{`"txt"`},
		TTL:     300,
		Created: mustParseTime("2020-05-06T11:46:07.641885Z"),
	}
	assert.Equal(t, expected, newRecord)
}

func TestRecordsService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient("token")
	client.BaseURL = server.URL

	mux.HandleFunc("/domains/example.dedyn.io/rrsets/_acme-challenge/TXT/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Records.Delete("example.dedyn.io", "_acme-challenge", "TXT")
	require.NoError(t, err)
}

func TestRecordsService_Get(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient("token")
	client.BaseURL = server.URL

	mux.HandleFunc("/domains/example.dedyn.io/rrsets/_acme-challenge/TXT/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/records_get.json")
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

	record, err := client.Records.Get("example.dedyn.io", "_acme-challenge", "TXT")
	require.NoError(t, err)

	expected := &RRSet{
		Name:    "_acme-challenge.example.dedyn.io.",
		Domain:  "example.dedyn.io",
		SubName: "_acme-challenge",
		Type:    "TXT",
		Records: []string{`"txt"`},
		TTL:     300,
		Created: mustParseTime("2020-05-06T11:46:07.641885Z"),
	}
	assert.Equal(t, expected, record)
}

func TestRecordsService_Update(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient("token")
	client.BaseURL = server.URL

	mux.HandleFunc("/domains/example.dedyn.io/rrsets/_acme-challenge/TXT/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/records_update.json")
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

	updatedRecord, err := client.Records.Update("example.dedyn.io", "_acme-challenge", "TXT", []string{`"updated"`})
	require.NoError(t, err)

	expected := &RRSet{
		Name:    "_acme-challenge.example.dedyn.io.",
		Domain:  "example.dedyn.io",
		SubName: "_acme-challenge",
		Type:    "TXT",
		Records: []string{`"updated"`},
		TTL:     300,
		Created: mustParseTime("2020-05-06T11:46:07.641885Z"),
	}
	assert.Equal(t, expected, updatedRecord)
}

func TestRecordsService_GetAll(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := NewClient("token")
	client.BaseURL = server.URL

	mux.HandleFunc("/domains/example.dedyn.io/rrsets/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(rw, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open("./fixtures/records_getall.json")
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

	records, err := client.Records.GetAll("example.dedyn.io", nil)
	require.NoError(t, err)

	expected := []RRSet{
		{
			Name:    "example.dedyn.io.",
			Domain:  "example.dedyn.io",
			SubName: "",
			Type:    "A",
			Records: []string{"10.10.10.10"},
			TTL:     60,
			Created: mustParseTime("2020-05-06T11:46:07.641885Z"),
		},
		{
			Name:    "example.dedyn.io.",
			Domain:  "example.dedyn.io",
			SubName: "",
			Type:    "NS",
			Records: []string{"ns1.desec.io.", "ns2.desec.org."},
			TTL:     3600,
			Created: mustParseTime("2020-05-06T11:46:07.641885Z"),
		},
	}
	assert.Equal(t, expected, records)
}

func mustParseTime(value string) *time.Time {
	date, _ := time.Parse(time.RFC3339, value)
	return &date
}
