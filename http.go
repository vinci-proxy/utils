package utils

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// ProxyWriter helps to capture response headers and status code
// from the ServeHTTP. It can be safely passed to ServeHTTP handler,
// wrapping the real response writer.
type ProxyWriter struct {
	Code int
	W    http.ResponseWriter
}

// StatusCode defines the status code required.
func (p *ProxyWriter) StatusCode() int {
	if p.Code == 0 {
		// per contract standard lib will set this to http.StatusOK if not set
		// by user, here we avoid the confusion by mirroring this logic
		return http.StatusOK
	}
	return p.Code
}

// Header returns http.Header.
func (p *ProxyWriter) Header() http.Header {
	return p.W.Header()
}

// Write writes the given slice of bytes in the internal buffer.
func (p *ProxyWriter) Write(buf []byte) (int, error) {
	return p.W.Write(buf)
}

// WriteHeader writes the HTTP response status code.
func (p *ProxyWriter) WriteHeader(code int) {
	p.Code = code
	p.W.WriteHeader(code)
}

// Flush flushes the cached data, if possible.
func (p *ProxyWriter) Flush() {
	if f, ok := p.W.(http.Flusher); ok {
		f.Flush()
	}
}

// WriterStub implements a http.ResponseWriter compatible interface desiged for stub during testing.
type WriterStub struct {
	Code    int
	Body    []byte
	Headers http.Header
}

// NewWriterStub creates a new WriterStub which implements a http.ResponseWriter interface.
func NewWriterStub() *WriterStub {
	return &WriterStub{Code: http.StatusOK, Headers: make(http.Header)}
}

// Header returns http.Header.
func (p *WriterStub) Header() http.Header {
	return p.Headers
}

// Write writes the given slice of bytes in the internal buffer.
func (p *WriterStub) Write(buf []byte) (int, error) {
	p.Body = append(p.Body, buf...)
	return len(p.Body), nil
}

// WriteHeader writes the HTTP response status code.
func (p *WriterStub) WriteHeader(code int) {
	p.Code = code
}

// BufferWriter represents an HTTP writable entity.
type BufferWriter struct {
	Code int
	H    http.Header
	W    io.WriteCloser
}

// NewBufferWriter creates a new writer buffer.
func NewBufferWriter(w io.WriteCloser) *BufferWriter {
	return &BufferWriter{
		W: w,
		H: make(http.Header),
	}
}

// Close closes the used WriteCloser.
func (b *BufferWriter) Close() error {
	return b.W.Close()
}

// Header returns rw.Header.
func (b *BufferWriter) Header() http.Header {
	return b.H
}

// Write writes the giben
func (b *BufferWriter) Write(buf []byte) (int, error) {
	return b.W.Write(buf)
}

// WriteHeader sets rw.Code.
func (b *BufferWriter) WriteHeader(code int) {
	b.Code = code
}

type nopWriteCloser struct {
	io.Writer
}

func (*nopWriteCloser) Close() error { return nil }

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

// CopyURL provides update safe copy by avoiding shallow copying User field
func CopyURL(i *url.URL) *url.URL {
	out := *i
	if i.User != nil {
		out.User = &(*i.User)
	}
	return &out
}

// CopyHeaders copies http headers from source to destination, it
// does not overide, but adds multiple headers
func CopyHeaders(dst, src http.Header) {
	for k, vv := range src {
		dst[k] = append([]string{}, vv...)
	}
}

// HasHeaders determines whether any of the header names is present in the http headers
func HasHeaders(names []string, headers http.Header) bool {
	for _, h := range names {
		if headers.Get(h) != "" {
			return true
		}
	}
	return false
}

// RemoveHeaders removes the header with the given names from the headers map
func RemoveHeaders(headers http.Header, names ...string) {
	for _, h := range names {
		headers.Del(h)
	}
}

// IsWebsocketRequest determines if the specified HTTP request is a websocket handshake request.
func IsWebsocketRequest(req *http.Request) bool {
	return ConstainsHeader(req, "Connection", "upgrade") && ConstainsHeader(req, "Upgrade", "websocket")
}

// ConstainsHeader checks if the given header field is present if the given HTTP request.
func ConstainsHeader(req *http.Request, name, value string) bool {
	if name == "" || value == "" {
		return false
	}
	items := strings.Split(req.Header.Get(name), ",")
	for _, item := range items {
		if value == strings.ToLower(strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}

// EnsureTransporterFinalized will ensure that when the HTTP client is GCed
// the runtime will close the idle connections (so that they won't leak)
// this function was adopted from Hashicorp's go-cleanhttp package.
func EnsureTransporterFinalized(httpTransport *http.Transport) {
	runtime.SetFinalizer(&httpTransport, func(transportInt **http.Transport) {
		(*transportInt).CloseIdleConnections()
	})
}

// DefaultTransport returns a new http.Transport with the same default values
// as http.DefaultTransport, but with idle connections and keepalives disabled.
func DefaultTransport() *http.Transport {
	transport := DefaultPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// DefaultPooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func DefaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
	}
	return transport
}
