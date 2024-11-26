package desec

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseCursor(t *testing.T) {
	testCases := []struct {
		desc     string
		header   string
		expected *Cursors
	}{
		{
			desc:     "all cursors",
			header:   `<https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=>; rel="first", <https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=:prev_cursor>; rel="prev", <https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=:next_cursor>; rel="next"`,
			expected: &Cursors{First: "", Prev: ":prev_cursor", Next: ":next_cursor"},
		},
		{
			desc:     "first page",
			header:   `<https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=>; rel="first", <https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=:next_cursor>; rel="next"`,
			expected: &Cursors{First: "", Prev: "", Next: ":next_cursor"},
		},
		{
			desc:     "last page",
			header:   `<https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=>; rel="first", <https://desec.io/api/v1/domains/{domain}/rrsets/?cursor=:prev_cursor>; rel="prev"`,
			expected: &Cursors{First: "", Prev: ":prev_cursor", Next: ""},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			h := http.Header{}
			h.Set("Link", test.header)

			cursor, err := parseCursor(h)
			require.NoError(t, err)

			require.Equal(t, test.expected, cursor)
		})
	}
}
