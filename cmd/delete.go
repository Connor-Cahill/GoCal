package cmd

import (
	"strings"

	calevents "github.com/connor-cahill/goCal/services"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Removes specific event from calendar.",
	Run: func(cmd *cobra.Command, args []string) {
		index := strings.Join(args, " ")
		calevents.DeleteItem(index)
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(deleteCmd)
}
