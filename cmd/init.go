package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [skill-name]",
	Short: "Create a new SKILL.md template",
	Example: `  pour init my-skill
  pour init code-reviewer`,
	Args: cobra.ExactArgs(1),
	Run:  runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	skillName := args[0]

	ui.SectionHeader("Creating New Skill")

	dir := filepath.Join("skills", skillName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		ui.Error("Failed to create directory: " + err.Error())
		return
	}

	skillPath := filepath.Join(dir, "SKILL.md")
	if _, err := os.Stat(skillPath); err == nil {
		ui.Error("SKILL.md already exists at " + skillPath)
		return
	}

	template := fmt.Sprintf(`---
name: %s
description: A short description of what this skill does
tags:
  - general
agents:
  - claude
  - cursor
  - windsurf
---

# %s

## Instructions

Describe what the AI agent should do when this skill is active.

## Examples

Provide examples of how to use this skill.

## Notes

Any additional context or constraints.
`, skillName, skillName)

	if err := os.WriteFile(skillPath, []byte(template), 0644); err != nil {
		ui.Error("Failed to create SKILL.md: " + err.Error())
		return
	}

	ui.Success("Created " + ui.SkillBadge(skillName))
	fmt.Println()
	ui.Info("File: " + ui.HighlightStyle.Render(skillPath))
	ui.Muted("Edit the file to add your skill's instructions")
	fmt.Println()

	preview := fmt.Sprintf(`skills/
└── %s/
    └── SKILL.md  ← edit this`, skillName)

	box := lipglossBox(preview)
	fmt.Println(box)
	fmt.Println()

	ui.Muted("When ready, share your repo and users can install with:")
	fmt.Println("  " + ui.HighlightStyle.Render("pour add your-username/your-repo -a claude"))
	fmt.Println()
}

func lipglossBox(content string) string {
	_ = lipgloss.Color("#000") // ensure import used
	return ui.BoxStyle.Render(ui.MutedStyle.Render(content))
}
