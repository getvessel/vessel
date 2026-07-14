package utils

import (
	"fmt"
	"os"
	"strings"
)

func GenerateSslipDomain(projectNameOrID string, hostIP string, wildcardDomain string) string {
	if wildcardDomain == "" {
		wildcardDomain = os.Getenv("VESSL_WILDCARD_DOMAIN")
	}
	if wildcardDomain != "" {
		cleanName := sanitizeDomainName(projectNameOrID)
		return fmt.Sprintf("http://%s.%s", cleanName, strings.TrimPrefix(wildcardDomain, "*."))
	}

	if hostIP == "" {
		hostIP = os.Getenv("VESSL_HOST_IP")
	}
	if hostIP == "" {
		hostIP = "127.0.0.1"
	}
	cleanIP := strings.ReplaceAll(strings.TrimSpace(hostIP), ".", "-")
	cleanName := sanitizeDomainName(projectNameOrID)
	return fmt.Sprintf("http://%s.%s.sslip.io", cleanName, cleanIP)
}

func sanitizeDomainName(name string) string {
	clean := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
	if len(clean) > 32 {
		clean = clean[:32]
	}
	return strings.Trim(clean, "-")
}
