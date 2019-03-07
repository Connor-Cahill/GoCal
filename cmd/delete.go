package cmd

import (
	"strings"
    "strconv"
    "log"
	calevents "github.com/connor-cahill/goCal/services"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Removes specific event from calendar.",
	Run: func(cmd *cobra.Command, args []string) {
		indexStr := strings.Join(args, " ")
        index, err := strconv.Atoi(indexStr)
        if err != nil {
            log.Fatalln(err)
        }
		calevents.DeleteItem(index)
	},
}

//	Adds show command to root command
func init() {
	RootCmd.AddCommand(deleteCmd)
}
