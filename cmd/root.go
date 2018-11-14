package cmd

import (
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var endPoint string
var logLevel string
var apiToken string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hassio-cli",
	Short: "A brief description of your application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// set loglevel if posible
		logrusLevel, err := log.ParseLevel(viper.GetString("log-level"))

		if err == nil {
			log.SetLevel(logrusLevel)
		}
		log.WithFields(log.Fields{
			"cfgFile":  viper.GetString("config"),
			"endpoint": viper.GetString("endpoint"),
			"logLevel": viper.GetString("log-level"),
			"apiToken": viper.GetString("api-token"),
		}).Debugln("Debug flags")

	},
}

// Execute represents the entrypoint for when called without any subcommand
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error while executing rootCmd: %s", err)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.homeassistant.yaml)")
	rootCmd.PersistentFlags().StringVar(&endPoint, "endpoint", "", "Endpoint for hassio supervisor ( default is 'hassio' )")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Log level defaults to Warn")
	rootCmd.PersistentFlags().StringVar(&apiToken, "api-token", "", "Hassio api token")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("api-token", rootCmd.PersistentFlags().Lookup("api-token"))

	viper.SetDefault("endpoint", "hassio")
	viper.SetDefault("log-level", "warn")
	viper.SetDefault("api-token", "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("HASSIO")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// set loglevel if posible
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))

	if err == nil {
		log.SetLevel(logLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Error while finding home directory: %s", err)
		}

		// Search config in home directory with name ".homeassistant" (without extension).
		viper.AddConfigPath(home)
		log.WithField("homedir", home).Debug("Adding homedir to searchpath")
		viper.SetConfigName(".homeassistant")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("configfile", viper.ConfigFileUsed()).Info("Using configfile")
	} else {
		log.Info("No configfile found")
	}
}
