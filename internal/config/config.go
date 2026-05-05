package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GitHubToken string            `json:"github_token,omitempty"`
	Agents      map[string]string `json:"agents,omitempty"`
}

type InstalledSkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Agent       string `json:"agent"`
	Repo        string `json:"repo"`
	InstalledAt string `json:"installed_at"`
	Path        string `json:"path"`
}

type SkillsDB struct {
	Skills []InstalledSkill `json:"skills"`
}

// AgentPaths maps agent flag values to their project-level skill directories.
var AgentPaths = map[string]string{
	// AiderDesk
	"aider-desk": ".aider-desk/skills",
	// Amp / Kimi Code CLI / Replit / Universal
	"amp":       ".agents/skills",
	"kimi-cli":  ".agents/skills",
	"replit":    ".agents/skills",
	"universal": ".agents/skills",
	// Antigravity
	"antigravity": ".agents/skills",
	// Augment
	"augment": ".augment/skills",
	// IBM Bob
	"bob": ".bob/skills",
	// Claude Code (canonical name + legacy alias)
	"claude-code": ".claude/skills",
	"claude":      ".claude/skills",
	// OpenClaw
	"openclaw": "skills",
	// Cline / Dexto / Warp
	"cline": ".agents/skills",
	"dexto": ".agents/skills",
	"warp":  ".agents/skills",
	// CodeArts Agent
	"codearts-agent": ".codeartsdoer/skills",
	// CodeBuddy
	"codebuddy": ".codebuddy/skills",
	// Codemaker
	"codemaker": ".codemaker/skills",
	// Code Studio
	"codestudio": ".codestudio/skills",
	// Codex
	"codex": ".agents/skills",
	// Command Code
	"command-code": ".commandcode/skills",
	// Continue
	"continue": ".continue/skills",
	// Cortex Code
	"cortex": ".cortex/skills",
	// Crush
	"crush": ".crush/skills",
	// Cursor
	"cursor": ".agents/skills",
	// Deep Agents
	"deepagents": ".agents/skills",
	// Devin
	"devin": ".devin/skills",
	// Droid
	"droid": ".factory/skills",
	// Firebender
	"firebender": ".agents/skills",
	// ForgeCode
	"forgecode": ".forge/skills",
	// Gemini CLI
	"gemini-cli": ".agents/skills",
	// GitHub Copilot
	"github-copilot": ".agents/skills",
	"copilot":        ".agents/skills",
	// Goose
	"goose": ".goose/skills",
	// Junie
	"junie": ".junie/skills",
	// iFlow CLI
	"iflow-cli": ".iflow/skills",
	// Kilo Code
	"kilo": ".kilocode/skills",
	// Kiro CLI
	"kiro-cli": ".kiro/skills",
	// Kode
	"kode": ".kode/skills",
	// MCPJam
	"mcpjam": ".mcpjam/skills",
	// Mistral Vibe
	"mistral-vibe": ".vibe/skills",
	// Mux
	"mux": ".mux/skills",
	// OpenCode
	"opencode": ".agents/skills",
	// OpenHands
	"openhands": ".openhands/skills",
	// Pi
	"pi": ".pi/skills",
	// Qoder
	"qoder": ".qoder/skills",
	// Qwen Code
	"qwen-code": ".qwen/skills",
	// Rovo Dev
	"rovodev": ".rovodev/skills",
	// Roo Code
	"roo": ".roo/skills",
	// Tabnine CLI
	"tabnine-cli": ".tabnine/agent/skills",
	// Trae / Trae CN
	"trae":    ".trae/skills",
	"trae-cn": ".trae/skills",
	// Windsurf
	"windsurf": ".windsurf/skills",
	// Zencoder
	"zencoder": ".zencoder/skills",
	// Neovate
	"neovate": ".neovate/skills",
	// Pochi
	"pochi": ".pochi/skills",
	// AdaL
	"adal": ".adal/skills",
	// Vercel skills subdirectories
	"skills-curated":      "skills/.curated",
	"skills-experimental": "skills/.experimental",
	"skills-system":       "skills/.system",
}

// AgentGlobalPaths maps agent flag values to their global (home) skill directories.
var AgentGlobalPaths = map[string]string{
	"aider-desk":     ".aider-desk/skills",
	"amp":            ".config/agents/skills",
	"kimi-cli":       ".config/agents/skills",
	"replit":         ".config/agents/skills",
	"universal":      ".config/agents/skills",
	"antigravity":    ".gemini/antigravity/skills",
	"augment":        ".augment/skills",
	"bob":            ".bob/skills",
	"claude-code":    ".claude/skills",
	"claude":         ".claude/skills",
	"openclaw":       ".openclaw/skills",
	"cline":          ".agents/skills",
	"dexto":          ".agents/skills",
	"warp":           ".agents/skills",
	"codearts-agent": ".codeartsdoer/skills",
	"codebuddy":      ".codebuddy/skills",
	"codemaker":      ".codemaker/skills",
	"codestudio":     ".codestudio/skills",
	"codex":          ".codex/skills",
	"command-code":   ".commandcode/skills",
	"continue":       ".continue/skills",
	"cortex":         ".snowflake/cortex/skills",
	"crush":          ".config/crush/skills",
	"cursor":         ".cursor/skills",
	"deepagents":     ".deepagents/agent/skills",
	"devin":          ".config/devin/skills",
	"droid":          ".factory/skills",
	"firebender":     ".firebender/skills",
	"forgecode":      ".forge/skills",
	"gemini-cli":     ".gemini/skills",
	"github-copilot": ".copilot/skills",
	"copilot":        ".copilot/skills",
	"goose":          ".config/goose/skills",
	"junie":          ".junie/skills",
	"iflow-cli":      ".iflow/skills",
	"kilo":           ".kilocode/skills",
	"kiro-cli":       ".kiro/skills",
	"kode":           ".kode/skills",
	"mcpjam":         ".mcpjam/skills",
	"mistral-vibe":   ".vibe/skills",
	"mux":            ".mux/skills",
	"opencode":       ".config/opencode/skills",
	"openhands":      ".openhands/skills",
	"pi":             ".pi/agent/skills",
	"qoder":          ".qoder/skills",
	"qwen-code":      ".qwen/skills",
	"rovodev":        ".rovodev/skills",
	"roo":            ".roo/skills",
	"tabnine-cli":    ".tabnine/agent/skills",
	"trae":           ".trae/skills",
	"trae-cn":        ".trae-cn/skills",
	"windsurf":       ".codeium/windsurf/skills",
	"zencoder":       ".zencoder/skills",
	"neovate":             ".neovate/skills",
	"pochi":               ".pochi/skills",
	"adal":                ".adal/skills",
	"skills-curated":      "skills/.curated",
	"skills-experimental": "skills/.experimental",
	"skills-system":       "skills/.system",
}

// KnownAgents returns the list of all supported agent names.
func KnownAgents() []string {
	seen := map[string]bool{}
	var agents []string
	for k := range AgentPaths {
		if !seen[k] {
			seen[k] = true
			agents = append(agents, k)
		}
	}
	return agents
}

func PourDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pour")
}

func ConfigPath() string {
	return filepath.Join(PourDir(), "config.json")
}

func SkillsDBPath() string {
	return filepath.Join(PourDir(), "skills.json")
}

func EnsurePourDir() error {
	dirs := []string{
		PourDir(),
		filepath.Join(PourDir(), "cache"),
		filepath.Join(PourDir(), "skills"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig() (*Config, error) {
	cfg := &Config{Agents: make(map[string]string)}
	data, err := os.ReadFile(ConfigPath())
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	return cfg, json.Unmarshal(data, cfg)
}

func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
}

func LoadSkillsDB() (*SkillsDB, error) {
	db := &SkillsDB{}
	data, err := os.ReadFile(SkillsDBPath())
	if os.IsNotExist(err) {
		return db, nil
	}
	if err != nil {
		return nil, err
	}
	return db, json.Unmarshal(data, db)
}

func SaveSkillsDB(db *SkillsDB) error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SkillsDBPath(), data, 0644)
}

func (db *SkillsDB) AddSkill(skill InstalledSkill) {
	for i, s := range db.Skills {
		if s.Name == skill.Name && s.Agent == skill.Agent {
			db.Skills[i] = skill
			return
		}
	}
	db.Skills = append(db.Skills, skill)
}

func (db *SkillsDB) RemoveSkill(name, agent string) bool {
	for i, s := range db.Skills {
		if s.Name == name && s.Agent == agent {
			db.Skills = append(db.Skills[:i], db.Skills[i+1:]...)
			return true
		}
	}
	return false
}
