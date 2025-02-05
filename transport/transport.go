package transport

import "github.com/miekg/dns"

type Transport interface {
	Exchange(*dns.Msg) (*dns.Msg, error)
	Close() error
}

type Type string

const (
	TypePlain    Type = "plain"
	TypeTCP      Type = "tcp"
	TypeTLS      Type = "tls"
	TypeHTTP     Type = "http"
	TypeQUIC     Type = "quic"
	TypeDNSCrypt Type = "dnscrypt"
)

// Interface guards
var (
	_ Transport = (*Plain)(nil)
	_ Transport = (*TLS)(nil)
	_ Transport = (*HTTP)(nil)
	_ Transport = (*ODoH)(nil)
	_ Transport = (*QUIC)(nil)
	_ Transport = (*DNSCrypt)(nil)
)
