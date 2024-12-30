/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

// rootCmd represents the base command when called without any subcommands

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if cfg.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	profile := false
	var rootCmd = &cobra.Command{
		Use:   "web",
		Short: "web is a web server",
		Long:  `web is a web server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}
			f, pErr := os.Create(("cpu.pprof"))
			if pErr != nil {
				return pErr
			}

			_ = pprof.StartCPUProfile(f)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}
			pprof.StopCPUProfile()
			f, pErr := os.Create(("mem.pprof"))
			if pErr != nil {
				return pErr
			}
			defer f.Close()

			runtime.GC()
			err := pprof.WriteHeapProfile(f)
			return err
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfg.HttpAddr, "http-addr", "a", cfg.HttpAddr, "HTTP address")
	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "enable profiling")

	rootCmd.AddCommand(APICmd(ctx, cfg))
	rootCmd.AddCommand(WorkerCmd(ctx, cfg))

	if err = rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
