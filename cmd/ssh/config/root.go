package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is largely clerical, will link sub commands here
var (
	debug = false

	Config = &cobra.Command{
		Use: "config",
		Short: "A profiling system for the SSH config file",
		Long:  "A profiling system for the SSH config file",
		Run: func(cmd *cobra.Command, args []string) {
			// If nothing is specified print help
			// This is because I want a uniform CLI,
			// So the user must use a command
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
)

// Execute is called on by the main process
func Execute() {
	cobra.OnInitialize(initConfig)

	// Debug - Add global flag
	Config.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Show debugging information")

	// Print errors
	if err := Config.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func initConfig() {

	// Write or load config file: ~/.historian.yaml
	if c, err := os.UserConfigDir(); err != nil {
		log.Fatalln(err)
	} else {
		dir := filepath.Join(c, "ccp")
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}

		viper.AddConfigPath(dir)
		viper.SetConfigName("profiles")
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()

		// Debug - Print configuration file
		if err := viper.ReadInConfig(); err == nil && debug {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// Write or read config
	if err := viper.SafeWriteConfig(); err != nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("cant read config: %s\n", err)
		}
	}

	// Update during execution
	viper.WatchConfig()

	// Debug - notify on configuration change
	if debug {
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name, e.Op)
		})
	}
}
