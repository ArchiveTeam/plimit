/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/csv"
	"github.com/spf13/viper"
	"io"
	"log"
	"modernc.org/mathutil"
	"net/http"
	"plimit/pkg/limitmgr"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func readCSVFromUrl(url string) []map[string]string {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("Unable to fetch status: %s\n", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	rows := []map[string]string{}
	var header []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicf("Unable to parse haproxy csv: %s\n", err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}

	return rows
}

type parsedHAProxy struct {
	SessionsCurrent int
	SessionsLimit   int
}

func parseHAProxyStats(url string) *parsedHAProxy {
	data := readCSVFromUrl(url)

	for _, row := range data {
		if row["# pxname"] == "s3" && row["svname"] == "BACKEND" {
			slim, err := strconv.Atoi(row["slim"])
			if err != nil {
				log.Panicf("Unable to parse slim (%s) as number: %s", row["slim"], err)
			}
			scur, err := strconv.Atoi(row["scur"])
			if err != nil {
				log.Panicf("Unable to parse scur (%s) as number: %s", row["scur"], err)
			}
			return &parsedHAProxy{
				SessionsLimit:   slim,
				SessionsCurrent: scur,
			}
		}
	}
	log.Panicf("Unable to find totals in haproxy!\n")
	return nil
}

func autoAdjustOnce(ctx context.Context, url string, mgr *limitmgr.LimitManager) {
	log.Println("Doing autoadjustment...")

	mgr.CollectGarbage(ctx)

	log.Println("Fetching stats...")
	stats := parseHAProxyStats(url)
	log.Println("Fetched!")
	log.Printf("Current sessions: %v / %v\n", stats.SessionsCurrent, stats.SessionsLimit)

	connectionsNotOurs := stats.SessionsCurrent - int(mgr.GetCurrentConnectionCount(ctx))
	log.Printf("Connections that are not ours: %v\n", connectionsNotOurs)

	maxLoad := mgr.GetAutoscaleMaxLoad(ctx)
	log.Printf("Configured max percentage of connections we will allow: %d%%\n", maxLoad)

	allowedLimit := stats.SessionsLimit * maxLoad / 100
	log.Printf("How many connections we will permit in total: %v\n", allowedLimit)

	hardLimit := mgr.GetAutoscaleHardLimit(ctx)
	log.Printf("Configured hard limit: %v\n", hardLimit)

	connectionsAvailableForUs := mathutil.Clamp(allowedLimit-connectionsNotOurs, 0, hardLimit)

	log.Printf("We will set our limit to: %v\n", connectionsAvailableForUs)
	mgr.SetLimit(ctx, int64(connectionsAvailableForUs))

	log.Println("Autoadjustment complete!")
}

// autoadjustCmd represents the autoadjust command
var autoadjustCmd = &cobra.Command{
	Use:        "autoadjust [flags] haproxy-url",
	Short:      "Automatically adjust the limit based on haproxy.",
	Long:       `Automatically adjust the limit based on haproxy.`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"haproxy url"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		mgr := limitmgr.NewLimitManagerFromViper()
		for {
			autoAdjustOnce(ctx, args[0], mgr)

			r := viper.GetInt("repeat")

			if r == 0 {
				break
			} else {
				log.Printf("Sleeping for %v seconds...\n", r)
				time.Sleep(time.Duration(r) * time.Second)
			}
		}
		cancel()
	},
}

func init() {
	flags := autoadjustCmd.Flags()
	flags.IntP("repeat", "r", 0, "Set to a non-zero value to repeat the adjustment every n seconds.")
	viper.BindPFlag("repeat", flags.Lookup("repeat"))

	haproxyCmd.AddCommand(autoadjustCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// autoadjustCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// autoadjustCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
