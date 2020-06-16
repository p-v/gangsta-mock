package cmd

import (
	"gangsta-mock/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(serverCmd)
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start gangsta-mock server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server commands",
}
