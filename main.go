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
	state := cmd.Execute()
	state.PrintResponse()
	ux.Close()
	os.Exit(state.ExitCode)
}

/*
@TODO - Offer prompt when finding a file during link creation.

@TODO - Consider using https://github.com/vektra/gitreader (or variant), to store repo binaries in GitHub.

*/

