package main

import (
	"github.com/lunafinancialgroup/xsd2go/cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
