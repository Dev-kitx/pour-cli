package cmd

import (
	"os"

	"github.com/Dev-kitx/pour-cli/internal/config"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pour",
	Short: "AI Skill Manager",
	Long:  "pour — install and manage skills for any AI agent",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintWelcomeHint()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show pour version",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner(Version)
	},
}

func Execute() {
	config.EnsurePourDir()
	if err := rootCmd.Execute(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(searchCmd)
}
