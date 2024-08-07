package cmd

import (
	"log"
	"os"

	"github.com/daguflow/dagu/internal/config"
	"github.com/daguflow/dagu/internal/frontend"
	"github.com/daguflow/dagu/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		// nolint:line-length-limit
		Long: `dagu server [--dags=<DAGs dir>] [--host=<host>] [--port=<port>]`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("host", cmd.Flags().Lookup("host"))
			_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))
		},
		Run: func(cmd *cobra.Command, _ []string) {
			cfg, err := config.Load()
			if err != nil {
				// nolint
				log.Fatalf("Configuration load failed: %v", err)
			}
			logger := logger.NewLogger(logger.NewLoggerArgs{
				LogLevel:  cfg.LogLevel,
				LogFormat: cfg.LogFormat,
			})

			logger.Info("Server initialization", "host", cfg.Host, "port", cfg.Port)

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, logger)
			server := frontend.New(cfg, logger, cli)
			if err := server.Serve(cmd.Context()); err != nil {
				// nolint
				logger.Error("Server initialization failed", "error", err)
				os.Exit(1)
			}
		},
	}

	bindServerCommandFlags(cmd)
	return cmd
}

func bindServerCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"dags", "d", "", "location of DAG files (default is $HOME/.config/dagu/dags)",
	)
	cmd.Flags().StringP("host", "s", "", "server host (default is localhost)")
	cmd.Flags().StringP("port", "p", "", "server port (default is 8080)")
}
