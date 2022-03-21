package cmd

import (
	server "echoserver/internal"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var bannerText string = `
                       _/
    _/_/      _/_/_/  _/_/_/      _/_/      _/_/_/    _/_/    _/  _/_/  _/      _/    _/_/    _/  _/_/
 _/_/_/_/  _/        _/    _/  _/    _/  _/_/      _/_/_/_/  _/_/      _/      _/  _/_/_/_/  _/_/
_/        _/        _/    _/  _/    _/      _/_/  _/        _/          _/  _/    _/        _/
 _/_/_/    _/_/_/  _/    _/    _/_/    _/_/_/      _/_/_/  _/            _/        _/_/_/  _/
`

var rootCmd = &cobra.Command{
	Use:   "echoserver",
	Short: "An HTTP echo server implementation in Go",
	Long: `An HTTP echo server implementation in Go. This will echo back to you the status
code and information you send it. This is useful for testing infrastructure.`,
	Run: func(cmd *cobra.Command, args []string) {
		banner, err := cmd.Flags().GetBool("banner")
		host, err := cmd.Flags().GetString("host")
		port, err := cmd.Flags().GetInt("port")
		metricsPort, err := cmd.Flags().GetInt("metrics-port")
		wait, err := cmd.Flags().GetInt("wait")
		writeTimeout, err := cmd.Flags().GetInt("write-timeout")
		readTimeout, err := cmd.Flags().GetInt("read-timeout")
		idleTimeout, err := cmd.Flags().GetInt("idle-timeout")
		livenessPath, err := cmd.Flags().GetString("liveness-path")
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Could not get options: %s", err))
			os.Exit(1)
		}

		if banner == true {
			fmt.Println(bannerText)
		}

		err = server.RunServer(&server.Options{
			LivenessFilePath: livenessPath,
			Host:             host,
			Port:             port,
			MetricsPort:      metricsPort,
			Wait:             wait,
			WriteTimeout:     writeTimeout,
			ReadTimeout:      readTimeout,
			IdleTimeout:      idleTimeout,
		})

		if err != nil {
			log.Fatal().Msg("Error shutting down server")
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.echo.yaml)")
	rootCmd.Flags().String("liveness-path", "/tmp/echoserver-live", "Path to put the liveness file.")
	rootCmd.Flags().Bool("banner", true, "Start without the banner.")
	rootCmd.Flags().String("host", "localhost", "The host to run the server on.")
	rootCmd.Flags().Int("port", 8080, "The port to run the server on.")
	rootCmd.Flags().Int("metrics-port", 9000, "The metrics port to run the server on. This can be the same as the host port, but is set to a separate port by default for security.")
	rootCmd.Flags().Int("wait", 15, "Wait timeout in seconds.")
	rootCmd.Flags().Int("write-timeout", 15, "Write timeout in seconds.")
	rootCmd.Flags().Int("read-timeout", 15, "Read timeout in seconds.")
	rootCmd.Flags().Int("idle-timeout", 60, "Idle timeout in seconds.")

	viper.BindPFlag("banner", rootCmd.Flags().Lookup("banner"))
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("metrics-port", rootCmd.Flags().Lookup("metrics-port"))
	viper.BindPFlag("wait", rootCmd.Flags().Lookup("wait"))
	viper.BindPFlag("write-timeout", rootCmd.Flags().Lookup("write-timeout"))
	viper.BindPFlag("read-timeout", rootCmd.Flags().Lookup("read-timeout"))
	viper.BindPFlag("idle-timeout", rootCmd.Flags().Lookup("idle-timeout"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".echoserver")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
