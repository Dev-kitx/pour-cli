package cmd

import (
	"fmt"

	"github.com/Dev-kitx/pour-cli/internal/installer"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var removeAgentFlag string
var removeGlobalFlag bool

var removeCmd = &cobra.Command{
	Use:     "remove [skill]",
	Aliases: []string{"rm", "uninstall"},
	Short:   "Remove an installed skill",
	Example: `  pour remove code-reviewer -a claude
  pour rm sql-expert -a cursor`,
	Args: cobra.ExactArgs(1),
	Run:  runRemove,
}

func init() {
	removeCmd.Flags().StringVarP(&removeAgentFlag, "agent", "a", "", "target AI agent")
	removeCmd.Flags().BoolVarP(&removeGlobalFlag, "global", "g", false, "remove from global install")
	removeCmd.MarkFlagRequired("agent")
}

func runRemove(cmd *cobra.Command, args []string) {
	skillName := args[0]

	ui.SectionHeader("Removing Skill")

	spinner := ui.NewSpinner(fmt.Sprintf("Removing %s from %s ...", ui.SkillBadge(skillName), ui.AgentBadge(removeAgentFlag)))

	if err := installer.Uninstall(skillName, removeAgentFlag, removeGlobalFlag); err != nil {
		ui.SpinnerFail(spinner, "Failed to remove skill")
		ui.Error(err.Error())
		return
	}

	installer.RecordUninstall(skillName, removeAgentFlag)
	ui.SpinnerSuccess(spinner, fmt.Sprintf("Removed %s from %s", ui.SkillBadge(skillName), ui.AgentBadge(removeAgentFlag)))
	fmt.Println()
}
