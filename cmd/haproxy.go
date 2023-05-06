/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// haproxyCmd represents the haproxy command
var haproxyCmd = &cobra.Command{
	Use:   "haproxy",
	Short: "Subcommands for handling haproxy autoadjustment.",
	Long:  `Subcommands for handling haproxy autoadjustment.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("haproxy called")
	//},
}

func init() {
	rootCmd.AddCommand(haproxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// haproxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// haproxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
