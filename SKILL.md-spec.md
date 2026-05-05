# SKILL.md Specification

A `SKILL.md` file defines a reusable skill that can be installed into any AI agent using `pour`.

## File Location

A `SKILL.md` can live:
- In the **root** of a repo — for single-skill repos
- Inside a **named subfolder** — for multi-skill repos

```
# single skill
my-repo/
└── SKILL.md

# multiple skills
my-repo/
├── code-reviewer/
│   └── SKILL.md
└── email-drafter/
    └── SKILL.md
```

## Format

A `SKILL.md` is a Markdown file with optional YAML frontmatter.

```markdown
---
name: code-reviewer
description: Reviews code for bugs, readability, and best practices
version: 1.0.0
tags:
  - coding
  - review
agents:
  - claude
  - cursor
---

# Code Reviewer

You are an expert code reviewer...
```

## Frontmatter Fields

| Field | Required | Description |
|---|---|---|
| `name` | Yes | Unique identifier, lowercase with hyphens |
| `description` | Yes | One-line summary shown in `pour list` and `pour search` |
| `version` | No | Semver string e.g. `1.0.0` |
| `tags` | No | List of keywords for search |
| `agents` | No | List of agents this skill is optimized for |

## Body

Everything after the closing `---` of the frontmatter is the skill content — the instructions, prompt, or rules that get installed into the AI agent.

Write it as if you're directly instructing the AI:

```markdown
---
name: sql-expert
description: Writes optimized SQL queries for PostgreSQL
version: 1.0.0
tags: [sql, database, postgres]
agents: [claude, cursor]
---

You are an expert PostgreSQL engineer.

When asked to write SQL:
- Always use CTEs for complex queries
- Add indexes where appropriate
- Explain query plans when relevant
```

## Naming Rules

- `name` must be **lowercase**
- Use **hyphens** not underscores: `code-reviewer` not `code_reviewer`
- Must be **unique** within the same repo
- Folder name should match `name` in the frontmatter

## Publishing

To share your skill, push your repo to GitHub. Users install it with:

```bash
pour add your-username/your-repo -a claude
```

No registration needed. GitHub is the registry.
