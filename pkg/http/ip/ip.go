package ip

import (
	"fmt"
	"net"
	"strings"
)

var ErrorHTTPInvalidIP = fmt.Errorf("invalid IP")

// Verify if provided IP is valid. Return error if it is not
func IsValidIP(ip string) error {
	if strings.TrimSpace(ip) == "" {
		return ErrorHTTPInvalidIP
	}

	if parsed := net.ParseIP(ip); parsed == nil {
		return ErrorHTTPInvalidIP
	}
	return nil
}

// Format address & port to create a valid web address
func FormatAddress(addr string, port int) string {
	// Initialize server
	address := addr
	if port > 0 {
		address += fmt.Sprintf(":%d", port)
	}
	return strings.TrimSpace(address)
}
