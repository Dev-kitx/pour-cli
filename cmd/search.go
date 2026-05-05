package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Dev-kitx/pour-cli/internal/config"
	gh "github.com/Dev-kitx/pour-cli/internal/github"
	"github.com/Dev-kitx/pour-cli/internal/ui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search GitHub for repos containing skills",
	Example: `  pour search
  pour search code review
  pour search ui design`,
	Run: runSearch,
}

var searchAgentFlag string

func init() {
	searchCmd.Flags().StringVarP(&searchAgentFlag, "agent", "a", "", "target AI agent for install (skips agent picker)")
}

func runSearch(cmd *cobra.Command, args []string) {
	query := strings.Join(args, " ")

	if query == "" {
		ui.SectionHeader("Searching for skills on GitHub")
	} else {
		ui.SectionHeader("Searching for \"" + query + "\"")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		ui.Error("Failed to load config: " + err.Error())
		return
	}

	spinner := ui.NewSpinner("Searching skills.sh ...")
	results, err := gh.SearchSkillsSH(query)
	if err != nil || len(results) == 0 {
		ui.SpinnerFail(spinner, "skills.sh unavailable, trying registry ...")
		fmt.Println()
		spinner = ui.NewSpinner("Searching registry ...")
		results, err = gh.SearchRegistry(query)
		if err != nil || len(results) == 0 {
			ui.SpinnerFail(spinner, "Registry unavailable, falling back to GitHub search")
			fmt.Println()
			spinner = ui.NewSpinner("Searching GitHub ...")
			results, err = gh.SearchSkills(query, cfg.GitHubToken)
			if err != nil {
				ui.SpinnerFail(spinner, "Search failed")
				ui.Error(err.Error())
				return
			}
		}
	}
	ui.SpinnerSuccess(spinner, fmt.Sprintf("Found %d repos", len(results)))
	fmt.Println()

	if len(results) == 0 {
		ui.Muted("No results found. Try a different query.")
		return
	}

	repo := runRepoPicker(results)
	if repo == "" {
		return
	}

	agent := searchAgentFlag
	if agent == "" {
		fmt.Println()
		ui.SectionHeader("Choose an agent")
		agent = runAgentPicker()
		if agent == "" {
			ui.Muted("No agent selected. Cancelled.")
			return
		}
	}

	fmt.Println()
	ui.SectionHeader("Installing Skills")
	installFromRepo(repo, agent, "", false)
}

// ── repo picker ───────────────────────────────────────────────────────────────

type repoItem struct {
	result gh.SearchResult
}

func (i repoItem) Title() string {
	stars := ""
	if i.result.Stars > 0 {
		stars = ui.MutedStyle.Render(fmt.Sprintf("  ★ %d", i.result.Stars))
	}
	return i.result.Repo + stars
}

func (i repoItem) Description() string {
	if i.result.Description == "" {
		return "No description"
	}
	if len(i.result.Description) > 70 {
		return i.result.Description[:70] + "..."
	}
	return i.result.Description
}

func (i repoItem) FilterValue() string {
	return i.result.Repo + " " + i.result.Description
}

type repoPickerModel struct {
	list     list.Model
	selected string
	quitting bool
}

func (m repoPickerModel) Init() tea.Cmd { return nil }

func (m repoPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(repoItem); ok {
				m.selected = item.result.Repo
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

func (m repoPickerModel) View() string {
	if m.quitting {
		return ""
	}
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#9B59B6")).Bold(true).
		Render("Select a repo to install from")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).
		Render("  enter: install  •  type to filter  •  q: cancel")
	return title + "\n\n" + m.list.View() + "\n" + hint + "\n"
}

func runRepoPicker(results []gh.SearchResult) string {
	items := make([]list.Item, len(results))
	for i, r := range results {
		items[i] = repoItem{result: r}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FF6B9D")).
		BorderForeground(lipgloss.Color("#9B59B6"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#00D4FF")).
		BorderForeground(lipgloss.Color("#9B59B6"))

	l := list.New(items, delegate, 70, 14)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)

	result, err := tea.NewProgram(repoPickerModel{list: l}).Run()
	if err != nil {
		return ""
	}
	final := result.(repoPickerModel)
	if final.quitting {
		return ""
	}
	return final.selected
}
