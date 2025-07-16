package cmd

import (
	"EvoBot/backend/server"

	"github.com/spf13/cobra"
)

var configFile string

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/opt/evobot/conf/config.yaml", "Config file path")
}

var RootCmd = &cobra.Command{
	Use:   "EvoBot",
	Short: "EvoBot",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start(configFile)
		return nil
	},
}
