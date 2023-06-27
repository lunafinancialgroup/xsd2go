package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocomply/xsd2go/pkg/xsd2go"
	"github.com/urfave/cli"
)

// Execute ...
func Execute() error {
	app := cli.NewApp()
	app.Name = "GoComply XSD2Go"
	app.Usage = "Automatically generate golang xml parser based on XSD"
	app.Commands = []cli.Command{
		convert,
	}

	return app.Run(os.Args)
}

var convert = cli.Command{
	Name:      "convert",
	Usage:     "convert XSD to golang code to parse xml files generated by given xsd",
	ArgsUsage: "XSD-FILE GO-MODULE-IMPORT OUTPUT-DIR",
	Before: func(c *cli.Context) error {
		if c.NArg() != 3 {
			return cli.NewExitError("Exactly 3 arguments are required", 1)
		}

		for _, override := range c.StringSlice("xmlns-override") {
			if !strings.Contains(override, "=") {
				return cli.NewExitError(
					fmt.Sprintf("Invalid xmlns-override: '%s', expecting form of XMLNS=GOPKGNAME", override),
					1)
			}
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		err := xsd2go.Convert(xsd2go.Params{
			XsdPath:         c.Args()[0],
			GoModuleImport:  c.Args()[1],
			OutputDir:       c.Args()[2],
			OutputFile:      c.String("output-file"),
			TemplatePackage: c.String("template-package"),
			TemplateName:    c.String("template-name"),
			XmlnsOverrides:  c.StringSlice("xmlns-override"),
		})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output-file",
			Usage: "Defines the file name to generate code to. Example: --output-file='sample.go'",
		},
		cli.StringFlag{
			Name:  "template-package",
			Usage: "Defines where the templates exist in the packaged application. Use \".\" for root-level package. Example: --template-package='rtp'",
		},
		cli.StringFlag{
			Name:  "template-name",
			Usage: "Defines template to use for the packaged application. Example: --template-name='types.tmpl'",
		},
		cli.StringSliceFlag{
			Name:  "xmlns-override",
			Usage: "Allows to explicitly set gopackage name for given XMLNS. Example: --xmlns-override='http://www.w3.org/2000/09/xmldsig#=xml_signatures'",
		},
	},
}
