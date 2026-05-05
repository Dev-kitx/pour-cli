package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Dev-kitx/pour-cli/internal/config"
)

type Skill struct {
	Name        string
	Description string
	Content     string
	Path        string
}

type ghContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

type SkillMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type SearchResult struct {
	Repo        string
	Description string
	Stars       int
	Skills      []string
	URL         string
}

func FetchSkills(repo, token string) ([]Skill, error) {
	owner, repoName, ref, err := parseRepo(repo)
	if err != nil {
		return nil, err
	}
	return scanRepoForSkills(owner, repoName, ref, token)
}

// knownSkillDirs returns the unique skill base directories across all known agents.
// e.g. ".claude/skills", ".windsurf/skills", ".agents/skills", etc.
func knownSkillDirs() []string {
	seen := map[string]bool{}
	var dirs []string
	for _, path := range config.AgentPaths {
		if !seen[path] {
			seen[path] = true
			dirs = append(dirs, path)
		}
	}
	return dirs
}

// scanRepoForSkills checks every known agent skill directory for SKILL.md files,
// then falls back to root-level and flat single-skill layouts.
func scanRepoForSkills(owner, repo, ref, token string) ([]Skill, error) {
	refParam := ""
	if ref != "" {
		refParam = "?ref=" + ref
	}

	base := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", owner, repo)
	var skills []Skill
	seen := map[string]bool{}

	addSkill := func(name, content string) {
		if !seen[name] {
			seen[name] = true
			skills = append(skills, parseSkill(name, content))
		}
	}

	// 1. Check root SKILL.md (single-skill repos)
	rootEntries, err := fetchContents(base+refParam, token)
	if err != nil {
		return nil, err
	}
	for _, e := range rootEntries {
		if e.Type == "file" && e.Name == "SKILL.md" {
			if content, err := decodeContent(e); err == nil {
				addSkill(repo, content)
			}
		}
	}

	// 2. Check <name>/SKILL.md at root (flat layout, e.g. fix/SKILL.md)
	// Skip known source dirs to avoid unnecessary API calls on large repos.
	skipDirs := map[string]bool{
		"src": true, "lib": true, "pkg": true, "cmd": true, "internal": true,
		"dist": true, "build": true, "out": true, "bin": true, "vendor": true,
		"node_modules": true, "packages": true, "test": true, "tests": true,
		"docs": true, "examples": true, "scripts": true, "tools": true,
	}
	for _, e := range rootEntries {
		if e.Type == "dir" && !skipDirs[e.Name] {
			url := fmt.Sprintf("%s/%s/SKILL.md%s", base, e.Path, refParam)
			if entries, err := fetchContents(url, token); err == nil && len(entries) > 0 {
				if content, err := decodeContent(entries[0]); err == nil {
					addSkill(e.Name, content)
				}
			}
		}
	}

	// 3. Check all known agent skill dirs: <agentDir>/SKILL.md and <agentDir>/<name>/SKILL.md
	for _, skillDir := range knownSkillDirs() {
		url := fmt.Sprintf("%s/%s%s", base, skillDir, refParam)
		entries, err := fetchContents(url, token)
		if err != nil {
			continue // this agent dir doesn't exist in this repo
		}
		// direct SKILL.md inside the agent skill dir
		for _, e := range entries {
			if e.Type == "file" && e.Name == "SKILL.md" {
				if content, err := decodeContent(e); err == nil {
					addSkill(skillDir, content)
				}
			}
		}
		for _, e := range entries {
			if e.Type != "dir" {
				continue
			}
			skillURL := fmt.Sprintf("%s/%s/SKILL.md%s", base, e.Path, refParam)
			skillEntries, err := fetchContents(skillURL, token)
			if err != nil || len(skillEntries) == 0 {
				continue
			}
			if content, err := decodeContent(skillEntries[0]); err == nil {
				addSkill(e.Name, content)
			}
		}
	}

	if len(skills) == 0 {
		return nil, fmt.Errorf("no skills found in %s/%s — make sure repo contains a SKILL.md", owner, repo)
	}

	return skills, nil
}

func FetchSkillContent(repo, skillName, token string) (string, error) {
	owner, repoName, ref, err := parseRepo(repo)
	if err != nil {
		return "", err
	}

	skills, err := scanRepoForSkills(owner, repoName, ref, token)
	if err != nil {
		return "", err
	}

	for _, s := range skills {
		if strings.EqualFold(s.Name, skillName) {
			return s.Content, nil
		}
	}

	return "", fmt.Errorf("skill '%s' not found in %s", skillName, repo)
}

const registryURL = "https://raw.githubusercontent.com/Dev-kitx/pour-skills-registry/main/index.json"

type registry struct {
	UpdatedAt string         `json:"updated_at"`
	Count     int            `json:"count"`
	Repos     []SearchResult `json:"repos"`
}

func SearchRegistry(query string) ([]SearchResult, error) {
	resp, err := http.Get(registryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("registry unavailable: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reg registry
	if err := json.Unmarshal(body, &reg); err != nil {
		return nil, err
	}

	if query == "" {
		return reg.Repos, nil
	}

	q := strings.ToLower(query)
	var results []SearchResult
	for _, r := range reg.Repos {
		if strings.Contains(strings.ToLower(r.Repo), q) ||
			strings.Contains(strings.ToLower(r.Description), q) {
			results = append(results, r)
		}
	}
	return results, nil
}

type ghSearchResponse struct {
	Items []struct {
		Repository struct {
			FullName    string `json:"full_name"`
			Description string `json:"description"`
			StarCount   int    `json:"stargazers_count"`
			HTMLURL     string `json:"html_url"`
		} `json:"repository"`
	} `json:"items"`
}

func SearchSkills(query, token string) ([]SearchResult, error) {
	q := "filename:SKILL.md"
	if query != "" {
		q += "+" + strings.ReplaceAll(query, " ", "+")
	}
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=10", q)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("rate limit exceeded — run `pour auth login` to add a GitHub token")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ghSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var results []SearchResult
	for _, item := range result.Items {
		repo := item.Repository.FullName
		if seen[repo] {
			continue
		}
		seen[repo] = true
		results = append(results, SearchResult{
			Repo:        repo,
			Description: item.Repository.Description,
			Stars:       item.Repository.StarCount,
			URL:         item.Repository.HTMLURL,
		})
	}
	return results, nil
}

type skillsSHResponse struct {
	Skills []struct {
		Name     string `json:"name"`
		Installs int    `json:"installs"`
		Source   string `json:"source"` // "owner/repo"
	} `json:"skills"`
}

// SearchSkillsSH queries the skills.sh public search API.
// Results are deduplicated by repo and sorted by total installs.
func SearchSkillsSH(query string) ([]SearchResult, error) {
	if len(strings.TrimSpace(query)) < 2 {
		return nil, fmt.Errorf("query must be at least 2 characters")
	}

	apiURL := fmt.Sprintf("https://skills.sh/api/search?q=%s&limit=50", url.QueryEscape(query))
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "pour-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("skills.sh unavailable: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result skillsSHResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// deduplicate by repo, collecting skill names and summing installs
	repoInstalls := map[string]int{}
	repoSkillNames := map[string][]string{}
	seen := map[string]bool{}
	var ordered []string
	for _, s := range result.Skills {
		if s.Source == "" {
			continue
		}
		repoInstalls[s.Source] += s.Installs
		repoSkillNames[s.Source] = append(repoSkillNames[s.Source], s.Name)
		if !seen[s.Source] {
			seen[s.Source] = true
			ordered = append(ordered, s.Source)
		}
	}

	var results []SearchResult
	for _, repo := range ordered {
		names := repoSkillNames[repo]
		desc := strings.Join(names, ", ")
		if len(desc) > 80 {
			desc = desc[:80] + "..."
		}
		results = append(results, SearchResult{
			Repo:        repo,
			Description: desc,
			Stars:       repoInstalls[repo],
			URL:         "https://github.com/" + repo,
		})
	}
	return results, nil
}

func fetchContents(url, token string) ([]ghContent, error) {
	body, err := doRequest(url, token)
	if err != nil {
		return nil, err
	}

	var single ghContent
	if err := json.Unmarshal(body, &single); err == nil && single.Name != "" {
		return []ghContent{single}, nil
	}

	var items []ghContent
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchFileContent(url, token string) (string, error) {
	body, err := doRequest(url, token)
	if err != nil {
		return "", err
	}

	var item ghContent
	if err := json.Unmarshal(body, &item); err != nil {
		return "", err
	}

	return decodeContent(item)
}

func decodeContent(item ghContent) (string, error) {
	if item.Encoding == "base64" {
		clean := strings.ReplaceAll(item.Content, "\n", "")
		decoded, err := base64.StdEncoding.DecodeString(clean)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
	if item.DownloadURL != "" {
		resp, err := http.Get(item.DownloadURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		return string(data), err
	}
	return "", fmt.Errorf("cannot decode content")
}

func doRequest(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("rate limit exceeded — run `pour auth login` to add a GitHub token")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// parseRepo handles owner/repo and owner/repo@ref formats.
func parseRepo(repo string) (owner, name, ref string, err error) {
	repo = strings.TrimPrefix(repo, "github.com/")
	repo = strings.TrimPrefix(repo, "https://github.com/")

	if idx := strings.Index(repo, "@"); idx != -1 {
		ref = repo[idx+1:]
		repo = repo[:idx]
	}

	parts := strings.Split(repo, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid repo format, use owner/repo or owner/repo@ref")
	}
	return parts[0], parts[1], ref, nil
}

func parseSkill(name, content string) Skill {
	description := extractDescription(content)
	return Skill{
		Name:        name,
		Description: description,
		Content:     content,
	}
}

func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	inFrontmatter := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			break
		}
		if inFrontmatter && strings.HasPrefix(line, "description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		}
	}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") {
			if len(line) > 80 {
				return line[:80] + "..."
			}
			return line
		}
	}
	return "No description"
}

func extractSkillName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}
