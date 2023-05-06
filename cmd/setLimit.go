/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"plimit/pkg/limitmgr"
	"strconv"
)

// setLimitCmd represents the setLimit command
var setLimitCmd = &cobra.Command{
	Use:        "set-limit [flags] new-limit",
	Short:      "Set the limit",
	Long:       `A`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"new-limit"},
	Run: func(cmd *cobra.Command, args []string) {
		newLimit, err := strconv.Atoi(args[0])

		if err != nil {
			log.Fatalf("Failed to parse new limit: %e\n", err)
		}

		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		mgr.SetLimit(ctx, int64(newLimit))

		cancelGlobal()
	},
}

func init() {
	rootCmd.AddCommand(setLimitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setLimitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setLimitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
