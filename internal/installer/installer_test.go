package installer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ui-design", "ui-design"},
		{"UI Design", "ui-design"},
		{"My Skill Name", "my-skill-name"},
		{"already-lower", "already-lower"},
		{"UPPER", "upper"},
	}
	for _, tc := range tests {
		got := sanitizeName(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestInstall_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	content := "# My Skill\nDo the thing."
	path, err := Install("my-skill", content, "claude", false)
	if err != nil {
		t.Fatalf("Install error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if string(data) != content {
		t.Errorf("file content = %q, want %q", string(data), content)
	}

	// path should be inside the temp dir
	rel, err := filepath.Rel(tmpDir, path)
	if err != nil || rel == "" {
		t.Errorf("installed path %q is not under tmpDir %q", path, tmpDir)
	}
}

func TestInstall_SanitizesFileName(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	path, err := Install("My Skill", "content", "claude", false)
	if err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if filepath.Base(path) != "my-skill.md" {
		t.Errorf("expected filename my-skill.md, got %q", filepath.Base(path))
	}
}

func TestUninstall_RemovesFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	if _, err := Install("my-skill", "content", "claude", false); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if err := Uninstall("my-skill", "claude", false); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}
}

func TestUninstall_MissingSkill(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	err := Uninstall("ghost-skill", "claude", false)
	if err == nil {
		t.Error("expected error for missing skill, got nil")
	}
}

func TestInstall_UnknownAgentFallback(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	defer func() { _ = os.Chdir(orig) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	path, err := Install("my-skill", "content", "unknown-agent-xyz", false)
	if err != nil {
		t.Fatalf("Install error: %v", err)
	}

	// should fall back to .pour/skills/<agent>/my-skill.md
	// use EvalSymlinks to handle macOS /var -> /private/var aliasing
	resolvedTmp, _ := filepath.EvalSymlinks(tmpDir)
	resolvedPath, _ := filepath.EvalSymlinks(path)
	expected := filepath.Join(resolvedTmp, ".pour", "skills", "unknown-agent-xyz", "my-skill.md")
	if resolvedPath != expected {
		t.Errorf("path = %q, want %q", resolvedPath, expected)
	}
}
