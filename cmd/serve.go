package cmd

import (
	"github.com/msvens/mchess/internal/api"
	"github.com/msvens/mchess/internal/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the mchess API server to handle player data requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.SetupLogger()
		api.StartServer(cfg)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
