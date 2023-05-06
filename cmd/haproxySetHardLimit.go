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

// setHardLimitCmd represents the setHardLimit command
var setHardLimitCmd = &cobra.Command{
	Use:        "set-hard-limit [flags] new-limit",
	Short:      "Set the hard limit beyond which the autoscaler will not go",
	Long:       `Set the hard limit beyond which the autoscaler will not go`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"new-limit"},
	Run: func(cmd *cobra.Command, args []string) {
		newLimit, err := strconv.Atoi(args[0])

		if err != nil {
			log.Panicf("Failed to parse new limit: %e\n", err)
		}

		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		mgr.SetAutoscaleHardLimit(ctx, int64(newLimit))

		cancelGlobal()
	},
}

func init() {
	haproxyCmd.AddCommand(setHardLimitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setHardLimitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setHardLimitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
