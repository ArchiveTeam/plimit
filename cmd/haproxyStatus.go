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
var haproxyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Current status for the haproxy autoscaler.",
	Long:  `Current status for the haproxy autoscaler.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelGlobal := context.WithCancel(context.Background())

		mgr := limitmgr.NewLimitManagerFromViper()

		w := viper.GetBool("h-watch")

		for {
			hardLimit := mgr.GetAutoscaleHardLimit(ctx)
			maxLoad := mgr.GetAutoscaleMaxLoad(ctx)
			log.Printf("Hard limit: %v\tMax load: %v%%\n", hardLimit, maxLoad)

			if !w {
				break
			}

			time.Sleep(1 * time.Second)
		}

		cancelGlobal()
	},
}

func init() {
	flags := haproxyStatusCmd.Flags()
	flags.BoolP("watch", "w", false, "Loop")
	viper.BindPFlag("h-watch", flags.Lookup("watch"))

	haproxyCmd.AddCommand(haproxyStatusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
