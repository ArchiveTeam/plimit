/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"plimit/pkg/limitmgr"
	"strconv"

	"github.com/spf13/cobra"
)

// setMaxLoadCmd represents the setMaxLoad command
var setMaxLoadCmd = &cobra.Command{
	Use:        "set-max-load [flags] <new-max-load>",
	Short:      "Set the new maximum load we will permit on haproxy in percentage.",
	Long:       `Set the new maximum load we will permit on haproxy in percentage.`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"new-max-load"},
	Run: func(cmd *cobra.Command, args []string) {
		newLimit, err := strconv.Atoi(args[0])

		if err != nil {
			log.Panicf("Failed to parse new max-load: %e\n", err)
		}

		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		mgr.SetAutoscaleMaxLoad(ctx, newLimit)

		cancelGlobal()
	},
}

func init() {
	haproxyCmd.AddCommand(setMaxLoadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setMaxLoadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setMaxLoadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
