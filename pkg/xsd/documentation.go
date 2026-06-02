package xsd

import "strings"

// Documentation is a single xsd:documentation node within an xsd:annotation.
type Documentation struct {
	Source string `xml:"source,attr"`
	Value  string `xml:",chardata"`
}

// definition returns the ISO 20022 "Definition" documentation from docs,
// whitespace-normalised to a single line; it returns "" when none is present.
func definition(docs []Documentation) string {
	for _, doc := range docs {
		if strings.EqualFold(doc.Source, "Definition") {
			return strings.Join(strings.Fields(doc.Value), " ")
		}
	}
	return ""
}
