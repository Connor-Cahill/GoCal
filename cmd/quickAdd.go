package cmd

import (
    "fmt"
    "strings"

    calevents "github.com/connor-cahill/goCal/services"
    "github.com/spf13/cobra"
)

var quickAdd = &cobra.Command{
    Use: "quickAdd",
    Short: "Given a event title creates a 1 hour starting at the current time.",
    Run: func(cmd *cobra.Command, args []string) {
        eventTitle := strings.Join(args, " ")
        fmt.Println(eventTitle)
        calevents.QuickAdd(eventTitle)
    },
}


// adds the quickAdd command to root command
func init() {
    RootCmd.AddCommand(quickAdd)
}
