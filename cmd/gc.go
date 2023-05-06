/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"plimit/pkg/limitmgr"
)

// gcCmd represents the gc command
var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Collect expired locks.",
	Long:  `Collect expired locks.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := limitmgr.NewLimitManagerFromViper()
		ctx, cancel := context.WithCancel(context.Background())
		mgr.CollectGarbage(ctx)
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(gcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
