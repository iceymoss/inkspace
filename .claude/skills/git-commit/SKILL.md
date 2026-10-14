---
name: git-commit
description: 为 InkSpace 草拟提交信息，或在用户明确要求提交时核对变更边界、验证并创建 Git 提交。只负责提交卫生，不替代实现或代码审查。
---

# Git Commit

此技能只处理**提交阶段**。！！！注意！！！：commit不需要你的署名，这个非常重要，仅要求拟写 commit message 时，不暂存、不提交；若请求包含推送或 PR，那是额外的对外操作，需分别核对权限。

1. 检查 `git status --short --branch`、`git diff`、`git diff --cached`、`git log -8 --oneline`。区分本任务修改与原有变更、生成文件及敏感数据；看不清提交边界时先询问，不使用 `git add .`、`git add -A` 或 `git commit -a`。
2. 根据变更选择简洁的英文 Conventional Commit 标题：`<type>[(scope)]: <imperative description>`；常用 `feat`、`fix`、`docs`、`chore`、`refactor`、`test`、`ci`，scope 仅在能澄清范围时使用。一个提交对应一个逻辑变更；需要时在正文说明动机，不编造 issue 或测试结果。
3. 如果用户明确要求创建提交，先确认所需的相关检查已运行并如实报告失败或未运行项；在默认分支上先创建任务分支。仅精确暂存本次提交的文件，再检查 `git diff --cached --check`、`git diff --cached` 和暂存路径；不要把用户无关的已暂存内容带入提交。无法安全区分时先确认。
4. 按当前会话及项目的归属要求结束提交信息；若没有更高优先级的不同指示，末尾使用 `Co-Authored-By: Claude Code <noreply@anthropic.com>`。不要恢复旧技能的“禁止 AI 署名”规则。不擅自 `--amend`、`--no-verify`、空提交、推送或创建 PR；hook 失败时先报告并处理原因，而非绕过。
5. 提交后查看 `git status --short` 与提交标题/哈希，汇报验证结果和剩余变更。

`/implement` 自身也规定了提交步骤；只在用户单独要求提交管理或拟写信息时用本技能，不在每次实施中重复跑完整流程。
