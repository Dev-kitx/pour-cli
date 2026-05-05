package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/config"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var infoAgentFlag string

var infoCmd = &cobra.Command{
	Use:   "info <skill>",
	Short: "Show details and content of an installed skill",
	Example: `  pour info code-reviewer -a claude
  pour info ui-design -a cursor`,
	Args: cobra.ExactArgs(1),
	Run:  runInfo,
}

func init() {
	infoCmd.Flags().StringVarP(&infoAgentFlag, "agent", "a", "", "agent the skill is installed for")
	infoCmd.MarkFlagRequired("agent")
}

func runInfo(cmd *cobra.Command, args []string) {
	skillName := args[0]

	db, err := config.LoadSkillsDB()
	if err != nil {
		ui.Error("Failed to load skills database: " + err.Error())
		return
	}

	var found *config.InstalledSkill
	for _, s := range db.Skills {
		s := s
		if strings.EqualFold(s.Name, skillName) && strings.EqualFold(s.Agent, infoAgentFlag) {
			found = &s
			break
		}
	}

	if found == nil {
		ui.Error(fmt.Sprintf("Skill '%s' not found for agent '%s'", skillName, infoAgentFlag))
		ui.Muted("Run `pour list` to see installed skills")
		return
	}

	ui.SectionHeader("Skill Info")

	header := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9B59B6")).
		Padding(1, 2).
		Width(60).
		Render(
			ui.HighlightStyle.Render("⚡ "+found.Name) + "\n" +
				ui.MutedStyle.Render(found.Description) + "\n\n" +
				ui.InfoStyle.Render("Agent  : ") + ui.AgentBadge(found.Agent) + "\n" +
				ui.InfoStyle.Render("Source : ") + ui.HighlightStyle.Render(found.Repo) + "\n" +
				ui.InfoStyle.Render("Path   : ") + ui.MutedStyle.Render(found.Path) + "\n" +
				ui.InfoStyle.Render("Installed : ") + ui.MutedStyle.Render(found.InstalledAt[:10]),
		)
	fmt.Println(header)
	fmt.Println()

	content, err := os.ReadFile(found.Path)
	if err != nil {
		ui.Error("Could not read skill file: " + err.Error())
		return
	}

	ui.SectionHeader("Content")

	lines := strings.Split(string(content), "\n")
	inFrontmatter := false
	frontmatterDone := false

	for _, line := range lines {
		if line == "---" && !frontmatterDone {
			if !inFrontmatter {
				inFrontmatter = true
				fmt.Println(ui.MutedStyle.Render(line))
			} else {
				inFrontmatter = false
				frontmatterDone = true
				fmt.Println(ui.MutedStyle.Render(line))
			}
			continue
		}
		if inFrontmatter {
			fmt.Println(ui.MutedStyle.Render(line))
			continue
		}
		if strings.HasPrefix(line, "# ") {
			fmt.Println(ui.HighlightStyle.Render(line))
		} else if strings.HasPrefix(line, "## ") {
			fmt.Println(ui.InfoStyle.Render(line))
		} else if strings.HasPrefix(line, "### ") {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#00D4FF")).Italic(true).Render(line))
		} else {
			fmt.Println(ui.TitleStyle.Render("  ") + line)
		}
	}
	fmt.Println()
}

