/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"time"
)

var errorDelay = 30 * time.Second
var lockDelay = 2 * time.Second
var refreshDelay = 60 * time.Second

var acquireScript = redis.NewScript(`
local key = KEYS[1]

local limit = redis.call("GET", "limiter:limit")
local num_locks = #redis.pcall('keys', 'limiter:locks:*')
if limit == nil then
	limit = 0
end
if num_locks < tonumber(limit) then
	redis.call('SET', 'limiter:locks:'..key, key, 'EX', '3600')
	return true
else
	return false
end
`)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Run a command with a connection lock",
	Long:  `AAAA`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("Command missing!")
		}

		redisConnString := viper.GetString("redis_url")
		opt, err := redis.ParseURL(redisConnString)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %e\n", err)
		}

		rdb := redis.NewClient(opt)

		ctx, cancelGlobal := context.WithCancel(context.Background())

		id, err := uuid.NewUUID()

		if err != nil {
			log.Fatalln(err)
		}

		var acquired bool
		for !acquired {
			log.Println("Attempting...")
			acquired, err = acquireScript.Run(ctx, rdb, []string{id.String()}).Bool()

			if err != nil && err != redis.Nil {
				log.Printf("Error during acquire: %v\n", err)
				time.Sleep(errorDelay)
			} else {
				if !acquired {
					log.Println("No more locks available. Sleeping...")
					time.Sleep(lockDelay)
				}
			}
		}

		defer func() {
			cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			log.Printf("Releasing lock %s...\n", id.String())
			rdb.Del(cctx, fmt.Sprintf("limiter:locks:%s", id.String()))
			cancel()
		}()

		log.Printf("Acquired lock: %s\n", id.String())

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(refreshDelay):
					break
				}

				err := rdb.Expire(ctx, fmt.Sprintf("limiter:locks:%s", id.String()), 3600*time.Second).Err()
				if err != nil {
					log.Println("Refreshing lock failed: %v\n", err)
				}
			}
		}()

		log.Println("Running command...")
		wrappedCmd := exec.CommandContext(ctx, args[0], args[1:]...)
		// redirect the output to terminal
		wrappedCmd.Stdout = os.Stdout
		wrappedCmd.Stderr = os.Stderr
		err = wrappedCmd.Run()
		if err != nil {
			log.Fatalf("Failed to execute wrapped command: %e\n", err)
		}

		log.Println("Command run complete!")

		cancelGlobal()

	},
}

func init() {
	rootCmd.AddCommand(lockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
