package github

import (
	"testing"
)

func TestParseRepo(t *testing.T) {
	tests := []struct {
		input    string
		owner    string
		name     string
		ref      string
		wantErr  bool
	}{
		{"owner/repo", "owner", "repo", "", false},
		{"owner/repo@main", "owner", "repo", "main", false},
		{"owner/repo@v1.2.3", "owner", "repo", "v1.2.3", false},
		{"github.com/owner/repo", "owner", "repo", "", false},
		{"https://github.com/owner/repo", "owner", "repo", "", false},
		{"https://github.com/owner/repo@feat", "owner", "repo", "feat", false},
		{"noslash", "", "", "", true},
		{"", "", "", "", true},
	}

	for _, tc := range tests {
		owner, name, ref, err := parseRepo(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseRepo(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseRepo(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if owner != tc.owner || name != tc.name || ref != tc.ref {
			t.Errorf("parseRepo(%q) = (%q, %q, %q), want (%q, %q, %q)",
				tc.input, owner, name, ref, tc.owner, tc.name, tc.ref)
		}
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "frontmatter description",
			content: `---
name: my-skill
description: Does something useful
---
# My Skill`,
			want: "Does something useful",
		},
		{
			name: "no frontmatter, first non-heading line",
			content: `# My Skill

This is the description line.`,
			want: "This is the description line.",
		},
		{
			name: "truncates at 80 chars",
			content: `---
---
` + "This is a very long description that exceeds eighty characters and should be cut",
			want: "This is a very long description that exceeds eighty characters and should be cut",
		},
		{
			name:    "truncates long line with ellipsis",
			content: "This is a very long description line that definitely exceeds eighty characters total here",
			want:    "This is a very long description line that definitely exceeds eighty characters t...",
		},
		{
			name:    "empty content",
			content: "",
			want:    "No description",
		},
		{
			name:    "only headings",
			content: "# Title\n## Subtitle",
			want:    "No description",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDescription(tc.content)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseSkill(t *testing.T) {
	content := `---
name: ui-design
description: Pixel-perfect UI components
---
# UI Design`

	skill := parseSkill("ui-design", content)

	if skill.Name != "ui-design" {
		t.Errorf("Name = %q, want %q", skill.Name, "ui-design")
	}
	if skill.Description != "Pixel-perfect UI components" {
		t.Errorf("Description = %q, want %q", skill.Description, "Pixel-perfect UI components")
	}
	if skill.Content != content {
		t.Errorf("Content mismatch")
	}
}
