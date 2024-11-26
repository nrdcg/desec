package desec

import (
	"net/http"
	"net/url"

	"github.com/peterhellberg/link"
)

type Cursors struct {
	First string
	Prev  string
	Next  string
}

func parseCursor(h http.Header) (*Cursors, error) {
	links := link.ParseHeader(h)

	c := &Cursors{}

	for s, l := range links {
		uri, err := url.ParseRequestURI(l.URI)
		if err != nil {
			return nil, err
		}

		query := uri.Query()

		switch s {
		case "first":
			c.First = query.Get("cursor")
		case "prev":
			c.Prev = query.Get("cursor")
		case "next":
			c.Next = query.Get("cursor")
		}
	}

	return c, nil
}
