package main

import (
	"github.com/gearboxworks/bootstrap/cmd"
	"github.com/gearboxworks/bootstrap/defaults"
	"github.com/newclarity/scribeHelpers/ux"
	"os"
	"strings"
)


func init() {
	_ = ux.Open(strings.ToUpper(defaults.BinaryName) + ": ")
}


func main() {
	exitCode := 0

	err := cmd.Execute()
	if err != nil {
		ux.PrintflnError("%s", err)
		exitCode = 1
	}

	ux.Close()
	os.Exit(exitCode)
}
