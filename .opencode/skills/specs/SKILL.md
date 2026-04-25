---
name: spec
description: Generate feature requirement specification documents. Triggers when user wants to develop a new feature, needs a spec, requirement analysis, or feature planning. This is the starting point for all development work.
---

Generate feature requirement specification documents.

## Input

The user describes what they want to do in one sentence (can be very rough).

## Steps

### 1. Requirements Interview (required)

Use the question tool to confirm each item one by one:

- **User story**: Who does what operation in what scenario, and what result do they expect?
- **Trigger mechanism**: Page entry, API call, CLI subcommand, scheduled task, webhook?
- **Core flow**: Happy path described step by step
- **Exception scenarios**: Failure, timeout, duplicate operations, concurrency?
- **Data**: What data is added or modified? Input/output formats?
- **Boundaries**: Explicitly state what is out of scope
- **Priority**: What's in the MVP? What goes to later iterations?

Wait for the user to answer each question before asking the next. Mark "unsure" items as TBD.

### 2. Code Research (automatic)

Use a Task subagent to research; don't make the user worry about it:
- Existing similar feature implementation patterns
- Related data models / type structures
- Route registration / command registration patterns
- Existing code that needs to be reused
- Project constraints (dependency restrictions, line count limits, Go version, etc.)

### 3. Generate Spec

Write to `docs/specs/{feature-name}.md`:

```markdown
# {Feature Name}

## Status
- Created: {date}
- Status: Draft

## Goal
(One sentence describing what problem to solve)

## Non-goals
- ...

## User Story
As {role}, I want to {action}, so that {value}

## Core Flow
1. ...

## Exception Handling
| Scenario | Handling |
|----------|----------|

## Technical Design

### Data Model
(structs, fields, types)

### API Endpoints / CLI Commands
| Method | Path/Command | Description |

### Implementation Steps (each step can be an independent commit)
1. [ ] Data model
2. [ ] Core logic
3. [ ] CLI command integration
4. [ ] Output rendering

### Existing Patterns Referenced
- {file path} — what was referenced

## Test Plan
- [ ] ...

## Open Questions
- ...

## MVP Scope
```

### 4. User Review

After generation, tell the user:
- Please review `docs/specs/{feature-name}.md`
- Confirm MVP scope and open questions
- Begin development after confirmation