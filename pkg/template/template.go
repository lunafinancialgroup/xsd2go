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

// gen1CPMITextTypes are the shared text types whose Validate enforces the 2026
// CPMI charset. For gen1 (2025) output the generated code must call ValidateV1
// (length only) instead, so a pre-cutover payload is not rejected by the stricter
// 2026 charset.
var gen1CPMITextTypes = map[string]bool{
	"Max4Text": true, "Max10Text": true, "Max16Text": true, "Max34Text": true,
	"Max35Text": true, "Max70Text": true, "Max105Text": true, "Max140Text": true,
	"Max350Text": true, "Max500Text": true, "Max2048Text": true,
	"Max4TextCPMI": true, "Max10TextCPMI": true, "Max16TextCPMI": true,
	"Max34TextCPMI": true, "Max35TextCPMI": true, "Max70TextCPMI": true,
	"Max105TextCPMI": true, "Max140TextCPMI": true, "Max350TextCPMI": true,
	"Max500TextCPMI": true, "Max2048TextCPMI": true,
}

// validateFn returns the validation method name to call on a field of type
// goTypeName: ValidateV1 for a shared CPMI text type in gen1 (2025) output,
// Validate otherwise (gen2 and all non-text types).
func validateFn(modulesPath, goTypeName string) string {
	if strings.HasSuffix(modulesPath, "/gen") && gen1CPMITextTypes[goTypeName] {
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
