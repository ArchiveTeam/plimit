/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"plimit/pkg/limitmgr"
	"time"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays statistics about the locks",
	Long:  `Displays statistics about the locks`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		w := viper.GetBool("watch")

		for {

			limit := mgr.GetLimit(ctx)
			active := mgr.GetCurrentConnectionCount(ctx)

			log.Printf("Active locks: %v / %v", active, limit)

			if !w {
				break
			}

			time.Sleep(1 * time.Second)
		}

		cancelGlobal()
	},
}

func init() {
	flags := statusCmd.Flags()
	flags.BoolP("watch", "w", false, "Loop")
	viper.BindPFlag("watch", flags.Lookup("watch"))

	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
