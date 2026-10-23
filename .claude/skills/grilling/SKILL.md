---
name: grilling
description: Grill the user relentlessly about a plan, decision, or idea. Use when the user wants to stress-test their thinking, or uses any 'grill' trigger phrases.
---

Interview the user relentlessly until you reach a shared understanding. Map this as a **design tree**: every decision branches into the decisions that hang off it.

Work the tree in **rounds**. The **frontier** is every decision whose prerequisites are already settled: the questions you can ask _now_ without guessing at answers you haven't heard yet. Ask the current frontier with the `question` tool. A round may contain one question or several related questions; choose the count from the decisions that are actually ready, never from a batch-size target. Then wait for the user's answers before the next round.

Build each tool question as follows:

- Use a short topic header such as `功能描述`, `触发入口`, `异常处理`, or `MVP 范围`. Never use `Q1`, `Q2`, or another sequence number as the header.
- Write a complete, self-contained question that explains the decision being made. Add a concrete example when it helps the user answer, especially for free-form requirements.
- Prefer concise, meaningful option labels over generic `是` / `否`. Use each option's description to explain its behavior or consequence.
- For a choice, provide the real, materially different alternatives the user might reasonably select. Never present only the recommended option.
- When one alternative is recommended, put it first and mark it as recommended in the option label. Keep the other viable alternatives visible and explain their tradeoffs; do not print a separate recommendation paragraph.
- Enable multiple selection when choices can coexist. Include `待定` only when deferring the decision is genuinely safe; explain what remains undecided.
- Keep custom answers available so the listed options do not constrain the user.

Use the interactive tool rather than rendering a questionnaire manually in chat. A good round should read like a compact set of product decisions, not a numbered interview transcript.

Each round the user answers reshapes the tree: settled decisions push the frontier outward and unblock questions that depended on them. Recompute the frontier and ask the next round. A question whose answer depends on another question still open in this round belongs to a _later_ round, not this one.

Finding _facts_ is your job, never the user's. When a frontier question needs a fact from the environment (filesystem, tools, etc.), dispatch a sub-agent to find it; don't ask the user for anything you could look up yourself. Don't block on it: a running exploration is an unsettled prerequisite, so only the questions downstream of it wait for the sub-agent to report; ask the rest of the frontier now. The _decisions_ are the user's: put each to them and wait.

The session is done when the frontier is empty: every branch of the design tree visited, nothing left silently assumed. Do not act on it until the user confirms you have reached a shared understanding.
