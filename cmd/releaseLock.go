/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
	"plimit/pkg/limitmgr"
)

// releaseLockCmd represents the releaseLock command
var releaseLockCmd = &cobra.Command{
	Use:        "release-lock [flags] lock-id",
	Short:      "Manually release a lock.",
	Long:       `Manually release a lock.`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := uuid.Parse(args[0])
		if err != nil {
			log.Panicf("Failed to parse uuid: %s\n", id)
		}

		ctx, cancel := context.WithCancel(context.Background())
		mgr := limitmgr.NewLimitManagerFromViper()
		mgr.ReleaseLock(ctx, id)
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(releaseLockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseLockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseLockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
