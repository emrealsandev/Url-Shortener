package security

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	ErrInvalidUrl        = errors.New("invalid_url")
	ErrUnsupportedScheme = errors.New("unsupported_scheme")
)

func NormalizeUrl(s string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil || u.Scheme == "" && u.Host == "" {
		return "", ErrInvalidUrl
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrUnsupportedScheme
	}

	host := strings.ToLower(u.Hostname())
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return "", ErrInvalidUrl
		}
	}

	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	return u.String(), nil
}
