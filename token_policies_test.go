package desec

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenPoliciesService_Get(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /auth/tokens/aaa/policies/rrsets/",
		fromFixtures("tokens_policy_get_all.json", http.StatusOK))

	tokens, err := client.TokenPolicies.Get(context.Background(), "aaa")
	require.NoError(t, err)

	expected := []TokenPolicy{
		{
			ID:              "7aed3f71-bc81-4f7e-90ae-8f0df0d1c211",
			Domain:          Pointer("example.com"),
			SubName:         Pointer("testing"),
			WritePermission: false,
		},
		{
			ID:              "fa6fdf60-6546-4cee-9168-5d144fe9339c",
			Domain:          Pointer("example.com"),
			SubName:         Pointer("testing"),
			Type:            Pointer("A"),
			WritePermission: true,
		},
	}
	assert.Equal(t, expected, tokens)
}

func TestTokenPoliciesService_GetOne(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /auth/tokens/aaa/policies/rrsets/fa6fdf60-6546-4cee-9168-5d144fe9339c/",
		fromFixtures("tokens_policy_get.json", http.StatusOK))

	tokens, err := client.TokenPolicies.GetOne(context.Background(), "aaa", "fa6fdf60-6546-4cee-9168-5d144fe9339c")
	require.NoError(t, err)

	expected := &TokenPolicy{
		ID:              "fa6fdf60-6546-4cee-9168-5d144fe9339c",
		Domain:          Pointer("example.com"),
		SubName:         Pointer("testing"),
		Type:            Pointer("A"),
		WritePermission: true,
	}
	assert.Equal(t, expected, tokens)
}

func TestTokenPoliciesService_Update(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("PATCH /auth/tokens/aaa/policies/rrsets/fa6fdf60-6546-4cee-9168-5d144fe9339c/",
		fromFixtures("tokens_policy_update.json", http.StatusOK))

	tokens, err := client.TokenPolicies.Update(context.Background(), "aaa", "fa6fdf60-6546-4cee-9168-5d144fe9339c", TokenPolicy{
		WritePermission: false,
	})
	require.NoError(t, err)

	expected := &TokenPolicy{
		ID:              "fa6fdf60-6546-4cee-9168-5d144fe9339c",
		Domain:          Pointer("example.com"),
		SubName:         Pointer("testing"),
		Type:            Pointer("A"),
		WritePermission: false,
	}
	assert.Equal(t, expected, tokens)
}

func TestTokenPoliciesService_GetAll(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("GET /auth/tokens/aaa/policies/rrsets/", fromFixtures("tokens_policy_get_all.json", http.StatusOK))

	tokens, err := client.TokenPolicies.GetAll(context.Background(), "aaa")
	require.NoError(t, err)

	expected := []TokenPolicy{
		{
			ID:              "7aed3f71-bc81-4f7e-90ae-8f0df0d1c211",
			Domain:          Pointer("example.com"),
			SubName:         Pointer("testing"),
			WritePermission: false,
		},
		{
			ID:              "fa6fdf60-6546-4cee-9168-5d144fe9339c",
			Domain:          Pointer("example.com"),
			SubName:         Pointer("testing"),
			Type:            Pointer("A"),
			WritePermission: true,
		},
	}
	assert.Equal(t, expected, tokens)
}

func TestTokenPoliciesService_Create(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("POST /auth/tokens/aaa/policies/rrsets/",
		fromFixtures("tokens_policy_create.json", http.StatusCreated))

	mux.HandleFunc("POST /auth/tokens/bbb/policies/rrsets/",
		fromFixtures("tokens_policy_create_empty.json", http.StatusCreated))

	testCases := []struct {
		desc     string
		tokenID  string
		policy   TokenPolicy
		expected *TokenPolicy
	}{
		{
			desc:    "all fields",
			tokenID: "aaa",
			policy: TokenPolicy{
				Domain:          Pointer("example.com"),
				SubName:         Pointer("testing"),
				Type:            Pointer("A"),
				WritePermission: true,
			},
			expected: &TokenPolicy{
				ID:              "2f133e8e-56a0-4b19-8e7e-f2e29c7ce263",
				Domain:          Pointer("example.com"),
				SubName:         Pointer("testing"),
				Type:            Pointer("A"),
				WritePermission: true,
			},
		},
		{
			desc:    "all null JSON policy",
			tokenID: "bbb",
			policy:  TokenPolicy{},
			expected: &TokenPolicy{
				ID:              "a563a574-33c9-45d1-9201-e5577b42aaf1",
				WritePermission: false,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			policies, err := client.TokenPolicies.Create(context.Background(), test.tokenID, test.policy)
			require.NoError(t, err)

			assert.Equal(t, test.expected, policies)
		})
	}
}

func TestTokenPoliciesService_Delete(t *testing.T) {
	client, mux := setupTest(t, "token")

	mux.HandleFunc("DELETE /auth/tokens/aaa/policies/rrsets/bbb/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	})

	err := client.TokenPolicies.Delete(context.Background(), "aaa", "bbb")
	require.NoError(t, err)
}
