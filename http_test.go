package utils

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
)

// Make sure copy does it right, so the copied url
// is safe to alter without modifying the other
func TestCopyUrl(t *testing.T) {
	urlA := &url.URL{
		Scheme:   "http",
		Host:     "localhost:5000",
		Path:     "/upstream",
		Opaque:   "opaque",
		RawQuery: "a=1&b=2",
		Fragment: "#hello",
		User:     &url.Userinfo{},
	}
	urlB := CopyURL(urlA)
	st.Expect(t, urlB, urlA)
	urlB.Scheme = "https"
	st.Reject(t, urlB, urlA)
}

// Make sure copy headers is not shallow and copies all headers
func TestCopyHeaders(t *testing.T) {
	source, destination := make(http.Header), make(http.Header)
	source.Add("a", "b")
	source.Add("c", "d")

	CopyHeaders(destination, source)

	st.Expect(t, destination.Get("a"), "b")
	st.Expect(t, destination.Get("c"), "d")

	// make sure that altering source does not affect the destination
	source.Del("a")
	st.Expect(t, source.Get("a"), "")
	st.Expect(t, destination.Get("a"), "b")
}

func TestHasHeaders(t *testing.T) {
	source := make(http.Header)
	source.Add("a", "b")
	source.Add("c", "d")
	st.Expect(t, HasHeaders([]string{"a", "f"}, source), true)
	st.Expect(t, HasHeaders([]string{"i", "j"}, source), false)
}

func TestRemoveHeaders(t *testing.T) {
	source := make(http.Header)
	source.Add("a", "b")
	source.Add("a", "m")
	source.Add("c", "d")
	RemoveHeaders(source, "a")
	st.Expect(t, source.Get("a"), "")
	st.Expect(t, source.Get("c"), "d")
}

func TestIsWebSocketRequest(t *testing.T) {
	headers := make(http.Header)
	headers.Set("Connection", "upgrade")
	headers.Set("Upgrade", "websocket")
	req := &http.Request{Header: headers}
	st.Expect(t, IsWebsocketRequest(req), true)

	headers = make(http.Header)
	headers.Set("Connection", "keep-alive")
	req = &http.Request{Header: headers}
	st.Expect(t, IsWebsocketRequest(req), false)
}

func TestConstainsHeader(t *testing.T) {
	tests := []struct {
		header  string
		value   string
		matches bool
	}{
		{"foo", "bar", true},
		{"bar", "foo", true},
		{"Baz", "baz", false},
		{"foo", "foo", false},
		{"foo", "", false},
		{"", "", false},
	}

	headers := make(http.Header)
	headers.Set("Foo", "bar")
	headers.Set("bar", "foo")
	req := &http.Request{Header: headers}

	for _, test := range tests {
		st.Expect(t, ConstainsHeader(req, test.header, test.value), test.matches)
	}
}
