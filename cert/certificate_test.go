package cert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"testing"
)

func TestFormatIssuerNamePreservesDomainComponents(t *testing.T) {
	rawIssuer, err := asn1.Marshal(pkix.RDNSequence{
		{{Type: asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 25}, Value: "local"}},
		{{Type: asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 25}, Value: "gov"}},
		{{Type: asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 25}, Value: "extgazt"}},
		{{Type: asn1.ObjectIdentifier{2, 5, 4, 3}, Value: "PRZEINVOICESCA4-CA"}},
	})
	if err != nil {
		t.Fatalf("marshal raw issuer: %v", err)
	}

	certificate := &x509.Certificate{
		RawIssuer: rawIssuer,
		Issuer: pkix.Name{
			CommonName: "PRZEINVOICESCA4-CA",
		},
	}

	got := FormatIssuerName(certificate)
	want := "CN=PRZEINVOICESCA4-CA, DC=extgazt, DC=gov, DC=local"
	if got != want {
		t.Fatalf("unexpected issuer name\nwant: %s\ngot:  %s", want, got)
	}
}
