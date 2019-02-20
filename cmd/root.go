package cmd

import (
	"github.com/spf13/cobra"
)

//RootCmd for our Google Cal CLI
var RootCmd = &cobra.Command{
	Use:   "gocal",
	Short: "A Google Calendar CLI",
}
