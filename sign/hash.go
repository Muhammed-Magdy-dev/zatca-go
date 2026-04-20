package sign

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

func HashBytesBase64(data []byte) string {
	sum := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func HashBytesBase64Hex(data []byte) string {
	return HashBytesBase64(data)
}

func HashInvoiceXML(xmlData []byte) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return "", err
	}

	root := doc.Root()
	if root == nil {
		return "", errors.New("xml root not found")
	}

	removeElementsForInvoiceHash(root)

	c14n := dsig.MakeC14N11Canonicalizer()
	canon, err := c14n.Canonicalize(root)
	if err != nil {
		return "", err
	}

	return HashBytesBase64(canon), nil
}

func removeElementsForInvoiceHash(root *etree.Element) {
	removeUBLExtensions(root)
	removeSignature(root)
	removeQRReference(root)
}
