---
name: implement-spec-in-seq
description: "Implement the result of /to-spec and /to-tickets in code in sequence."
disable-model-invocation: true
---

You have been provided a spec. This spec should have tickets associated with it, describing how to implement the spec.

The issue tracker should have been provided to you. If not, tell the user to run `/setup-matt-pocock-skills`.

The goal is the entire spec implemented on a single **integration branch**, with every ticket resolved the way the issue tracker closes work.

The tickets are a list of steps. They may have blocking relationships between them. They are meant to be implemented in sequence. The one that blocks the others should be implemented first. This also means the **frontier** of tickets are the ticket 1) not implemented 2) with the smallest ticket number.

Communication to and from subagents should be sparse. Communicate primarily through **context pointers**: to the spec, tickets, research notes, and previous commits. Don't duplicate information already available via pointers.

**Implementer subagents** should be run in the background.

**Code reviewer subagents** should be run in the background.

## Steps

1. Read the spec and tickets to understand the list of tasks.

2. (optional) Use an **exploration subagent** to conduct any exploration required by the tickets - relevant codebase files or external documentation. Ensure the exploration subagent can save files - it should save its markdown notes in a directory outside the repo, accessible by all future subagents. This lets **implementer subagents** focus on implementation rather than exploration.

3. Create the integration branch. If the issue tracker closes work through PRs, or the user asks for one, open a draft PR after the first merge in step 6 (a branch with no commits ahead of main can't open one), marked as closing the spec and tickets.

4. Use **implementer subagents** to implement each ticket, each in its own worktree on its own branch. Each implementer subagent:
  - confirms its worktree is based on the integration branch before starting, and resets onto it if not;
  - calls the Skill tool with `tdd` to build the ticket;
  - does not commit the changes before code review by another subagent is performed;
  - only commits the changes after addressing the code review findings;
  - reports that the first pass of the implementation is done

5. Once an **implementer subagent** reports that the first pass of the implementation is done, use **code reviewer subagent** and instruct **code reviewer subagent** to call the Skill tool with `code-review` to review the ticket. And pass the code review findings back to **implementer subagent** to improve the code base.

6. Once an **implementer subagent** completes, instruct the **implementer subagent** to merge the integration branch tip into its own branch and report ticket implementation is completed.

7. Merge the work in step 6 to the integration branch with a **merger subagent**.

8. Once all tickets are complete, call the Skill tool with `code-review` on the integration branch. Fix all issues raised by the code review in a single **implementer subagent**.

9. If a draft PR exists, mark it ready for review. Otherwise, resolve each ticket the way the issue tracker closes work, and report the integration branch.

10. Clean up all **implementer subagent** worktrees.
