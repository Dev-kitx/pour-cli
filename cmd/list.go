package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/config"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var listAgentFilter string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed skills",
	Example: `  pour list
  pour list -a claude`,
	Run: runList,
}

func init() {
	listCmd.Flags().StringVarP(&listAgentFilter, "agent", "a", "", "filter by agent")
}

func runList(cmd *cobra.Command, args []string) {
	db, err := config.LoadSkillsDB()
	if err != nil {
		ui.Error("Failed to load skills database: " + err.Error())
		return
	}

	skills := db.Skills
	if listAgentFilter != "" {
		var filtered []config.InstalledSkill
		for _, s := range skills {
			if s.Agent == listAgentFilter {
				filtered = append(filtered, s)
			}
		}
		skills = filtered
	}

	if len(skills) == 0 {
		empty := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(1, 4).
			Render(
				ui.MutedStyle.Render("No skills installed yet.\n\n") +
					ui.HighlightStyle.Render("  pour add <repo> -a <agent>") +
					ui.MutedStyle.Render("  to get started"),
			)
		fmt.Println(empty)
		return
	}

	ui.SectionHeader(fmt.Sprintf("Installed Skills (%d)", len(skills)))

	grouped := make(map[string][]config.InstalledSkill)
	for _, s := range skills {
		grouped[s.Agent] = append(grouped[s.Agent], s)
	}

	for agent, agentSkills := range grouped {
		fmt.Println("  " + ui.AgentBadge(agent))
		fmt.Println()

		for _, s := range agentSkills {
			desc := s.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			row := lipgloss.NewStyle().
				PaddingLeft(4).
				Render(
					ui.SkillBadge(s.Name) +
						"  " +
						ui.MutedStyle.Render(desc) +
						"\n" +
						lipgloss.NewStyle().PaddingLeft(4).Render(
							ui.InfoStyle.Render("from ") +
								ui.HighlightStyle.Render(s.Repo) +
								ui.MutedStyle.Render("  "+s.InstalledAt[:10]),
						),
				)
			fmt.Println(row)
			fmt.Println()
		}
	}
}
