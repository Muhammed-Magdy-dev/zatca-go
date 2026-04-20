package sign

import (
	"errors"

	"github.com/beevik/etree"
)

func HashSignedPropertiesXML(xmlData []byte) (string, error) {
	data, err := extractSignedProperties(xmlData)
	if err != nil {
		return "", err
	}

	return HashBytesBase64Hex(data), nil
}

func extractSignedProperties(xmlData []byte) ([]byte, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, err
	}

	root := doc.Root()
	if root == nil {
		return nil, errXMLRoot
	}

	el, path := findFirstElementWithAncestors(root, func(el *etree.Element) bool {
		if el == nil || el.Tag != "SignedProperties" {
			return false
		}
		for _, a := range el.Attr {
			if a.Key == "Id" && a.Value == "xadesSignedProperties" {
				return true
			}
		}
		return false
	})

	if el == nil {
		return nil, errors.New("SignedProperties not found")
	}

	copied := el.Copy()
	addInScopeNamespaces(copied, path)

	return canonicalize(copied)
}
