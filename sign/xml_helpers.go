package sign

import (
	"errors"
	"strings"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

var errXMLRoot = errors.New("xml root not found")

func canonicalize(el *etree.Element) ([]byte, error) {
	c14n := dsig.MakeC14N11Canonicalizer()
	return c14n.Canonicalize(el)
}

func findFirstElementWithAncestors(root *etree.Element, match func(*etree.Element) bool) (*etree.Element, []*etree.Element) {
	if root == nil {
		return nil, nil
	}

	var result *etree.Element
	var path []*etree.Element
	var stack []*etree.Element

	var walk func(*etree.Element)
	walk = func(el *etree.Element) {
		if el == nil || result != nil {
			return
		}

		stack = append(stack, el)

		if match(el) {
			result = el
			path = append([]*etree.Element(nil), stack...)
			stack = stack[:len(stack)-1]
			return
		}

		for _, c := range el.ChildElements() {
			walk(c)
			if result != nil {
				stack = stack[:len(stack)-1]
				return
			}
		}

		stack = stack[:len(stack)-1]
	}

	walk(root)
	return result, path
}

func addInScopeNamespaces(target *etree.Element, ancestors []*etree.Element) {
	if target == nil || len(ancestors) == 0 {
		return
	}

	inScope := map[string]string{}

	for _, el := range ancestors {
		for _, a := range el.Attr {
			if a.Space == "xmlns" {
				prefix := a.Key
				inScope[prefix] = a.Value
				continue
			}
			if a.Space == "" && a.Key == "xmlns" {
				inScope[""] = a.Value
			}
		}
	}

	hasDecl := func(prefix string) bool {
		for _, a := range target.Attr {
			if prefix == "" && a.Space == "" && a.Key == "xmlns" {
				return true
			}
			if prefix == "" && a.Space == "xmlns" && a.Key == "" {
				return true
			}
			if a.Space == "xmlns" && a.Key == prefix {
				return true
			}
		}
		return false
	}

	for prefix, uri := range inScope {
		if hasDecl(prefix) {
			continue
		}

		if prefix == "" {
			target.Attr = append(target.Attr, etree.Attr{
				Key:   "xmlns",
				Value: uri,
			})
			continue
		}

		target.Attr = append(target.Attr, etree.Attr{
			Space: "xmlns",
			Key:   prefix,
			Value: uri,
		})
	}
}

func removeUBLExtensions(root *etree.Element) {
	for i := 0; i < len(root.Child); {
		child, ok := root.Child[i].(*etree.Element)
		if ok && child.Tag == "UBLExtensions" {
			root.RemoveChildAt(i)
			continue
		}
		i++
	}
}

func removeSignature(root *etree.Element) {
	for i := 0; i < len(root.Child); {
		child, ok := root.Child[i].(*etree.Element)
		if ok && child.Tag == "Signature" {
			root.RemoveChildAt(i)
			continue
		}
		i++
	}
}

func removeQRReference(root *etree.Element) {
	for i := 0; i < len(root.Child); {
		child, ok := root.Child[i].(*etree.Element)
		if ok && child.Tag == "AdditionalDocumentReference" && isQRReference(child) {
			root.RemoveChildAt(i)
			continue
		}
		i++
	}
}

func isQRReference(el *etree.Element) bool {
	for _, c := range el.ChildElements() {
		if c.Tag == "ID" && strings.TrimSpace(c.Text()) == "QR" {
			return true
		}
	}
	return false
}
