package transport

import (
	"crypto/tls"
	"github.com/miekg/dns"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func httpTransport() *HTTP {
	return &HTTP{
		Server:    "https://cloudflare-dns.com/dns-query",
		TLSConfig: &tls.Config{},
		UserAgent: "",
		Method:    http.MethodGet,
		Timeout:   2 * time.Second,
		HTTP3:     false,
		NoPMTUd:   false,
	}
}

func TestTransportHTTP3(t *testing.T) {
	tp := httpTransport()
	tp.HTTP3 = true
	reply, err := tp.Exchange(validQuery())
	assert.Nil(t, err)
	assert.Greater(t, len(reply.Answer), 0)
}

func TestTransportHTTPInvalidResolver(t *testing.T) {
	tp := httpTransport()
	tp.Server = "https://example.com"
	_, err := tp.Exchange(validQuery())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unpacking DNS response")
}

func TestTransportHTTPServerError(t *testing.T) {
	go func() {
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}))
	}()

	tp := httpTransport()
	tp.Server = "http://localhost:8080"
	_, err := tp.Exchange(validQuery())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "got status code 500")
}

func TestTransportHTTPIDMismatch(t *testing.T) {
	go func() {
		http.ListenAndServe(":8085", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msg := dns.Msg{}
			msg.Id = 1
			buf, err := msg.Pack()
			if err != nil {
				t.Errorf("error packing DNS message: %s", err)
				return
			}
			w.Write(buf)
		}))
	}()
	tp := httpTransport()
	tp.Server = "http://localhost:8085"
	query := validQuery()
	reply, err := tp.Exchange(query)
	assert.Nil(t, err)
	assert.Equal(t, uint16(1), reply.Id)
	assert.NotEqual(t, 1, query.Id)
}
