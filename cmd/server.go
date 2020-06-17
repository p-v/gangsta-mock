package cmd

import (
	"gangsta-mock/server"
	"github.com/spf13/cobra"
)

func init() {
	startServerCmd.Flags().StringP("config", "c", "gangsta.yml", "Configuration file")
	serverCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(serverCmd)
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start gangsta-mock server",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		server.StartServer(configFile)
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server commands",
}
