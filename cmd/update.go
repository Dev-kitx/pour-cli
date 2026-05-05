package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/config"
	gh "github.com/Dev-kitx/pour-cli/internal/github"
	"github.com/Dev-kitx/pour-cli/internal/installer"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all installed skills to latest version",
	Example: `  pour update
  pour update -a claude`,
	Run: runUpdate,
}

var updateAgentFilter string

func init() {
	updateCmd.Flags().StringVarP(&updateAgentFilter, "agent", "a", "", "only update skills for a specific agent")
}

func runUpdate(cmd *cobra.Command, args []string) {
	ui.SectionHeader("Updating Skills")

	db, err := config.LoadSkillsDB()
	if err != nil {
		ui.Error("Failed to load skills database: " + err.Error())
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config: " + err.Error())
		return
	}

	skills := db.Skills
	if updateAgentFilter != "" {
		var filtered []config.InstalledSkill
		for _, s := range skills {
			if s.Agent == updateAgentFilter {
				filtered = append(filtered, s)
			}
		}
		skills = filtered
	}

	if len(skills) == 0 {
		ui.Info("No skills installed yet")
		return
	}

	updated, skipped, failed := 0, 0, 0

	for _, s := range skills {
		spinner := ui.NewSpinner(fmt.Sprintf("Updating %s ...", ui.SkillBadge(s.Name)))

		newContent, err := gh.FetchSkillContent(s.Repo, s.Name, cfg.GitHubToken)
		if err != nil {
			ui.SpinnerFail(spinner, fmt.Sprintf("Failed to fetch %s", s.Name))
			ui.Muted("  " + err.Error())
			failed++
			continue
		}

		oldContent, err := os.ReadFile(s.Path)
		if err == nil && string(oldContent) == newContent {
			ui.SpinnerSuccess(spinner, fmt.Sprintf("%s is already up to date", ui.SkillBadge(s.Name)))
			skipped++
			continue
		}

		if err == nil {
			printDiff(string(oldContent), newContent)
		}

		_, err = installer.Install(s.Name, newContent, s.Agent, false)
		if err != nil {
			ui.SpinnerFail(spinner, "Failed to update "+s.Name)
			failed++
			continue
		}

		ui.SpinnerSuccess(spinner, fmt.Sprintf("Updated %s", ui.SkillBadge(s.Name)))
		updated++
	}

	fmt.Println()
	summary := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Padding(0, 2).
		Render(
			ui.SuccessStyle.Render(fmt.Sprintf("✓ %d updated", updated)) + "  " +
				ui.MutedStyle.Render(fmt.Sprintf("— %d skipped", skipped)) + "  " +
				ui.ErrorStyle.Render(fmt.Sprintf("✗ %d failed", failed)),
		)
	fmt.Println(summary)
	fmt.Println()
}

func printDiff(old, new string) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	oldSet := map[string]bool{}
	newSet := map[string]bool{}
	for _, l := range oldLines {
		oldSet[l] = true
	}
	for _, l := range newLines {
		newSet[l] = true
	}

	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71"))
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C"))

	fmt.Println()
	for _, l := range oldLines {
		if !newSet[l] {
			fmt.Println(delStyle.Render("  - " + l))
		}
	}
	for _, l := range newLines {
		if !oldSet[l] {
			fmt.Println(addStyle.Render("  + " + l))
		}
	}
	fmt.Println()
}
