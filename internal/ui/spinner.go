package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

func NewSpinner(msg string) *pterm.SpinnerPrinter {
	spinner, _ := pterm.DefaultSpinner.
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		WithDelay(80 * time.Millisecond).
		WithText(InfoStyle.Render(msg)).
		Start()
	return spinner
}

func SpinnerSuccess(spinner *pterm.SpinnerPrinter, msg string) {
	spinner.Success(SuccessStyle.Render(msg))
}

func SpinnerFail(spinner *pterm.SpinnerPrinter, msg string) {
	spinner.Fail(ErrorStyle.Render(msg))
}

func PrintDivider() {
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Render("  ────────────────────────────────────────")
	fmt.Println(divider)
}

func PrintInstallSummary(skills []string, agent string) {
	fmt.Println()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(green).
		Padding(1, 2).
		Render(
			SuccessStyle.Render("Installation Complete!") + "\n\n" +
				InfoStyle.Render("Agent : ") + AgentBadge(agent) + "\n" +
				InfoStyle.Render("Skills: ") + formatSkillList(skills),
		)
	fmt.Println(box)
	fmt.Println()
}

func formatSkillList(skills []string) string {
	result := ""
	for _, s := range skills {
		result += SkillBadge(s) + " "
	}
	return result
}

func PrintWelcomeHint() {
	fmt.Println(MutedStyle.Render("  Run ") +
		HighlightStyle.Render("pour add <repo> -a <agent>") +
		MutedStyle.Render(" to install your first skill"))
	fmt.Println(MutedStyle.Render("  Run ") +
		HighlightStyle.Render("pour list") +
		MutedStyle.Render(" to see installed skills"))
	fmt.Println()
}
