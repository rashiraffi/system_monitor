/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rashiraffi/system_monitor.git/internal"
	"github.com/spf13/cobra"
)

// startMonitorCmd represents the startMonitor command
var startMonitorCmd = &cobra.Command{
	Use:   "startMonitor",
	Short: "Start monitoring the system resources",
	Long:  `Start monitoring the system resources and store the data in a database.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.StartMonitor()
	},
}

func init() {
	rootCmd.AddCommand(startMonitorCmd)
}
