package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Dev-kitx/pour-cli/internal/config"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Add a GitHub token for higher API rate limits",
	Run:   runAuthLogin,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current auth status",
	Run:   runAuthStatus,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	Run:   runAuthLogout,
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) {
	ui.SectionHeader("GitHub Authentication")

	ui.Info("A GitHub token allows 5000 API requests/hour instead of 60")
	ui.Muted("Create one at: github.com/settings/tokens")
	ui.Muted("Only needs 'public_repo' scope for public repositories")
	fmt.Println()

	fmt.Print(ui.HighlightStyle.Render("  Enter your GitHub token: "))
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		ui.Error("Failed to read token")
		return
	}
	token = strings.TrimSpace(token)

	if token == "" {
		ui.Error("Token cannot be empty")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config: " + err.Error())
		return
	}

	cfg.GitHubToken = token
	if err := config.SaveConfig(cfg); err != nil {
		ui.Error("Failed to save config: " + err.Error())
		return
	}

	ui.Success("GitHub token saved successfully")
	ui.Muted("Stored at " + config.ConfigPath())
	fmt.Println()
}

func runAuthStatus(cmd *cobra.Command, args []string) {
	ui.SectionHeader("Auth Status")

	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config")
		return
	}

	if cfg.GitHubToken != "" {
		masked := cfg.GitHubToken[:4] + strings.Repeat("*", len(cfg.GitHubToken)-8) + cfg.GitHubToken[len(cfg.GitHubToken)-4:]
		ui.Success("GitHub token configured")
		ui.Muted("Token: " + masked)
		ui.Muted("Rate limit: 5000 requests/hour")
	} else {
		ui.Info("No GitHub token configured")
		ui.Muted("Rate limit: 60 requests/hour")
		ui.Muted("Run `pour auth login` to add a token")
	}
	fmt.Println()
}

func runAuthLogout(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config")
		return
	}

	cfg.GitHubToken = ""
	if err := config.SaveConfig(cfg); err != nil {
		ui.Error("Failed to save config")
		return
	}

	ui.Success("Credentials removed")
	fmt.Println()
}
