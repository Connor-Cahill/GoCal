package cmd

import (
	"strings"

	calevents "github.com/connor-cahill/goCal/services"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Prints detailed information about inputted event",
	Run: func(cmd *cobra.Command, args []string) {
		event := strings.Join(args, " ")
		// fmt.Println(event)
		calevents.Find(event)
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(findCmd)
}
