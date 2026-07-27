package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/lunafinancialgroup/xsd2go/pkg/xsd"
)

func GenerateTypes(schema *xsd.Schema, outputDir string, outputFile string, templateName string) error {
	t, err := newTemplate(templateName)
	if err != nil {
		return err
	}

	packageName := schema.GoPackageName()
	dir := filepath.Join(outputDir, packageName)
	err = os.MkdirAll(dir, os.FileMode(0722))
	if err != nil {
		return err
	}
	goFile := filepath.Clean(filepath.Join(dir, outputFile))
	fmt.Printf("\tGenerating '%s'\n", goFile)
	f, err := os.Create(goFile)
	if err != nil {
		return fmt.Errorf("Could not create '%s': %v", goFile, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, schema); err != nil {
		return fmt.Errorf("Could not execute template: %v", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Could not gofmt output file\nError was: '%v'\nFile was:\n%s\n", err, buf.String())
	}

	_, err = f.Write(p)
	if err != nil {
		return err
	}

	return nil
}

// lengthOnlyTextTypes are the shared text types whose fields are validated for
// length only, without the 2026 CPMI charset. A field lands here when its XSD
// simpleType is the plain (non-"_CPMI") variant: every text field in gen1
// (2025), plus the length-only text fields in gen2 messages whose 2026 schema
// leaves them charset-unrestricted (e.g. pacs.002 MsgId, head.001 Id). The
// "...CPMI" GoTypeNames are deliberately absent, so a field generated from a
// "_CPMI" simpleType falls through to Validate (length + CPMI charset).
var lengthOnlyTextTypes = map[string]bool{
	"Max4Text": true, "Max10Text": true, "Max16Text": true, "Max34Text": true,
	"Max35Text": true, "Max70Text": true, "Max105Text": true, "Max140Text": true,
	"Max350Text": true, "Max500Text": true, "Max2048Text": true,
}

// validateFn returns the validation method name to call on a field of type
// goTypeName: ValidateV1 (length only) for a plain shared text type, Validate
// (length + CPMI charset) for the "_CPMI" text types and every other type.
func validateFn(modulesPath, goTypeName string) string {
	if lengthOnlyTextTypes[goTypeName] {
		return "ValidateV1"
	}
	return "Validate"
}

func newTemplate(templateName string) (*template.Template, error) {
	in, err := getFile(templateName)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	tempText, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	return template.New(templateName).Funcs(template.FuncMap{
		"title":      cases.Title(language.AmericanEnglish).String,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"split":      strings.Split,
		"validateFn": validateFn,
	}).Parse(string(tempText))
}

// getFile returns a fs.File either using pkger or the OS. This allows for templates outside the packaged program to be used.
func getFile(templateName string) (fs.File, error) {
	return os.Open(templateName)
}
