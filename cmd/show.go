package cmd

import (
	"github.com/connor-cahill/goCal/services"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "prints all upcoming events to terminal",
	Run: func(cmd *cobra.Command, args []string) {
		services.GetEventList()
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(showCmd)
}
