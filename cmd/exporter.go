/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"plimit/pkg/limitmgr"
	"time"
)

var (
	currentConnectionsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "plimit_current_connections",
		Help: "Number of currently available connections",
	})
	limitGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "plimit_current_limit",
		Help: "Currently configured limit",
	})
	hardLimitGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "plimit_haproxy_hard_limit",
		Help: "Currently configured hard limit for the haproxy autoscaler",
	})
	maxLoadGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "plimit_haproxy_max_load",
		Help: "Currently configured max load percentage of the haproxy autoscaler",
	})
)

func recordMetrics(ctx context.Context, mgr *limitmgr.LimitManager) {
	mgr.CollectGarbage(ctx)
	currentConnectionsGauge.Set(float64(mgr.GetCurrentConnectionCount(ctx)))
	limitGauge.Set(float64(mgr.GetLimit(ctx)))
	hardLimitGauge.Set(float64(mgr.GetAutoscaleHardLimit(ctx)))
	maxLoadGauge.Set(float64(mgr.GetAutoscaleMaxLoad(ctx)))
}

// exporterCmd represents the exporter command
var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "Prometheus exporter",
	Long:  `Prometheus exporter`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr := limitmgr.NewLimitManagerFromViper()
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(10 * time.Second):
					break
				}

				recordMetrics(ctx, mgr)
			}
		}()

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(viper.GetString("listen"), nil)

		cancel()
	},
}

func init() {
	flags := exporterCmd.Flags()
	flags.String("listen", ":9119", "Listen address for the exporter")
	viper.BindPFlag("listen", flags.Lookup("listen"))

	rootCmd.AddCommand(exporterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exporterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exporterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
