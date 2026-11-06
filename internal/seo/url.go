package seo

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"strconv"
	"strings"
	"unicode"
)

// NormalizePublicSiteURL validates an explicit origin for root-level public routes.
// A subpath, credentials, query or fragment is not supported. Local origins are
// allowed for explicit development/test configuration, never inferred from Host.
func NormalizePublicSiteURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || !validHTTPURL(u) || strings.ContainsAny(raw, "?#") || strings.Trim(u.Path, "/") != "" {
		return "", fmt.Errorf("PUBLIC_SITE_URL must be an explicit http(s) origin without credentials, path, query or fragment")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path, u.RawPath = "", ""
	return u.String(), nil
}

func validHTTPURL(u *url.URL) bool {
	if u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil || u.Opaque != "" {
		return false
	}
	if strings.ContainsAny(u.Hostname(), "\\ <>\"'") {
		return false
	}
	if strings.HasSuffix(u.Host, ":") {
		return false
	}
	if net.ParseIP(u.Hostname()) == nil {
		if len(u.Hostname()) > 253 {
			return false
		}
		for _, label := range strings.Split(strings.TrimSuffix(u.Hostname(), "."), ".") {
			if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
				return false
			}
			for _, r := range label {
				if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-') {
					return false
				}
			}
		}
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return false
		}
	}
	return true
}

// AbsoluteURL joins a validated site origin to a root-level public path.
func AbsoluteURL(base, publicPath string) string { return base + publicPath }

// ArticleURL returns the canonical absolute article URL for a persisted slug.
func ArticleURL(base, slug string) string { return AbsoluteURL(base, ArticlePath(slug)) }

// PublicResourceURL accepts only safe publicly addressable http(s) resources.
// Relative root paths use the configured origin; private application/storage
// paths, credentials, signed query URLs, fragments and local network targets are
// omitted. No DNS resolution or request is performed.
func PublicResourceURL(base, raw string) string {
	if raw == "" || strings.TrimSpace(raw) != raw || strings.ContainsAny(raw, "\\") || strings.IndexFunc(raw, unicode.IsControl) >= 0 {
		return ""
	}
	if strings.HasPrefix(raw, "/") && !strings.HasPrefix(raw, "//") {
		raw = base + raw
	}
	u, err := url.Parse(raw)
	if err != nil || !validHTTPURL(u) || strings.ContainsAny(raw, "?#") {
		return ""
	}
	if strings.Contains(u.Path, "\\") || strings.IndexFunc(u.Path, unicode.IsControl) >= 0 {
		return ""
	}
	host := strings.ToLower(u.Hostname())
	host = strings.TrimSuffix(host, ".")
	if host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") || !strings.Contains(host, ".") {
		return ""
	}
	if ip := net.ParseIP(host); ip != nil && (ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast()) {
		return ""
	}
	// WHATWG hosts ending in a number are parsed as IPv4, not DNS. Reject
	// noncanonical forms rather than trying to reproduce browser radix/overflow
	// rules. This includes hex, octal, short, mixed and single-number spellings,
	// as well as invalid numeric-ending DNS names browsers cannot load as DNS.
	if net.ParseIP(host) == nil && browserNumericHost(host) {
		return ""
	}
	clean := strings.ToLower(path.Clean(u.Path))
	for _, prefix := range []string{"/uploads/knowledge", "/api", "/dashboard", "/admin", "/favorites", "/profile", "/login"} {
		if clean == prefix || strings.HasPrefix(clean, prefix+"/") {
			return ""
		}
	}
	u.Path, u.RawPath = path.Clean(u.Path), ""
	return u.String()
}

func browserNumericHost(host string) bool {
	last := host[strings.LastIndex(host, ".")+1:]
	if strings.HasPrefix(last, "0x") {
		last = last[2:]
		return strings.IndexFunc(last, func(r rune) bool { return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') }) < 0
	}
	return last != "" && strings.IndexFunc(last, func(r rune) bool { return r < '0' || r > '9' }) < 0
}
