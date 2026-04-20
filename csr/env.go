package csr

func templateNameForEnv(env string) string {
	switch env {
	case "production":
		return "ZATCA-Code-Signing"
	case "nonProduction":
		return "PREZATCA-Code-Signing"
	default:
		return "PREZATCA-Code-Signing"
	}
}
