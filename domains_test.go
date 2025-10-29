package desec

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainsService_Create(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("POST /domains/", fromFixtures("domains_create.json", http.StatusCreated))

	newDomain, err := client.Domains.Create(context.Background(), "example.com")
	require.NoError(t, err)

	expected := &Domain{
		Name:       "example.com",
		MinimumTTL: 3600,
		Keys: []DomainKey{
			{
				DNSKey: "257 3 13 WFRl60...",
				DS: []string{
					"6006 13 1 8581e9...",
					"6006 13 2 f34b75...",
					"6006 13 3 dfb325...",
					"6006 13 4 2fdcf8...",
				},
				Flags:   257,
				KeyType: "csk",
			},
		},
		Created:   mustParseTime("2018-09-18T16:36:16.510368Z"),
		Published: mustParseTime("2018-09-18T17:21:38.348112Z"),
		Touched:   mustParseTime("2018-09-18T17:21:38.348112Z"),
	}
	assert.Equal(t, expected, newDomain)
}

func TestDomainsService_Delete(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("DELETE /domains/example.com/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.Domains.Delete(context.Background(), "example.com")
	require.NoError(t, err)
}

func TestDomainsService_Get(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /domains/example.com/", fromFixtures("domains_get.json", http.StatusOK))

	domain, err := client.Domains.Get(context.Background(), "example.com")
	require.NoError(t, err)

	expected := &Domain{
		Name:       "example.com",
		MinimumTTL: 3600,
		Keys: []DomainKey{
			{
				DNSKey: "257 3 13 WFRl60...",
				DS: []string{
					"6006 13 1 8581e9...",
					"6006 13 2 f34b75...",
					"6006 13 3 dfb325...",
					"6006 13 4 2fdcf8...",
				},
				Flags:   257,
				KeyType: "csk",
			},
		},
		Created:   mustParseTime("2018-09-18T16:36:16.510368Z"),
		Published: mustParseTime("2018-09-18T17:21:38.348112Z"),
		Touched:   mustParseTime("2018-09-18T17:21:38.348112Z"),
	}
	assert.Equal(t, expected, domain)
}

func TestDomainsService_GetResponsible(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /domains/", func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("owns_qname") != "git.dev.example.org" {
			http.Error(rw, "owns_qname not passed correctly", http.StatusBadRequest)
			return
		}

		fromFixtures("domains_getresponsible.json", http.StatusOK).ServeHTTP(rw, req)
	})

	domain, err := client.Domains.GetResponsible(context.Background(), "git.dev.example.org")
	require.NoError(t, err)

	expected := &Domain{
		Name:       "dev.example.org",
		MinimumTTL: 3600,
		Created:    mustParseTime("2022-11-12T18:01:35.454616Z"),
		Published:  mustParseTime("2022-11-12T18:03:19.516440Z"),
		Touched:    mustParseTime("2022-11-12T18:03:19.516440Z"),
	}
	assert.Equal(t, expected, domain)
}

func TestDomainsService_GetResponsible_error(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /domains/", func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte("[]"))
	})

	_, err := client.Domains.GetResponsible(context.Background(), "git.dev.example.org")

	var notFoundError *NotFoundError

	require.ErrorAs(t, err, &notFoundError)
}

func TestDomainsService_GetAll(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /domains/", fromFixtures("domains_getall.json", http.StatusOK))

	domains, err := client.Domains.GetAll(context.Background())
	require.NoError(t, err)

	expected := []Domain{
		{
			Name:       "example.org",
			MinimumTTL: 3600,
			Created:    mustParseTime("2020-05-13T11:35:40.954616Z"),
			Published:  mustParseTime("2020-05-13T12:25:19.816440Z"),
			Touched:    mustParseTime("2020-05-13T12:25:19.816440Z"),
		},
		{
			Name:       "example.dedyn.io",
			MinimumTTL: 60,
			Created:    mustParseTime("2020-05-05T23:17:36.101470Z"),
			Published:  mustParseTime("2020-05-06T12:13:06.138443Z"),
			Touched:    mustParseTime("2020-05-06T12:13:06.138443Z"),
		},
	}
	assert.Equal(t, expected, domains)
}
