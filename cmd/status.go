/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var lockCount = redis.NewScript(`
return #redis.pcall('keys', 'limiter:locks:*')
`)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays statistics about the locks",
	Long:  `Displays statistics about the locks`,
	Run: func(cmd *cobra.Command, args []string) {
		redisConnString := viper.GetString("redis_url")
		opt, err := redis.ParseURL(redisConnString)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %e\n", err)
		}

		rdb := redis.NewClient(opt)

		ctx, cancelGlobal := context.WithCancel(context.Background())

		w := viper.GetBool("watch")

		for {

			limit, err := rdb.Get(ctx, "limiter:limit").Int()
			if err != nil {
				if err == redis.Nil {
					limit = 0
				} else {
					log.Fatalf("Failed to fetch data: %e\n", err)
				}
			}

			count, err := lockCount.Run(ctx, rdb, []string{}).Int()

			if err != nil {
				log.Fatalf("Failed to fetch data: %e\n", err)
			}

			log.Printf("Active locks: %d / %d", count, limit)

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
