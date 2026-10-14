# InkSpace × Matt Pocock Skills：开发工作流使用手册

这份手册回答「现在该用哪个技能、下一步是什么、结果写到哪里」。InkSpace 将 [Matt Pocock Skills](https://www.aihero.dev/skills) 作为**可组合的工程技能**，不是每个任务都必须跑完的流水线。遇到不确定的入口，先用 `/ask-matt` 说明现状和目标；它是路由建议，不会替你自动完成整条流程。具体执行以本仓库 [`.claude/skills/`](../../.claude/skills/) 中的技能正文和 [项目规则](../../CLAUDE.md) 为准。

本仓库安装的是 [上游提交 `49dd158d`](https://github.com/mattpocock/skills/tree/49dd158d1076134a641b33efb035946536778336) 的 20 个 engineering 技能及主流程依赖 `/grilling`、`/handoff`，另有三个 InkSpace 本地辅助技能（见下文），不是自动更新的全局插件。上游 `/ask-matt` 还介绍了其他未安装的独立入口（如 `/grill-me`、`/to-questionnaire`、`/wait-what`、`/teach`）；不要把它们视为本项目可直接调用的命令。

## 一眼选路

| 你遇到的场景 | 推荐入口 | 何时转入下一步 |
| --- | --- | --- |
| 想法或需求还模糊，但一轮讨论能理清 | `/grill-with-docs` | 讨论中定下术语/决策；小改动直接 `/implement`，跨会话工作转 `/to-spec` |
| 清楚的小功能或具体修复 | `/implement`；只想先写一个行为的测试用 `/tdd` | 同意测试边界后实现、测试、审查 |
| 多会话新功能 | `/grill-with-docs` → `/to-spec` → `/to-tickets` | 批准纵向切片和依赖关系后，逐票 `/implement` |
| 可并行的整份规格 | `/implement-spec` | 已有 spec 和互相标明阻塞关系的 tickets，且同意创建集成分支/worktree |
| 外来 bug 报告或需求，尚未整理 | `/triage` | 分类、核实并经维护者确认后，形成可实施的任务；不要再 triage `/to-tickets` 产出的票 |
| 难复现的故障或回归 | `/diagnosing-bugs` | 先取得能对**这个故障**变红的快速反馈，再定位并补回归测试 |
| 大到无法在一次会话中厘清方向 | `/wayfinder` | 逐个解决决策票；方向明确后回 `/to-spec`，不直接跳到写代码 |
| 模块边界、可测试性或架构维护 | `/improve-codebase-architecture`、`/codebase-design` | 确定候选和测试 seam，再回需求流程 |
| 需要试出 UI/状态模型、查外部资料 | `/prototype`、`/research` | 将验证过的结论带回需求讨论；跨目录/会话用 `/handoff` |
| 实现完成或准备 PR | `/code-review`、`/pr`、`/retro` | 指定审查基点；验证后按权限决定是否提交、开 PR、复盘 |

## InkSpace 本地辅助技能

这三个入口由本项目维护，不属于上游的 22 个固定版本技能；它们处理主开发流程之外的工作，不另起一套需求 → 实施流程。

| 场景 | 入口 | 边界 |
| --- | --- | --- |
| 新会话或 `/clear` 后想知道仓库进度 | [`/catchup`](../../.claude/skills/catchup/SKILL.md) | 只读恢复可验证状态；要**为下一会话生成**交接文件则用 `/handoff` |
| 新成员想了解项目、找到开发入口 | [`/onboard`](../../.claude/skills/onboard/SKILL.md) | 按当前 README/规则介绍架构；具体需求仍用 `/ask-matt` 选路 |
| 只想拟写提交信息，或已完成工作后明确要求提交 | [`/git-commit`](../../.claude/skills/git-commit/SKILL.md) | 不代替 `/implement` 的测试/审查；仅在获准时提交，不自动推送或开 PR |

例如，以下每行都是**独立输入**，不是一条要整段执行的命令：

```text
/catchup 我刚切回 feat-workflow-mb，请根据当前分支和未提交变更告诉我已确认的进度；不要修改文件。
/onboard 我想从 web/blog 的文章列表开始贡献，请指给我架构、相关规则和测试入口。
/git-commit 只根据当前已完成的变更草拟提交信息，先不要暂存或提交。
```

本地启停和实际运行验证用已有的 `/run`，具体端口与命令见 [README 快速开始](../../README.md#-快速开始)；Go 测试按任务范围直接运行 `go test`。不用恢复旧的 `/run-services` 或 `/gotest`。需要选工程流程时用 `/ask-matt`，而不是把这些辅助技能当成自动编排器。

## 主流程：从想法到交付

```text
/grill-with-docs
    ├─ 一次会话能做完 → /implement
    └─ 多会话工作 → /to-spec → /to-tickets
                                ├─ 逐票 /implement（建议每票新上下文）
                                └─ /implement-spec（集成分支 + 可并行 worktree）
                  实现时尽可能 /tdd → /code-review → 按需 /pr → /retro
```

1. **澄清。** 例如输入 `/grill-with-docs 在 web/admin 增加知识库审核筛选`。它会调用 `/grilling` 逐轮追问，也会通过 `/domain-modeling` 按需记录术语和长期决策。用户负责确认目标、范围、数据及异常场景；Agent 负责查清现有代码事实。需要可运行的答案（页面外观、复杂状态流）时，先用 `/handoff` 转到原型会话、运行 `/prototype`，再把验证结论带回来。
2. **沉淀规格。** 多会话工作用 `/to-spec` 将**当前已经讨论好的内容**综合为规格；它不会替代前面的需求访谈，会在写入前就测试 seam 与用户对齐。本地发布位置是 `.scratch/<feature-slug>/spec.md`。短小、明确的需求无需机械写 spec。
3. **拆可交付的票。** `/to-tickets .scratch/<feature-slug>/spec.md` 提出贯穿数据/API/UI/测试的窄纵向切片（适用时），列出 `Blocked by` 依赖。先让用户判断粒度和阻塞边，再写入 `.scratch/<feature-slug>/issues/01-<slug>.md` 等**一票一文件**；不是按 model → service → handler → 页面简单横切。已经能直接实施的小任务可以略过。
4. **实施。** 可以 `/implement .scratch/<feature-slug>/issues/01-<slug>.md`，先完成阻塞票，再在新上下文处理下一票。它会在可用的测试边界调用 `/tdd`，定期运行检查，最后调用 `/code-review`，**原技能还要求向当前分支提交**。若票据确能独立推进且用户允许分支、worktree、子代理和合并，才考虑 `/implement-spec .scratch/<feature-slug>/spec.md`；后者会创建集成分支，并行处理可开始的票，最终统一审查。InkSpace 后端三个入口共用 `internal/` 和数据库，涉及同一模型/路由的票往往有真实依赖，不要为了并行而虚构独立性。
5. **复核并收尾。** 单独运行 `/code-review master` 时，`master` 是比较基点（应按实际目标分支/提交调整）；它分别检查仓库规范和来源规格，没有规格就只报告 Standards。修复问题、验证行为后可用 `/pr` 起草 PR 说明，复盘时用 `/retro` 寻找下次可复用的检查或文档改进。`/retro` 的上游第一步要求调用未安装的 `writing-for-agents`：**目前不能保证原样完整运行**；若真要用，先另行决定补装同版本依赖，或明确以人工方式按其现有检查项降级复盘，并说明跳过了依赖。

尽量让澄清 → spec → tickets 留在同一上下文；每个独立实施票再开启新会话。会话过长时在阶段边界 `/compact`，确实要换目录、换工具或交接给其他人时才用 `/handoff`。这些是指导，不是强制清空上下文的规则。

## 几个常见的 InkSpace 实例

以下代码块是**在 Claude Code 输入框中输入的示例**，不是 shell 命令；同一代码块的多行表示**每一步结束、检查产物并确认后，再输入下一条**，不要整块粘贴执行。示例中的 `.scratch/...` 路径只表示预期产物，实际名称以技能生成的文件为准。调用 `/implement` 前先说明是否允许提交：其上游技能最后要求 commit；若本次不允许提交，应先明确要求只实现与验证、不提交，并预期不能原样完成该技能最后一步。以下示例只是文档，不会实际调用技能或授权写代码。

### feat：大型新功能（目标可讨论，但需要多轮交付）

例如为知识库加入后台审核流程，涉及 `web/admin`、`cmd/admin` 和共享 `internal/`。先在同一会话中完成需求讨论、规格和拆票；不要一次粘贴三行：

```text
/grill-with-docs 规划知识库文档审核：哪些角色提交和审核？退回后如何修改？管理员页面和博客展示各有什么变化？请先梳理边界，不开始实现。
/to-spec 根据刚才确定的审核流程整理规格；请提出测试 seam 让我确认，发布到本地 Markdown tracker。
/to-tickets .scratch/knowledge-base-review/spec.md 请按可独立验证的纵向行为拆票，并先让我确认票的粒度和 Blocked by 关系。
```

然后**二选一**：依赖紧密时按阻塞顺序逐票开启新会话，例如 `/implement .scratch/knowledge-base-review/issues/01-<实际文件名>.md`；若票据确实可独立并行、且已允许建集成分支/worktree/合并，则用 `/implement-spec .scratch/knowledge-base-review/spec.md`。不要同时运行两种实施方式。破坏性数据库变更先对齐并准备迁移 SQL。

### feat：小型新功能（有少量待确认细节）

例如在博客文章列表增加「只看已收藏」筛选。先确认筛选语义，再判断一次会话能否完成；能完成就不必造 spec 和 tickets：

```text
/grill-with-docs 博客文章列表增加「只看已收藏」筛选：请确认未登录时的行为、筛选与分页的组合；先不要写代码。
/implement 按刚确认的范围实现「只看已收藏」筛选；沿用 web/blog 与 cmd/server 现有模式，先确认测试 seam，再实现和验证。
```

### feat：已经很清晰的新功能

例如用户已给定「后台友链列表新增按状态过滤」，请求/响应、默认值和验收行为都已明确。可以直接进入实现，不必重复访谈；想明确要求测试先行，也可单独调用 `/tdd`：

```text
/implement 在 web/admin 的友链列表增加状态筛选：默认显示全部，切换后重置到第一页，保留现有接口响应结构；请先核对仓库模式并确认测试 seam。
```

`/tdd 为友链状态筛选先写一个失败的行为测试，确认测试 seam 后按 red → green 逐个行为推进。` 是**替代的测试优先入口**，不是运行 `/implement` 后必须再手动调用的一步；`/implement` 本身会在可行时调用 `/tdd`。

### feat：仅后端功能

例如管理后台新增「禁用指定用户的分享链接」接口：

```text
/implement 为 cmd/admin 增加禁用指定用户知识库分享链接的接口：只有管理员可操作；先检查 internal/router、handler、service、models 的现有模式，确认鉴权、响应语义和测试 seam 后实现。
```

明确目标入口是 `cmd/admin`，业务逻辑放共享的 `internal/service`，统一走 `utils.*` 响应；若需求还不清楚，先 `/grill-with-docs`，不要猜接口行为。Go 测试优先选择不依赖真实 MySQL/Redis 的现有方式。

### feat：仅前端功能

例如博客首页为已有文章列表增加客户端的「返回顶部」按钮，不新增 API：

```text
/implement 只修改 web/blog：为文章列表增加滚动后出现的「返回顶部」按钮；沿用 Vue 3 <script setup> 和 Element Plus，保持中文文案，不改 web/admin 或后端接口。先确认交互与可验证行为。
```

若功能涉及接口或跨应用状态，先查清是否真的“仅前端”；不要为了贴合示例而绕开已有 axios 封装或 Pinia。

### feat：UI 设计方向尚未确定

例如希望改版博客文章详情页，但尚不确定目录、正文与作者信息的排布：

```text
/grill-with-docs 我想重做 web/blog 的文章详情页。请先讨论阅读目标、必须保留的交互和移动端约束，不改生产页面。
/prototype 在 web/blog 探索文章详情页的几种明显不同的 UI 布局，让我实际比较目录、正文与作者信息；原型仅用于验证设计，不直接替换正式页面。
```

`/prototype` 会创建可运行的临时代码，原技能还要求保存到原型分支并提交；应先明确分支与提交授权。确定方向后把结论带回原讨论（跨会话/目录时用 `/handoff`），再视规模用 `/to-spec` 或 `/implement` 落地，而不是把原型直接当成成品。

### fix：复现方式和原因都明确的 bug

例如文章列表在传入非法页码时返回错误分页数据：

```text
/implement 修复 cmd/server 文章列表在 page=0 时的分页行为：预期按第一页处理，当前返回空列表；先在现有测试 seam 补一个能复现的回归用例，再做最小修复并运行相关 Go 测试。
```

只需单项行为且想显式测试先行时，可改用 `/tdd 为 page=0 的分页行为先写失败测试，确认 seam 后做最小修复。`；无需为已确定的修复重新建一套 spec/tickets。

### fix：症状不明确，需要 Agent 协助排查

例如用户偶尔看不到刚发布的文章，但尚不知道是列表过滤、缓存还是时序问题：

```text
/diagnosing-bugs 用户反馈「刚发布的文章偶尔不在博客列表出现」。请先调查 cmd/server、internal/service 的列表和缓存路径，找出能针对这个故障变红的快速反馈；在确认原因前不要猜修复方案，也不要直接修改业务逻辑。
```

先拿到可重复的失败信号，再缩小范围、验证假设并补回归测试；如果问题已经明确，不必调用诊断技能。排查结果若揭示功能需求或设计决策，再回 `/grill-with-docs` → `/to-spec`，不要把未经证实的猜测直接写入规格。

### 外来反馈、超大任务与架构维护

外部反馈先存成可引用的本地 issue，再输入 `/triage .scratch/incoming-reports/issues/01-<实际文件名>.md`；先核实、建议 `bug`/`enhancement` 类别与状态并等待维护者确认，**不要**再 triage `/to-tickets` 产出的票。若任务是「重新梳理知识库权限与版本模型」且连方案都说不清，可用 `/wayfinder 规划知识库权限与版本模型的决策地图，先明确目标和未知问题，不直接改代码。`；在 `.scratch/<effort>/map.md` 按会话解决决策票，方向明确后回 `/to-spec`。

需要查外部 API 事实时用 `/research 调研目标 API 的官方文档与限制，保存带出处的笔记供后续设计讨论。`；要发现项目模块边界问题用 `/improve-codebase-architecture`，选定候选后可借 `/codebase-design` 比较接口和测试 seam。人工才能完成的凭据、仪表盘或迁移步骤，可请 `/wizard` 生成向导，但先审阅脚本及其对 `.env`/外部系统的影响。

## 本地配置、产物与权限

- **配置已就绪。** [`CLAUDE.md`](../../CLAUDE.md#agent-skills) 已指向 [本地任务记录](issue-tracker.md)、[状态词汇](triage-labels.md) 和 [领域文档约定](domain.md)；不必每次运行 `/setup-matt-pocock-skills`。要切换 tracker 或重做配置时再用它。`.scratch/` 尚未创建，且**未被 `.gitignore` 忽略**；新产物是否纳入版本控制须显式决定，不能假定它是自动清理的临时文件。
- **历史规格并未迁移。** [`docs/specs/`](../specs/) 里的业务规格可供参考，但新 `/to-spec` 的默认发布目录是 `.scratch/<feature-slug>/spec.md`。本地任务状态写在文件开头的 `Status:`，五种默认状态见 [triage 映射](triage-labels.md)；`/wayfinder` 用 `.scratch/<effort>/map.md` 和 `issues/NN-<slug>.md`，其 `claimed/resolved` 是决策票状态，不是普通实现票的五态。外部 GitHub issue 不会因配置本地 tracker 而自动同步；本地 `bug/enhancement` 分类及关闭语义也没有专门的字段/平台操作规范，使用 `/triage` 时应在本地文件明确记载并与用户对齐。
- **项目规则优先。** 三个 Go 入口 `cmd/server`（博客）、`cmd/admin`（后台）、`cmd/scheduler`（定时任务）共享 `internal/`；请求沿 router → handler → service → models；两个独立 Vue 应用在 `web/blog`、`web/admin`，包管理用 pnpm、界面中文单语。详见 [架构说明](../../CLAUDE.md)、[Go](../../.claude/rules/go-backend.md)、[数据库](../../.claude/rules/database-model.md)、[前端](../../.claude/rules/frontend.md)、[调度器](../../.claude/rules/scheduler.md)、[测试](../../.claude/rules/test-files.md)，不要让上游通用示例覆盖实际约定。
- **检查与副作用。** Go 常用 `go test ./... -count=1`、`go vet ./...`、`go build ./...`；各前端目录执行 `pnpm build`，博客可执行 `pnpm test` 和 `pnpm lint`。后台的 `pnpm lint` 实际带 `--fix` 且扫描 `..`，可能改动其他文件：先确认范围/权限，再运行或改用不写文件的检查。数据库集成测试须隔离测试库。测试通过不等于授权提交或对外发布：branch/worktree、commit、push、issue/label 更新和 PR 均需按当次请求与权限执行。`/pr` 是 PR 文案规范，不等于已经创建 PR；对外发布前再次确认。

## 参考资料

- [AI Hero — Skills](https://www.aihero.dev/skills)：上游面向用户的用途与主流程介绍。
- [mattpocock/skills 固定版本 README](https://github.com/mattpocock/skills/blob/49dd158d1076134a641b33efb035946536778336/README.md)：安装方式、技能分类及注意事项。
- [本地 `/ask-matt` 路由说明](../../.claude/skills/ask-matt/SKILL.md)：场景和阶段边界；执行前仍应阅读对应技能的 `SKILL.md`，不要只依据概览推断副作用。
