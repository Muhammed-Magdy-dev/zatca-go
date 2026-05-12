package csr

import "strings"

func templateNameForEnv(env string) string {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "production", "prod", "core":
		return "ZATCA-Code-Signing"
	case "nonproduction", "non-production", "preproduction", "pre-production", "preprod", "pre-prod", "simulation", "sandbox", "developer-portal":
		return "PREZATCA-Code-Signing"
	default:
		return "PREZATCA-Code-Signing"
	}
}
