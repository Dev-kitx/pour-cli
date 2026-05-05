package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	purple    = lipgloss.Color("#9B59B6")
	pink      = lipgloss.Color("#FF6B9D")
	cyan      = lipgloss.Color("#00D4FF")
	green     = lipgloss.Color("#2ECC71")
	yellow    = lipgloss.Color("#F1C40F")
	white     = lipgloss.Color("#FFFFFF")
	darkGray  = lipgloss.Color("#2D2D2D")
	lightGray = lipgloss.Color("#888888")

	TitleStyle = lipgloss.NewStyle().
			Foreground(pink).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(cyan)

	MutedStyle = lipgloss.NewStyle().
			Foreground(lightGray)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 1)

	SkillStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(purple).
			Padding(0, 1).
			Bold(true)

	AgentStyle = lipgloss.NewStyle().
			Foreground(darkGray).
			Background(cyan).
			Padding(0, 1).
			Bold(true)
)

func PrintBanner(version string) {
	pterm.Println()

	art := `
██████╗  ██████╗ ██╗   ██╗██████╗
██╔══██╗██╔═══██╗██║   ██║██╔══██╗
██████╔╝██║   ██║██║   ██║██████╔╝
██╔═══╝ ██║   ██║██║   ██║██╔══██╗
██║     ╚██████╔╝╚██████╔╝██║  ██║
╚═╝      ╚═════╝  ╚═════╝ ╚═╝  ╚═╝`

	fmt.Println(TitleStyle.Render(art))
	fmt.Println()

	tagline := BoxStyle.Render(
		SubtitleStyle.Render("✦ ") +
			TitleStyle.Render("AI Skill Manager") +
			SubtitleStyle.Render(" ✦") +
			"  " +
			MutedStyle.Render(version),
	)
	fmt.Println(tagline)
	fmt.Println()
}

func Success(msg string) {
	fmt.Println(SuccessStyle.Render("  ✓ " + msg))
}

func Error(msg string) {
	fmt.Println(ErrorStyle.Render("  ✗ " + msg))
}

func Info(msg string) {
	fmt.Println(InfoStyle.Render("  → " + msg))
}

func Muted(msg string) {
	fmt.Println(MutedStyle.Render("    " + msg))
}

func Highlight(msg string) {
	fmt.Println(HighlightStyle.Render("  ★ " + msg))
}

func SkillBadge(name string) string {
	return SkillStyle.Render(name)
}

func AgentBadge(name string) string {
	return AgentStyle.Render(name)
}

func SectionHeader(title string) {
	fmt.Println()
	line := lipgloss.NewStyle().
		Foreground(purple).
		Bold(true).
		Render("━━━ " + title + " ━━━")
	fmt.Println(line)
	fmt.Println()
}

func PrintSkillCard(name, description, author string) {
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(purple).
		Padding(0, 2).
		Width(50).
		Render(
			HighlightStyle.Render("⚡ "+name) + "\n" +
				MutedStyle.Render(description) + "\n" +
				InfoStyle.Render("by "+author),
		)
	fmt.Println(card)
}
