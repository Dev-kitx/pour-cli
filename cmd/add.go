package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/config"
	gh "github.com/Dev-kitx/pour-cli/internal/github"
	"github.com/Dev-kitx/pour-cli/internal/installer"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	agentFlag  string
	globalFlag bool
	skillFlag  string
)

var addCmd = &cobra.Command{
	Use:   "add [repo]",
	Short: "Install skills from a GitHub repo",
	Example: `  pour add vercel-labs/agent-skills -a claude
  pour add john/my-skills -a cursor -s code-reviewer
  pour add owner/repo -a claude -g`,
	Args: cobra.ExactArgs(1),
	Run:  runAdd,
}

func init() {
	addCmd.Flags().StringVarP(&agentFlag, "agent", "a", "", "target AI agent (claude, cursor, windsurf …)")
	addCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "install globally")
	addCmd.Flags().StringVarP(&skillFlag, "skill", "s", "", "install a specific skill by name")
}

func runAdd(cmd *cobra.Command, args []string) {
	repo := args[0]

	agent := agentFlag
	if agent == "" {
		ui.SectionHeader("Choose an agent")
		agent = runAgentPicker()
		if agent == "" {
			ui.Muted("No agent selected. Cancelled.")
			return
		}
		fmt.Println()
	}

	ui.SectionHeader("Installing Skills")
	installFromRepo(repo, agent, skillFlag, globalFlag)
}

// installFromRepo is shared between add and search commands.
func installFromRepo(repo, agent, skillFilter string, global bool) {
	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config: " + err.Error())
		return
	}

	spinner := ui.NewSpinner("Fetching skills from " + ui.HighlightStyle.Render(repo) + " ...")
	skills, err := gh.FetchSkills(repo, cfg.GitHubToken)
	if err != nil {
		ui.SpinnerFail(spinner, "Failed to fetch skills")
		ui.Error(err.Error())
		return
	}
	ui.SpinnerSuccess(spinner, fmt.Sprintf("Found %d skill(s)", len(skills)))

	var toInstall []gh.Skill

	if skillFilter != "" {
		for _, s := range skills {
			if strings.EqualFold(s.Name, skillFilter) {
				toInstall = []gh.Skill{s}
				break
			}
		}
		if len(toInstall) == 0 {
			ui.Error(fmt.Sprintf("Skill '%s' not found in %s", skillFilter, repo))
			return
		}
	} else if len(skills) == 1 {
		toInstall = skills
	} else {
		selected := runSkillPicker(skills)
		if selected == nil {
			ui.Muted("No skills selected. Cancelled.")
			return
		}
		toInstall = selected
	}

	fmt.Println()
	var installed []string

	for _, skill := range toInstall {
		spin := ui.NewSpinner("Installing " + ui.SkillBadge(skill.Name) + " ...")

		content := skill.Content
		if content == "" {
			content, err = gh.FetchSkillContent(repo, skill.Name, cfg.GitHubToken)
			if err != nil {
				ui.SpinnerFail(spin, "Failed to fetch "+skill.Name)
				continue
			}
		}

		path, err := installer.Install(skill.Name, content, agent, global)
		if err != nil {
			ui.SpinnerFail(spin, "Failed to install "+skill.Name+": "+err.Error())
			continue
		}

		installer.RecordInstall(skill.Name, skill.Description, agent, repo, path)
		ui.SpinnerSuccess(spin, "Installed "+ui.SkillBadge(skill.Name))
		installed = append(installed, skill.Name)
	}

	if len(installed) > 0 {
		ui.PrintInstallSummary(installed, agent)
	}
}

// ── agent picker ─────────────────────────────────────────────────────────────

type agentItem struct{ name string }

func (i agentItem) Title() string       { return i.name }
func (i agentItem) Description() string { return "" }
func (i agentItem) FilterValue() string { return i.name }

type agentPickerModel struct {
	list     list.Model
	selected string
	quitting bool
}

func (m agentPickerModel) Init() tea.Cmd { return nil }

func (m agentPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(agentItem); ok {
				m.selected = item.name
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m agentPickerModel) View() string {
	if m.quitting {
		return ""
	}
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#9B59B6")).Bold(true).
		Render("Select an agent")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).
		Render("  enter: select  •  type to filter  •  q: cancel")
	return title + "\n\n" + m.list.View() + "\n" + hint + "\n"
}

func runAgentPicker() string {
	agents := config.KnownAgents()
	sort.Strings(agents)

	items := make([]list.Item, len(agents))
	for i, a := range agents {
		items[i] = agentItem{name: a}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FF6B9D")).
		BorderForeground(lipgloss.Color("#9B59B6"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#e2e2f0"))
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 40, 14)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)

	result, err := tea.NewProgram(agentPickerModel{list: l}).Run()
	if err != nil {
		return ""
	}
	return result.(agentPickerModel).selected
}

// ── skill picker ──────────────────────────────────────────────────────────────

type skillItem struct {
	skill gh.Skill
}

func (i skillItem) Title() string       { return "⚡ " + i.skill.Name }
func (i skillItem) Description() string { return i.skill.Description }
func (i skillItem) FilterValue() string { return i.skill.Name }

type pickerModel struct {
	list     list.Model
	selected []gh.Skill
	quitting bool
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(skillItem); ok {
				m.selected = []gh.Skill{item.skill}
			}
			return m, tea.Quit
		case "a":
			for _, item := range m.list.Items() {
				if s, ok := item.(skillItem); ok {
					m.selected = append(m.selected, s.skill)
				}
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickerModel) View() string {
	if m.quitting {
		return ""
	}
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#9B59B6")).Bold(true).
		Render("Select a skill to install")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).
		Render("  enter: install selected  •  a: install all  •  q: cancel")
	return title + "\n\n" + m.list.View() + "\n" + hint + "\n"
}

func runSkillPicker(skills []gh.Skill) []gh.Skill {
	items := make([]list.Item, len(skills))
	for i, s := range skills {
		items[i] = skillItem{skill: s}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FF6B9D")).
		BorderForeground(lipgloss.Color("#9B59B6"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#00D4FF")).
		BorderForeground(lipgloss.Color("#9B59B6"))

	l := list.New(items, delegate, 60, 14)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)

	result, err := tea.NewProgram(pickerModel{list: l}).Run()
	if err != nil {
		return nil
	}
	final := result.(pickerModel)
	if final.quitting {
		return nil
	}
	return final.selected
}
