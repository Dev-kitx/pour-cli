package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Dev-kitx/pour-cli/internal/config"
)

func Install(skillName, content, agent string, global bool) (string, error) {
	destDir, err := resolveDestDir(agent, global)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fileName := sanitizeName(skillName) + ".md"
	destPath := filepath.Join(destDir, fileName)

	if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write skill: %w", err)
	}

	return destPath, nil
}

func Uninstall(skillName, agent string, global bool) error {
	destDir, err := resolveDestDir(agent, global)
	if err != nil {
		return err
	}

	fileName := sanitizeName(skillName) + ".md"
	destPath := filepath.Join(destDir, fileName)

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found for agent '%s'", skillName, agent)
	}

	return os.Remove(destPath)
}

func RecordInstall(skillName, description, agent, repo, path string) error {
	if err := config.EnsurePourDir(); err != nil {
		return err
	}

	db, err := config.LoadSkillsDB()
	if err != nil {
		return err
	}

	db.AddSkill(config.InstalledSkill{
		Name:        skillName,
		Description: description,
		Agent:       agent,
		Repo:        repo,
		InstalledAt: time.Now().Format(time.RFC3339),
		Path:        path,
	})

	return config.SaveSkillsDB(db)
}

func RecordUninstall(skillName, agent string) error {
	db, err := config.LoadSkillsDB()
	if err != nil {
		return err
	}
	db.RemoveSkill(skillName, agent)
	return config.SaveSkillsDB(db)
}

func resolveDestDir(agent string, global bool) (string, error) {
	agentKey := strings.ToLower(agent)

	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		subPath, ok := config.AgentGlobalPaths[agentKey]
		if !ok {
			subPath = filepath.Join(".pour", "skills", agentKey)
		}
		return filepath.Join(home, subPath), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	subPath, ok := config.AgentPaths[agentKey]
	if !ok {
		subPath = filepath.Join(".pour", "skills", agentKey)
	}
	return filepath.Join(cwd, subPath), nil
}

func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}
