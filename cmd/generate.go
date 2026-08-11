package cmd

import (
	"context"
	"fmt"
	"log"

	"gh-task-gen/pkg/csv"
	"gh-task-gen/pkg/github"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	csvPath string
	owner   string
	repo    string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate GitHub issues from a CSV file",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("token")
		if token == "" {
			return fmt.Errorf("token is required. Pass via --token or GH_TASK_TOKEN env var")
		}

		owner = viper.GetString("owner")
		repo = viper.GetString("repo")
		if owner == "" || repo == "" {
			return fmt.Errorf("owner and repo are required. Pass via --owner, --repo or GH_TASK_OWNER, GH_TASK_REPO")
		}

		fmt.Printf("Processing CSV file: %s...\n", csvPath)

		parseCSV := csv.NewParser()
		tasks, err := parseCSV(csvPath)
		if err != nil {
			return fmt.Errorf("parse csv: %w", err)
		}

		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("initalize github client: %w", err)
		}

		ctx := context.Background()

		for _, task := range tasks {
			if err := task.Validate(); err != nil {
				log.Printf("Warning: skipping invalid task: %v", err)
				continue
			}

			result, err := client.CreateIssue(ctx, owner, repo, task)
			if err != nil {
				log.Printf("Warning: failed to create issue '%s': %v", task.Title, err)
				continue
			}

			log.Printf("Created #%d: %s\n", result.Number, result.URL)
		}

		log.Printf("Done!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&csvPath, "file", "f", "", "Path to the target CSV file (required)")
	generateCmd.Flags().StringVarP(&owner, "owner", "o", "", "GitHub repository owner")
	generateCmd.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository name")

	_ = viper.BindPFlag("owner", generateCmd.Flags().Lookup("owner"))
	_ = viper.BindPFlag("repo", generateCmd.Flags().Lookup("repo"))

	_ = generateCmd.MarkFlagRequired("file")
}
