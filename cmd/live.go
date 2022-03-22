package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Check for liveness file.",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("liveness-path")
		if err != nil {
			log.Fatal().Msg("Could not get liveness-path flag")
			os.Exit(1)
		}

		if _, err := os.Stat(path); err == nil {
			log.Trace().Msg("Liveness probe successful")
			os.Exit(0)
		} else if errors.Is(err, os.ErrNotExist) {
			log.Warn().Msg("Liveness file does not exist")
			os.Exit(1)
		} else {
			log.Error().Msgf("Liveness file is in a super-position: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)
}
