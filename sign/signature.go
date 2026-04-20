package sign

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/Muhammed-Magdy-dev/zatca-go/cert"
	"github.com/beevik/etree"
)

func BuildSignatureValue(xmlData []byte, privateKeyPEM []byte) (string, error) {
	key, err := cert.ParseECPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}

	signedInfo, err := ExtractSignedInfoXML(xmlData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(signedInfo)

	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

func ExtractSignedInfoXML(xmlData []byte) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, err
	}

	root := doc.Root()
	if root == nil {
		return nil, errXMLRoot
	}

	signedInfo, path := findFirstElementWithAncestors(root, func(el *etree.Element) bool {
		return el.Space == "ds" && el.Tag == "SignedInfo"
	})

	if signedInfo == nil {
		return nil, errors.New("SignedInfo not found")
	}

	copied := signedInfo.Copy()
	addInScopeNamespaces(copied, path)

	return canonicalize(copied)
}
