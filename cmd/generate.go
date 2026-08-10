package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var csvPath string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate GitHub issues from a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.GetString("token")
		if token == "" {
			log.Fatal("Error: token is required. Pass via --token or GH_TASK_TOKEN env var")
		}

		fmt.Printf("Processing CSV file: %s...\n", csvPath)
		// TODO: Call csv parser into structs
		// TODO: Call github client to send issues to GitHub API
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&csvPath, "file", "f", "", "Path to the target CSV file (required)")
	_ = generateCmd.MarkFlagRequired("file")
}
