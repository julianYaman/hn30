package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

var proxyClient = &http.Client{
	// Prevent the client from following redirects automatically, as we need to validate
	// the redirect location to prevent SSRF.
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

const maxImageSize = 10 * 1024 * 1024 // 10 MB

// isPrivateIP checks if an IP address is in a private, loopback, or link-local range.
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7", // Unique local addresses
	}

	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// validateURL performs security checks on the URL to be proxied.
func validateURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme: %s", parsedURL.Scheme)
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return nil, fmt.Errorf("could not resolve hostname: %s", parsedURL.Hostname())
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return nil, fmt.Errorf("denied: URL resolves to a private IP address")
		}
	}

	return parsedURL, nil
}

func imageProxyHandler(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("url")
	if imageURL == "" {
		http.Error(w, "Missing image URL", http.StatusBadRequest)
		return
	}

	// 1. Validate the URL to prevent SSRF
	validatedURL, err := validateURL(imageURL)
	if err != nil {
		LogError("Proxy: SSRF validation failed for %s: %v", imageURL, err)
		http.Error(w, fmt.Sprintf("Invalid or forbidden URL: %v", err), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", validatedURL.String(), nil)
	if err != nil {
		LogError("Proxy: Failed to create request for %s: %v", imageURL, err)
		http.Error(w, "Failed to create image request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", customUserAgent)

	resp, err := proxyClient.Do(req)
	if err != nil {
		LogError("Proxy: Failed to fetch image %s: %v", imageURL, err)
		http.Error(w, "Failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 2. Enforce a size limit to prevent DoS
	if resp.ContentLength > maxImageSize {
		LogError("Proxy: Image %s exceeds max size of %d bytes", imageURL, maxImageSize)
		http.Error(w, "Image exceeds maximum size", http.StatusRequestEntityTooLarge)
		return
	}

	// 3. Use a header allowlist instead of copying all headers
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	if resp.ContentLength > 0 {
		w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Stream the image body to the client, respecting the size limit as a fallback
	ltdReader := &io.LimitedReader{R: resp.Body, N: maxImageSize}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, ltdReader)
}
