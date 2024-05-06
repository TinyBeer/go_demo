package main

import (
	"learn_cobra/cmd"

	"github.com/spf13/cobra"
)

func main() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
