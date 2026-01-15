package cmd

import (
	"fmt"
	"os"

	"github.com/msvens/mchess/internal/config"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mchess",
	Short: "A caching proxy for the Swedish Chess Federation API",
	Long: `mchess is a backend service that provides caching and batch operations
for the Swedish Chess Federation (schack.se) API.

It reduces API calls, provides batch endpoints for fetching multiple players,
and serves historical rating data efficiently.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config init for version command
		if cmd.Name() == "version" {
			return nil
		}
		return config.Init(cfgFile)
	},
	// Default action: run the server
	Run: func(cmd *cobra.Command, args []string) {
		serveCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./config.yaml)")
}
