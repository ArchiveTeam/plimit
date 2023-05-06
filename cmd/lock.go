/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"plimit/pkg/limitmgr"
	"time"
)

var errorDelay = 30 * time.Second
var lockDelay = 2 * time.Second
var refreshDelay = 60 * time.Second
var lockDuration = time.Hour

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock [flags] command...",
	Short: "Run a command with a connection lock",
	Long:  `Run a command with a connection lock`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelGlobal := context.WithCancel(context.Background())
		id, err := uuid.NewUUID()
		if err != nil {
			log.Panicln(err)
		}

		mgr := limitmgr.NewLimitManagerFromViper()

		mgr.CollectGarbage(ctx)

		for {
			acquired, err := mgr.TryAcquireLock(ctx, id, lockDuration)

			if err != nil && err != redis.Nil {
				log.Printf("Error during acquire: %v\n", err)
				time.Sleep(errorDelay)
			} else {
				if !acquired {
					log.Println("No more locks available. Sleeping...")
					time.Sleep(lockDelay)
				} else {
					break
				}
			}
		}

		defer func() {
			cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			mgr.ReleaseLock(cctx, id)
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

				mgr.RefreshLock(ctx, id, lockDuration)
			}
		}()

		log.Println("Running command...")
		wrappedCmd := exec.CommandContext(ctx, args[0], args[1:]...)
		// redirect the output to terminal
		wrappedCmd.Stdout = os.Stdout
		wrappedCmd.Stderr = os.Stderr
		err = wrappedCmd.Run()
		if err != nil {
			log.Panicf("Failed to execute wrapped command: %e\n", err)
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
