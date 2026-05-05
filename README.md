<div align="center">

<p align="center"><img src="docs/logo.svg" alt="pour" width="420"/></p>

**Install skills for any AI agent in one command.**

[![Latest Release](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2FDev-kitx%2Fpour-cli%2Fmain%2F.github%2Fbadges%2Frelease.json&style=flat-square)](https://github.com/Dev-kitx/pour-cli/releases/latest)
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-9B59B6?style=flat-square)](LICENSE)
[![Open Source](https://img.shields.io/badge/open%20source-yes-FF6B9D?style=flat-square)](#)

```sh
brew install pour
```

[**Docs & Demo →**](https://dev-kitx.github.io/pour-cli)

</div>

---

## What is pour?

`pour` is an open-source CLI — built like `brew`, for AI skills.

A **skill** is a `SKILL.md` file in any public GitHub repo. It contains instructions, prompts, or context that teaches an AI agent how to do something — write code in your style, follow your conventions, use your stack. `pour` finds those files and installs them into the right place for whatever agent you're using.

**GitHub is the registry.** No central server. No signup. Just repos.

---

## Install

```sh
# Homebrew (macOS / Linux)
brew install pour

# Build from source
git clone https://github.com/Dev-kitx/pour-cli
cd pour-cli && go build -o pour . && sudo mv pour /usr/local/bin/
```

---

## Usage

```sh
# Install all skills from a repo
pour add owner/repo -a claude

# Install a specific skill
pour add owner/repo -s ui-design -a cursor

# Install globally (all projects)
pour add owner/repo -a windsurf -g

# List installed skills
pour list
pour list -a claude

# Remove a skill
pour remove ui-design -a claude

# Show skill details
pour info ui-design

# Re-fetch and update all skills
pour update

# Search GitHub for skills
pour search "typescript testing"

# Create a new SKILL.md to share
pour init my-skill

# Add a GitHub token (higher rate limits)
pour auth login
```

---

## Supported Agents

| Agent          | Flag                              |
|----------------|-----------------------------------|
| Claude Code    | `claude` / `claude-code`          |
| Cursor         | `cursor`                          |
| Windsurf       | `windsurf`                        |
| GitHub Copilot | `github-copilot`                  |
| Cline          | `cline`                           |
| Roo Code       | `roo`                             |
| OpenCode       | `opencode`                        |
| OpenHands      | `openhands`                       |
| Codex CLI      | `codex`                           |
| Augment        | `augment`                         |
| Amp            | `amp`                             |
| Tabnine        | `tabnine-cli`                     |
| Continue       | `continue`                        |
| Gemini CLI     | `gemini-cli`                      |
| **+50 more**   | [see full list →](docs/spec.html) |

---

## How skills work

A skill is a single `SKILL.md` file with YAML frontmatter:

```markdown
---
name: ui-design
description: Pixel-perfect UI components with design system conventions
version: 1.0.0
authors: [your-name]
tags: [ui, design, react]
agents: [claude, cursor, windsurf]
---

# UI Design Skill

When building UI components, always...
```

`pour add owner/repo` scans the repo root and subdirectories for every `SKILL.md` automatically. No fixed folder structure required.

```
your-skills-repo/
├── ui-design/
│   └── SKILL.md        ✓ found
├── api-patterns/
│   └── SKILL.md        ✓ found
└── skills/
    └── testing/
        └── SKILL.md    ✓ found (2 levels deep)
```

→ Full specification: [`SKILL.md Spec`](docs/spec.html)

---

## Where skills are installed

Skills land in the agent's native context folder so they're picked up automatically:

| Agent       | Project path       | Global path          |
|-------------|--------------------|----------------------|
| Claude Code | `.claude/skills/`  | `~/.claude/skills/`  |
| Cursor      | `.cursor/rules/`   | `~/.cursor/rules/`   |
| Windsurf    | `.windsurf/rules/` | `~/.windsurf/rules/` |
| Cline       | `.cline/rules/`    | `~/.cline/rules/`    |

---

## Features

- **Any GitHub repo** — point `pour` at any public repo, it finds `SKILL.md` files automatically
- **Interactive picker** — when a repo has multiple skills, a terminal UI lets you choose
- **Multi-agent** — one skill, any agent, the `-a` flag handles routing
- **No auth required** — works with public repos out of the box (60 req/hr)
- **Global or local** — install per-project or globally with `-g`
- **Update tracking** — `pour update` re-fetches all skills and shows a colored diff of what changed
- **Single binary** — no runtime, no deps, just a fast Go binary

---

## Share your skills

```sh
# 1. Create a SKILL.md
pour init my-skill

# 2. Push to GitHub
git add . && git commit -m "add skill" && git push

# 3. Anyone can install it
pour add your-username/your-repo -a claude
```

---

## Contributing

### Clone and build

```sh
git clone https://github.com/Dev-kitx/pour-cli
cd pour-cli
go build -o pour .
```

### Run tests

```sh
go test ./...
```

### Run locally

```sh
go run . add --help
go run . search react
go run . version
```

### Build release binary

```sh
go build -ldflags "-s -w -X github.com/Dev-kitx/pour-cli/cmd.Version=dev" -o pour .
```

PRs welcome. Open an issue to discuss before large changes.

---

<div align="center">
  <sub>Built with Go · <a href="https://github.com/spf13/cobra">Cobra</a> · <a href="https://github.com/charmbracelet/bubbletea">Bubbletea</a> · <a href="https://github.com/charmbracelet/lipgloss">Lipgloss</a></sub>
</div>
