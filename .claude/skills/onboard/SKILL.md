---
name: onboard
description: 帮助新成员了解 InkSpace 的架构、开发环境、测试和当前工作流。用户要求项目介绍、新人上手或寻找第一个开发入口时使用；只做导航，不代替需求或实现技能。
---

# Onboard

依据当前仓库提供简短、按用户目标定制的项目导航，不启动服务或修改环境。

1. 阅读 `CLAUDE.md`、`README.md`、`docs/agents/WORKFLOW.md`，以及用户目标对应的 `.claude/rules/` 文件。涉及命令或测试时核对对应入口、脚本和测试文件，不把历史描述当成当前事实。
2. 简要介绍三个 Go 入口 `cmd/server`、`cmd/admin`、`cmd/scheduler` 共享 `internal/`；请求沿 router → handler → service → models；`web/blog` 和 `web/admin` 是独立的 Vue 应用；开发依赖 MySQL/Redis。按需给出 README 中对应服务的启动方式和现有测试命令。
3. 根据目标指出一个可以开始阅读的现有路径和相关规则，不重复整份 `CLAUDE.md`。日常任务选路看 `docs/agents/WORKFLOW.md` 或 `/ask-matt`；要实际启动并验证应用，用 `/run`；要构建功能，用相应的工程技能。不要推荐已经删除的旧命令。
4. 如果用户尚未说明想做的方向，只提供架构地图及几个可选入口，不把“新人介绍”变成需求访谈或未经请求的开发工作。

已有 Go 单测，不能声称仓库没有测试。所有路径、端口与依赖以当前文件为准。
