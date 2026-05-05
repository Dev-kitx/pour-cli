package config

import (
	"testing"
)

func TestSkillsDB_AddSkill(t *testing.T) {
	db := &SkillsDB{}

	s1 := InstalledSkill{Name: "ui-design", Agent: "claude", Repo: "acme/skills"}
	db.AddSkill(s1)

	if len(db.Skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(db.Skills))
	}
	if db.Skills[0].Name != "ui-design" {
		t.Errorf("Name = %q, want %q", db.Skills[0].Name, "ui-design")
	}

	// adding a different agent creates a second entry
	s2 := InstalledSkill{Name: "ui-design", Agent: "cursor", Repo: "acme/skills"}
	db.AddSkill(s2)
	if len(db.Skills) != 2 {
		t.Fatalf("expected 2 skills after different agent, got %d", len(db.Skills))
	}

	// adding same name+agent updates in place
	s1Updated := InstalledSkill{Name: "ui-design", Agent: "claude", Repo: "acme/skills-v2"}
	db.AddSkill(s1Updated)
	if len(db.Skills) != 2 {
		t.Fatalf("expected still 2 skills after update, got %d", len(db.Skills))
	}
	if db.Skills[0].Repo != "acme/skills-v2" {
		t.Errorf("Repo = %q, want %q", db.Skills[0].Repo, "acme/skills-v2")
	}
}

func TestSkillsDB_RemoveSkill(t *testing.T) {
	db := &SkillsDB{
		Skills: []InstalledSkill{
			{Name: "ui-design", Agent: "claude"},
			{Name: "testing", Agent: "claude"},
		},
	}

	removed := db.RemoveSkill("ui-design", "claude")
	if !removed {
		t.Error("RemoveSkill: expected true, got false")
	}
	if len(db.Skills) != 1 {
		t.Fatalf("expected 1 skill remaining, got %d", len(db.Skills))
	}
	if db.Skills[0].Name != "testing" {
		t.Errorf("remaining skill = %q, want %q", db.Skills[0].Name, "testing")
	}

	// remove non-existent
	removed = db.RemoveSkill("ui-design", "claude")
	if removed {
		t.Error("RemoveSkill: expected false for missing skill, got true")
	}

	// wrong agent doesn't remove
	db.AddSkill(InstalledSkill{Name: "testing", Agent: "cursor"})
	removed = db.RemoveSkill("testing", "windsurf")
	if removed {
		t.Error("RemoveSkill: expected false for wrong agent, got true")
	}
	if len(db.Skills) != 2 {
		t.Errorf("expected 2 skills unchanged, got %d", len(db.Skills))
	}
}

func TestKnownAgents(t *testing.T) {
	agents := KnownAgents()
	if len(agents) == 0 {
		t.Fatal("KnownAgents returned empty list")
	}

	// spot-check a few canonical names
	must := []string{"claude", "cursor", "windsurf"}
	set := make(map[string]bool, len(agents))
	for _, a := range agents {
		set[a] = true
	}
	for _, name := range must {
		if !set[name] {
			t.Errorf("KnownAgents missing %q", name)
		}
	}
}

func TestAgentPathsConsistency(t *testing.T) {
	// every key in AgentPaths should have a non-empty path
	for agent, path := range AgentPaths {
		if path == "" {
			t.Errorf("AgentPaths[%q] is empty", agent)
		}
	}
	for agent, path := range AgentGlobalPaths {
		if path == "" {
			t.Errorf("AgentGlobalPaths[%q] is empty", agent)
		}
	}
}
