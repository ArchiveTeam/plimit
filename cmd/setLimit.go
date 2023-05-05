/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

// setLimitCmd represents the setLimit command
var setLimitCmd = &cobra.Command{
	Use:   "set-limit",
	Short: "Set the limit",
	Long:  `A`,
	Run: func(cmd *cobra.Command, args []string) {
		redisConnString := viper.GetString("redis_url")
		opt, err := redis.ParseURL(redisConnString)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %e\n", err)
		}

		if len(args) < 1 {
			log.Fatalln("Not enough args")
		}

		newLimit, err := strconv.Atoi(args[0])

		if err != nil {
			log.Fatalf("Failed to parse new limit: %e\n", err)
		}

		rdb := redis.NewClient(opt)

		ctx, cancelGlobal := context.WithCancel(context.Background())

		err = rdb.Set(ctx, "limiter:limit", newLimit, 0).Err()
		if err != nil {
			log.Fatalf("Failed to set limit: %e\n", err)
		}

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
