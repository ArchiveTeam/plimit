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

// haproxySetMinLimitCmd represents the haproxySetMinLimit command
var haproxySetMinLimitCmd = &cobra.Command{
	Use:   "set-min-limit [flags] <new-limit>",
	Short: "Set the lower bound for the autoscale limit",
	Long:  `Set the lower bound for the autoscale limit`,
	Run: func(cmd *cobra.Command, args []string) {
		newLimit, err := strconv.Atoi(args[0])

		if err != nil {
			log.Panicf("Failed to parse new limit: %e\n", err)
		}

		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		mgr.SetAutoscaleMinLimit(ctx, int64(newLimit))

		cancelGlobal()
	},
}

func init() {
	haproxyCmd.AddCommand(haproxySetMinLimitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// haproxySetMinLimitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// haproxySetMinLimitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
