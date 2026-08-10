// Package cmd
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "gh-task-gen",
	Short: "Generate GitHub Issues from a CSV file",
	Long:  "A CLI tool that parses a structured CSV file and bulk-creates GitHub issues across repositories and projects",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Initalization
func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is .gh-task-gen.yaml)")
	rootCmd.PersistentFlags().String("token", "", "GitHub Personal Access Token")

	// Bind flags to Viper
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
}

// Config Initalization
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		path, err := os.UserConfigDir()

		if err == nil {
			viper.AddConfigPath(path)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".gh-task-gen")
		}
	}

	// e.g. GH_TASK_TOKEN
	viper.SetEnvPrefix("GH_TASK")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
