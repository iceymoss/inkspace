# 知识库文档只读详情与多文件类型支持

## 状态
- 创建日期: 2026-08-13
- 状态: 开发中（步骤 1-6 已完成）
- 关联规格: `docs/specs/workspace-knowledge-base.md`

## 目标
将知识库从“点击 Markdown 文档直接编辑”升级为统一的文件型知识库：所有条目默认进入权限受控的只读详情页，有编辑权限的成员通过明确的“编辑”按钮修改文本类文件；同时支持主流 Office、PDF、图片、文本、代码、结构化数据及压缩文件的上传、安全预览、下载和版本冲突检测。

## 非目标
- 不做 OT、CRDT、多人光标或实时协同编辑。
- 不在浏览器内原生编辑 Office、PDF、图片和压缩文件；这些类型只支持预览、下载和有权限时替换原件。
- 不执行上传的脚本、宏、二进制程序或压缩包内容。
- 不允许在页面中直接运行 HTML、SVG、JavaScript 等主动内容；文本/代码一律转义后展示。
- 不在 MVP 中做压缩包内文件的递归解压与逐项预览，只展示安全生成的只读目录清单。
- 不引入微服务、消息队列、Kafka 或 WebSocket；转换任务由现有 `cmd/scheduler` 和 MySQL 持久化任务承载。
- 不新增知识库管理后台；本期入口位于 `web/blog`。全局系统管理员是否可绕过知识库权限不在本期范围。
- 不改变现有 Wiki 发布、发布到博客和 token 分享的产品语义；它们只适配新的详情与文件类型能力。
- 不支持音频、视频、CAD、PSD、电子书等未列入 MVP 白名单的专用在线预览器；未知类型保留受控下载能力。
- 不对存量数据库做删列、改类型或新增唯一约束等破坏性迁移。

## 用户故事
- 作为公开知识库访客，我想点击文档后直接阅读只读内容，以便不登录也能浏览已公开、已发布的知识。
- 作为私有知识库成员，我想在权限允许时查看私有文档，并且不会因为一次点击误入编辑态或误改内容。
- 作为知识库编辑者，我想从详情页明确点击“编辑”，保存后回到最新只读预览，以便区分阅读和创作场景。
- 作为知识库管理员或所有者，我想上传、替换、发布和管理多种文件，以便知识库不局限于 Markdown。
- 作为文本和代码文件编辑者，我想在站内修改源内容，并在内容已被他人更新时收到冲突提示，以免覆盖他人的修改。
- 作为知识库读者，我想在线查看 Office、PDF、图片、代码、JSON、YAML 和压缩包目录，以便尽量减少下载后再打开的操作。

## 角色与权限

### 知识库角色
| 角色 | 私有库查看 | 新建/上传 | 编辑文本 | 替换文件 | 发布/分享 | 成员管理 | 删除空间 |
|------|------------|-----------|----------|----------|-----------|----------|----------|
| owner | 是 | 是 | 是 | 是 | 是 | 是 | 是 |
| admin | 是 | 是 | 是 | 是 | 是 | 是 | 否 |
| editor | 是 | 是 | 是 | 是 | 是 | 否 | 否 |
| viewer | 是 | 否 | 否 | 否 | 否 | 否 | 否 |
| public visitor | 仅公开库中已发布条目 | 否 | 否 | 否 | 否 | 否 | 否 |

- `Workspace.OwnerID` 继续表示唯一所有者；所有者不重复写入成员表。
- 私有知识库仅 owner 和有效成员可读；公开知识库中的已发布条目所有人可读。
- “公开知识库所有人可读”指公开库内已发布版本无需登录即可读取；草稿仍只对成员可见，避免编辑中的内容被提前公开。
- 登录用户读取公开内容时，如果同时具备编辑权限，详情响应返回 `can_edit=true` 并显示编辑按钮。
- service 层统一计算权限；前端只根据服务端返回的 capabilities 控制按钮，不把隐藏按钮当作安全边界。
- 非成员访问私有资源或无权限资源返回 404，不泄露资源是否存在；已识别成员执行越权写操作返回 403。

## 支持的文件类型

### MVP 类型矩阵
| 类别 | 扩展名示例 | 上传 | 站内预览 | 在线编辑 | 下载 | 预览策略 |
|------|------------|------|----------|----------|------|----------|
| Markdown | `.md`, `.markdown` | 是 | 是 | 是 | 是 | Vditor 静态渲染，服务端安全清洗 |
| 纯文本 | `.txt`, `.log`, `.ini`, `.conf`, `.env` | 是 | 是 | 是 | 是 | UTF-8 文本，等宽字体展示 |
| 结构化文本 | `.json`, `.jsonl`, `.yaml`, `.yml`, `.xml`, `.toml`, `.csv` | 是 | 是 | 是 | 是 | 文本预览、语法高亮；JSON/YAML 保存前校验语法 |
| 主流代码 | `.go`, `.js`, `.jsx`, `.ts`, `.tsx`, `.vue`, `.py`, `.java`, `.kt`, `.kts`, `.c`, `.h`, `.cpp`, `.hpp`, `.cs`, `.rs`, `.php`, `.rb`, `.swift`, `.scala`, `.sh`, `.bash`, `.zsh`, `.fish`, `.sql`, `.html`, `.css`, `.scss`, `.less`, `.dockerfile` | 是 | 是 | 是 | 是 | 纯文本读取与语法高亮，绝不执行 |
| PDF | `.pdf` | 是 | 是 | 否 | 是 | 原文件通过鉴权流式读取，在同源 PDF viewer 中展示 |
| 图片 | `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`, `.svg` | 是 | 是 | 否 | 是 | 光栅图片直接预览；SVG 清洗或栅格化后预览，不直接执行原始 SVG |
| Word | `.doc`, `.docx`, `.odt`, `.rtf` | 是 | 是 | 否 | 是 | 10 MiB 内 DOCX 可由受控浏览器组件预览；其他格式或超限文件后台转换为 PDF |
| Excel | `.xls`, `.xlsx`, `.ods` | 是 | 是 | 否 | 是 | 后台转换为 PDF；多工作表均需输出 |
| PowerPoint | `.ppt`, `.pptx`, `.odp` | 是 | 是 | 否 | 是 | 后台转换为 PDF |
| 压缩文件 | `.zip`, `.rar`, `.7z`, `.tar`, `.gz`, `.tgz` | 是 | 是 | 否 | 是 | 后台生成目录清单，只展示路径和解压后大小，不解压供下载 |
| 其他白名单附件 | 配置允许的其他扩展名 | 是 | 降级页 | 否 | 是 | 显示元数据、无预览提示和受控下载 |

### 文件限制
- 默认单文件上限为 100 MiB，允许通过配置调低；接口和反向代理限制必须保持一致。
- 可在线编辑的文本源文件上限为 2 MiB；超过上限仍可预览前 2 MiB并下载，但不可在线编辑。
- 文本统一按 UTF-8 保存；上传的 UTF-8 BOM 可移除。非 UTF-8 文本在转换任务中识别并转为 UTF-8，识别失败则预览失败但保留下载。
- 压缩包目录扫描必须限制最大条目数 10,000、目录层级 50、累计声明解压大小 2 GiB，并防止 Zip Slip 路径；超限时标记预览失败，不展开文件。
- 文件类型由服务端综合扩展名、文件签名和 MIME 白名单判定，不信任客户端 `Content-Type`。
- 宏文档可作为原件保存和转换，但绝不执行宏。转换进程必须禁用宏、网络访问并设置资源和时间限制。

## 核心流程

### A. 默认只读详情
1. 用户从工作空间目录树、文档列表或搜索结果点击任意条目。
2. 私有管理入口导航到 `/dashboard/docs/:id`；公开 Wiki 继续使用 `/wiki/docs/:id`，两者复用统一预览组件。
3. 后端根据工作空间公开状态、发布状态、登录用户身份和成员角色校验查看权限。
4. 详情接口返回条目元数据、文件类型、预览状态、版本号和 capabilities。
5. 前端按类型渲染 Markdown、文本/代码、PDF、图片、Office 转换 PDF、压缩包目录或降级下载页。
6. 有编辑权限且类型可编辑时显示“编辑”；不可在线编辑但可管理时显示“替换文件”。

### B. 文本类编辑与保存
1. 用户在详情页点击“编辑”，进入 `/dashboard/docs/:id/edit`。
2. 编辑接口返回源内容与 `revision`；Markdown 使用 Vditor，其他文本/代码使用代码编辑器或纯文本编辑器。
3. 自动保存和手动保存都携带最后读取的 `revision`。
4. service 在事务中按 `id + revision` 条件更新；成功后原子递增 `revision` 并返回新值。
5. 手动保存生成 `DocVersion`；自动保存不生成版本，修正当前实现与原知识库规格的偏差。
6. 保存成功后返回 `/dashboard/docs/:id` 并展示最新只读预览；如仅自动保存则保持编辑态。

### C. 文件上传与替换
1. editor/admin/owner 在当前目录点击“上传文件”，或新建菜单中选择“文档/文件”。
2. 前端以 multipart 提交工作空间、目录和原文件；进度条只表示原件上传进度，不代表预览转换完成。
3. 后端校验权限、大小、扩展名、签名和 MIME，生成服务端 object key，通过 Storage 接口保存原件。
4. 事务创建 `Attachment` 与 `Doc`。可直接预览的格式标记 `ready`；Office、压缩包及需规范化的格式创建 `DocumentConversionJob(pending)`。
5. 用户立即进入只读详情。转换中显示状态并允许下载原件。
6. scheduler 领取任务、读取原件、在隔离临时目录转换，将预览产物写回 Storage，并更新任务和附件状态。
7. 前端轮询详情状态；成功后自动展示预览，失败则展示失败原因摘要、重试按钮（管理员/编辑者）和原件下载。
8. 替换文件时创建新 Attachment，并在同一事务中将 Doc 指向新附件、递增 revision、保留旧附件用于版本/回滚；旧文件不能立即物理删除。

### D. Office 与压缩包预览
1. Office 文档使用无头 LibreOffice 转换为 PDF，不调用外部第三方在线预览服务。
2. PDF 产物通过受控预览接口返回，不能依赖公开静态 `/uploads` URL。
3. 压缩文件由隔离任务读取目录元数据，规范化路径后生成 JSON 清单；不把压缩内容展开到公开目录。
4. 详情页以树或表格展示压缩包内路径、文件/目录类型、压缩大小和声明解压大小。
5. 转换任务超时、崩溃或服务器重启后可基于数据库状态重试，不丢任务。

### E. 公开与分享访问
1. 公开知识库仍只展示已发布条目；草稿和转换中的新版本不对公众暴露。
2. 文件型条目发布时固定当前附件/预览版本；后续替换原件会将条目降为草稿，重新发布后公众才看到新版本。
3. token 分享继续允许分享私有或草稿条目，但响应使用收敛的只读 DTO，并通过同一受控预览/下载权限票据访问文件。
4. 公开访客永远看不到“编辑”“替换”“重试转换”等管理操作。

## 异常处理
| 场景 | 处理方式 |
|------|----------|
| 未登录访问公开库已发布条目 | 正常返回只读详情 |
| 未登录或非成员访问私有库 | 返回 404，不泄露资源存在性 |
| viewer 尝试编辑、替换或上传 | 返回 403，前端提示无操作权限 |
| 文档读取后被其他请求更新 | 保存条件不匹配，返回 HTTP 409 和当前 revision；保留本地内容，提供刷新、复制本地内容和查看最新版操作 |
| 自动保存遇到 409 | 立即暂停后续自动保存并提示冲突，不自动覆盖服务器内容 |
| 文件扩展名、签名与 MIME 不匹配 | 拒绝上传并删除已写入的临时文件 |
| 文件超过大小限制 | HTTP 413；前端在上传前同步提示 |
| Storage 写入成功但数据库事务失败 | 尝试删除孤立对象；删除失败写日志并由清理任务回收 |
| 数据库记录成功但 Storage 对象不存在 | 详情显示文件不可用；记录错误，允许有权限者重新上传 |
| Office 转换超时或失败 | 标记 failed，保留原件下载；允许手动重试且限制最大自动重试次数 |
| scheduler 转换中崩溃 | running 任务超过租约后回到 pending，按 attempts 重试 |
| PDF/图片预览失败 | 降级为文件信息和原件下载，不阻塞详情页其他操作 |
| 文本编码无法识别 | 标记预览失败，保留原件下载；不返回乱码内容 |
| JSON/YAML 保存语法错误 | 阻止保存并返回行列信息；普通代码和纯文本不做语义校验 |
| 压缩包加密、损坏、路径越界或超过扫描限制 | 标记目录预览失败，不解压；原件仍可按权限下载 |
| SVG 含脚本、事件或外链 | 清洗或栅格化；无法安全处理时只允许下载 |
| 删除或替换仍被版本引用的附件 | 仅解除当前引用，保留对象；无引用后由异步清理任务删除 |
| 用户重复提交上传 | 每次创建独立 Doc；客户端可用请求 id 防止网络重试造成重复创建 |

## 技术设计

### 总体结构
```text
HTTP 请求
  -> handler（绑定参数、取 user_id）
  -> service（权限、事务、文件策略、revision）
  -> models + Storage

cmd/scheduler
  -> DocumentConversionTask
  -> MySQL 领取持久化任务
  -> Storage.Open 原件
  -> Converter / ArchiveInspector
  -> Storage.Put 预览产物
```

### 数据模型

#### WorkspaceMember（新增表）
`internal/models/workspace_member.go`

| 字段 | 类型 / tag | 说明 |
|------|-------------|------|
| ID | uint primarykey | |
| WorkspaceID | uint `gorm:"uniqueIndex:idx_workspace_user;not null"` | 工作空间 |
| UserID | uint `gorm:"uniqueIndex:idx_workspace_user;not null"` | 成员 |
| Role | string `gorm:"size:20;index;not null"` | `admin/editor/viewer` |
| CreatedAt / UpdatedAt / DeletedAt | | 软删除 |

- owner 由 `Workspace.OwnerID` 表示，不写入本表。
- 新增组合唯一索引属于新增表，不影响存量数据。
- 被移除成员再次加入时恢复并更新原成员记录，不插入违反组合唯一索引的新记录。

#### Doc（追加字段）
`internal/models/doc.go`

| 字段 | 类型 / tag | 说明 |
|------|-------------|------|
| Kind | string `gorm:"size:20;index;default:'markdown';not null"` | `markdown/text/code/file` |
| Language | string `gorm:"size:50"` | 文本/代码高亮语言 |
| AttachmentID | *uint `gorm:"index"` | 文件型文档当前原件；Markdown 和站内新建文本可为空 |
| PublishedAttachmentID | *uint `gorm:"index"` | 文件型文档上次发布的原件；防止替换后的草稿文件提前公开 |
| PublishedRevision | uint64 `gorm:"default:0;not null"` | 上次发布的 revision，0 表示从未发布 |
| Revision | uint64 `gorm:"default:1;not null"` | 乐观锁版本，每次内容、标题或原件变化递增 |

- 存量文档通过默认值视为 `markdown`，继续使用 `Content` 和 `ContentHTML`。
- 文本/代码的规范化 UTF-8 内容以 `Doc.Content` 为当前真源；上传时的 Attachment 作为原始导入件保留。每次手动保存和自动保存都以更新后的 `Content` 动态导出下载内容，文件版本由 DocVersion 保存，不继续维护一份可漂移的“当前文本附件”。
- 文件型文档发布时将当前 `AttachmentID` 和 `Revision` 快照到 `PublishedAttachmentID`、`PublishedRevision`；替换原件后当前文档降为草稿，公开接口继续使用发布快照直到重新发布。
- `DocResponse` 增加 `kind`、`language`、`revision`、`attachment`、`preview` 与 `capabilities`。

#### Attachment（复用并追加字段）
`internal/models/attachment.go`

| 字段 | 类型 / tag | 说明 |
|------|-------------|------|
| ObjectKey | string `gorm:"size:500;index"` | Storage 内稳定键；替代把公开 URL 当真源 |
| Checksum | string `gorm:"size:64;index"` | SHA-256，用于完整性和去重判断 |
| PreviewStatus | string `gorm:"size:20;index;default:'none'"` | `none/pending/running/ready/failed` |
| PreviewObjectKey | string `gorm:"size:500"` | PDF、清洗图片或清单产物键 |
| PreviewMimeType | string `gorm:"size:100"` | 预览产物 MIME |
| PreviewError | string `gorm:"size:500"` | 对内错误摘要，不向公开访客暴露 |

- 保留现有字段兼容已有上传记录；知识库新文件必须使用 `ObjectKey`，不得写入公开静态目录后直接返回永久 URL。
- `URL` 仅可作为公开资源或兼容字段，不用于私有知识库鉴权。

#### DocumentConversionJob（新增表）
`internal/models/document_conversion_job.go`

| 字段 | 类型 / tag | 说明 |
|------|-------------|------|
| ID | uint primarykey | |
| AttachmentID | uint `gorm:"index;not null"` | 原件 |
| JobType | string `gorm:"size:30;index;not null"` | `office_to_pdf/archive_manifest/text_normalize/svg_sanitize` |
| Status | string `gorm:"size:20;index;not null"` | `pending/running/succeeded/failed` |
| Attempts | int `gorm:"default:0;not null"` | 已尝试次数 |
| MaxAttempts | int `gorm:"default:3;not null"` | 最大自动重试次数 |
| LeaseUntil | *time.Time `gorm:"index"` | 崩溃恢复租约 |
| OutputObjectKey | string `gorm:"size:500"` | 转换产物 |
| ErrorMessage | string `gorm:"size:1000"` | 内部诊断信息 |
| StartedAt / FinishedAt | *time.Time | |
| CreatedAt / UpdatedAt | | |

#### DocVersion（追加字段）
`internal/models/doc_version.go`

| 字段 | 类型 / tag | 说明 |
|------|-------------|------|
| Revision | uint64 `gorm:"default:1;not null"` | 快照对应 revision |
| Kind | string `gorm:"size:20;default:'markdown'"` | 快照类型 |
| Language | string `gorm:"size:50"` | 文本语言 |
| AttachmentID | *uint `gorm:"index"` | 文件版本原件 |

- 文件替换也生成版本，以支持查看旧版本和回滚。
- 删除 DocVersion 前必须重新计算附件引用；仍被当前 Doc、其他版本或发布快照引用时不得删除对象。

#### API 能力 DTO（不落表）
```json
{
  "can_view": true,
  "can_edit": true,
  "can_replace": true,
  "can_download": true,
  "can_publish": true,
  "can_share": true,
  "can_delete": true,
  "can_manage_members": false
}
```

### Storage 接口
在现有 `pkg/uploader` 基础上演进，不再让 handler 固定构造本地上传器：

```go
type Storage interface {
    Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
    Open(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    Stat(ctx context.Context, key string) (*ObjectInfo, error)
}
```

- 提供 `LocalStorage` 与 `TencentCOSStorage`，由配置工厂选择。
- 本地实现必须校验规范化路径始终位于 storage root 下。
- COS 实现启用前修复当前重复 `Object.Put`、固定 `image/png` 和超时上下文失效问题。
- 私有下载由 API 鉴权后流式返回；对象存储可由 service 返回短时签名 URL，但不得返回永久公开 URL。
- 配置新增项同步到 `config/config.example.yaml`；本地和 COS 的密钥只通过环境变量提供。

### 转换与隔离
- 新增 `internal/converter` 小接口，LibreOffice 实现使用 `exec.CommandContext`，禁止 shell 拼接。
- 每个任务使用独立临时目录和 LibreOffice profile，设置转换超时、CPU/内存/进程限制；容器内禁止出网并使用非 root 用户。
- `Dockerfile` 或独立 scheduler 镜像安装 LibreOffice、必要字体和压缩清单工具；HTTP 服务镜像无需承担转换依赖时可拆分构建 target，但仍属于同一仓库单体部署。
- 转换产物不是可信输入。HTML/SVG 必须清洗，PDF 通过同源 viewer 展示并设置安全响应头。
- scheduler 通过事务条件更新领取任务；`running` 且租约过期的任务可重试，确保进程重启不丢失。

### API 接口

#### 文档详情与编辑
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/docs/:id` | 可选登录 | 统一详情；公开已发布或有成员权限可读，返回预览与 capabilities |
| GET | `/api/docs/:id/edit` | editor+ | 取可编辑源内容与 revision；非文本类型返回 422 |
| PUT | `/api/docs/:id` | editor+ | 手动保存，body 必含 revision，成功建版本 |
| PUT | `/api/docs/:id/autosave` | editor+ | 自动保存，body 必含 revision，不建版本 |
| GET | `/api/docs/:id/preview` | 按文档查看权限 | 流式返回 PDF/图片/清单等预览产物，支持 Range |
| GET | `/api/docs/:id/download` | 按文档查看权限 | 下载原件，设置安全文件名和 `Content-Disposition: attachment` |
| POST | `/api/docs/:id/retry-preview` | editor+ | 失败后重建转换任务 |
| POST | `/api/docs/:id/file` | editor+ | 替换原件并生成版本，multipart 含 revision |

#### 文档创建与上传
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/docs` | editor+ | 新建 Markdown/文本类文档，增加 kind/language/revision |
| POST | `/api/workspaces/:id/files` | editor+ | 上传文件并创建 Doc；multipart 含 catalog_id、title、file、request_id |

#### 成员管理
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/workspaces/:id/members` | member | 成员列表 |
| POST | `/api/workspaces/:id/members` | owner/admin | 按用户添加成员并指定角色 |
| PUT | `/api/workspaces/:id/members/:userId` | owner/admin | 修改角色；admin 不可修改 owner |
| DELETE | `/api/workspaces/:id/members/:userId` | owner/admin | 移除成员；不可移除 owner |

#### 转换状态响应
```json
{
  "status": "pending",
  "kind": "office_pdf",
  "mime_type": "application/pdf",
  "error": "",
  "updated_at": "2026-08-13T12:00:00Z"
}
```

### HTTP 与安全要求
- `GET /api/docs/:id` 使用 `OptionalAuthMiddleware`，但 service 必须同时校验公开发布条件或成员关系。
- 文件响应设置 `X-Content-Type-Options: nosniff`、明确的 `Content-Type` 和内容安全策略。
- 下载文件名使用 RFC 5987 安全编码，移除控制字符和路径分隔符。
- PDF/图片预览支持 Range；对象存储签名 URL 的有效期不超过 5 分钟并绑定只读 GET。
- 文本预览返回 JSON 内容或 `text/plain`，禁止把上传 HTML 当作 `text/html` 返回。
- 上传、转换重试和下载沿用 Redis 限流模式，并记录用户、文档、附件、结果和耗时日志。

### 前端设计
- 新增 `/dashboard/docs/:id` 对应 `views/user/DocDetail.vue`；保留 `/dashboard/docs/:id/edit`。
- `WorkspaceDetail.vue` 的文档行、目录节点、键盘 Enter/Space 和搜索结果默认进入详情页；菜单保留明确“编辑”。
- 新建 `components/docs/DocumentViewer.vue` 按服务端类型与状态分派：
  - `MarkdownPreview.vue`
  - `TextCodePreview.vue`
  - `PdfPreview.vue`
  - `ImagePreview.vue`
  - `ArchiveManifest.vue`
  - `FileFallback.vue`
- Markdown 只读渲染复用 `web/admin/src/views/admin/ArticleView.vue` 的 `Vditor.preview()` 模式，并统一 `WikiDoc.vue`、`ShareDoc.vue` 与版本详情的渲染组件。
- 编辑页按 `kind` 分派 Vditor 或 CodeMirror 6；不能编辑的类型不初始化编辑器。CodeMirror 6 按需加载语言扩展，并复用于开发类文件的只读语法高亮。
- JSON/YAML 提供受限的只读树与源码切换，CSV 提供最多 500 行、100 列的表格与源码切换；DOCX 在 10 MiB 内使用按需加载的 `docx-preview`，禁用 AltChunk HTML 和非安全链接协议。
- 新建菜单提供“新建文本文件”和“上传文件”：直接新建时按文件名白名单推断 Markdown、纯文本、代码或结构化数据的 `kind/language`，例如 `README.md`、`styles.css`、`worker.rs`、`config.json`；二进制类型必须上传真实原件。上传 `.md` 后解析为 UTF-8 `Doc.Content`，与直接新建归一到同一预览、编辑和版本流程。
- 详情页固定显示返回、标题、类型、大小、更新时间、发布状态和权限允许的操作。转换状态在正文区域展示，不用全屏阻塞。
- 移动端详情优先保证正文阅读、下载和编辑入口；目录侧栏折叠为抽屉，PDF/代码预览可横向滚动。

### 部署决策
- 默认限制确定为：单文件 100 MiB、文本在线编辑 2 MiB、压缩包累计声明解压大小 2 GiB；配置可调低，MVP 不允许配置超过这些安全上限。
- LibreOffice 和压缩清单工具只安装到独立 `scheduler-runtime` 构建 target。现有 `server-runtime`、`admin-runtime` 不安装转换依赖，避免扩大 HTTP 服务镜像和攻击面。
- scheduler 仍使用现有 `cmd/scheduler` 二进制和共享 internal 代码，不新增服务协议；本地存储部署时继续挂载共享文件卷，对象存储部署时通过 Storage 读取原件和写入产物。
- 成员管理首期只支持按用户名精确查找并添加已注册站内用户，不发送邮件、不生成邀请链接，也不自动注册账号。

### 实现步骤（每步可独立 commit）
1. [x] **权限模型**：新增 `WorkspaceMember`、角色常量、权限 service 与成员 API；将现有 owner 查询逐步改为统一授权函数。
2. [x] **详情与并发模型**：为 Doc 增加 `Kind/Language/AttachmentID/Revision`，为 DocVersion 增加文件快照字段，实现 `GET /api/docs/:id` 和 409 乐观锁。
3. [x] **只读详情前端**：新增详情路由与 `DocDetail.vue`，把列表/树默认点击改为查看，增加显式编辑按钮并修正面包屑。
4. [x] **统一 Markdown 预览**：抽取只读组件，复用到私有详情、公开 Wiki、分享页和版本历史。
5. [x] **Storage 接口**：演进 `pkg/uploader` 为可读写删除存储，实现本地/COS provider，修复 COS，补全配置和私有流式下载。
6. [x] **附件与上传 service**：接入 `Attachment` 表，增加文件策略、签名/MIME/大小校验和工作空间文件上传 API。
7. [x] **文本与代码**：实现文本规范化、语言识别、预览和在线编辑；JSON/YAML 保存校验；与 revision/版本历史接通。
8. [ ] **直接预览类型**：实现 PDF Range 预览、光栅图片预览和 SVG 安全处理。（PDF 与光栅图片鉴权预览已完成，Range 与 SVG 安全产物待完成）
9. [ ] **持久化转换任务**：新增 `DocumentConversionJob`、converter 接口和 scheduler 任务，实现租约、重试与状态查询。
10. [ ] **Office 预览**：部署 LibreOffice 转 PDF，覆盖 Word/Excel/PowerPoint/OpenDocument，前端轮询并展示 PDF。
11. [ ] **压缩包清单**：实现 zip/rar/7z/tar/gz 安全目录扫描与限制，前端树/表格展示。
12. [ ] **文件替换与生命周期**：文件版本、回滚、引用计数、孤立对象清理和转换产物清理。
13. [ ] **公开/分享适配**：文件型发布快照、公开预览、收敛分享 DTO 和受控文件访问。
14. [ ] **联调与加固**：限流、安全响应头、失败恢复、桌面/移动端体验、迁移验证及部署文档。

### 参考的现有模式
- `docs/specs/workspace-knowledge-base.md` — 当前知识库业务、分享、目录和版本基线。
- `web/blog/src/router/index.js` — 当前只有 `/dashboard/docs/:id/edit`，需增加详情路由。
- `web/blog/src/views/user/WorkspaceDetail.vue` — 当前文档行和树节点直接调用编辑路由；复用目录树组装。
- `web/blog/src/views/user/DocEdit.vue` — 当前编辑、自动保存、版本、分享和发布流程。
- `web/blog/src/views/wiki/WikiDoc.vue` — 已有公开只读文档样式。
- `web/blog/src/views/ShareDoc.vue` — 已有 token 只读状态页。
- `web/admin/src/views/admin/ArticleView.vue` — Vditor 静态 Markdown 预览范式。
- `internal/models/doc.go` — 当前 Markdown Doc 及 DTO。
- `internal/models/attachment.go` — 已有但未接入业务的文件元数据模型。
- `internal/service/doc_service.go` — owner 鉴权、事务行锁、版本、发布和 Markdown 渲染。
- `internal/service/public_wiki_service.go` — 公开条件与 Bluemonday HTML 清洗。
- `internal/handler/upload_handler.go` — 现有 multipart、限额和图片上传模式；业务应下沉到 service。
- `pkg/uploader/uploader.go` — 本地/COS uploader 抽象起点。
- `pkg/uploader/local.go` — 本地文件写入实现。
- `pkg/uploader/tencent_cos.go` — COS 实现，接入前需修复重复上传和 MIME 问题。
- `internal/scheduler/scheduler.go`、`cmd/scheduler/main.go` — 现有持久化转换任务的运行外壳。
- `internal/router/blog.go` — protected、public、optional auth 路由及当前公开 `/uploads` 静态目录。

## 测试计划

### 后端单元与集成
- [ ] 权限矩阵覆盖 owner/admin/editor/viewer/public visitor 对私有、公开、草稿和已发布文档的读写结果。
- [ ] 非成员访问私有文档返回 404；成员越权写返回 403。
- [ ] 两个客户端基于同一 revision 保存，只有一个成功，另一个稳定返回 409 且内容未覆盖。
- [ ] autosave 不创建 DocVersion；手动保存、文件替换和回滚创建正确快照。
- [ ] 存量 Doc 未显式设置 Kind 时按 Markdown 正常读取。
- [ ] 扩展名、客户端 MIME 和文件签名不一致时拒绝上传。
- [ ] 文本 UTF-8、BOM、非 UTF-8 转换、超 2 MiB和二进制伪装场景。
- [ ] JSON/YAML 有效与无效保存校验；HTML/JS 只作为转义文本返回。
- [ ] 本地 Storage 路径越界防护、读写删除和对象不存在；COS 使用 mock 验证单次 Put 与正确 MIME。
- [ ] 私有文件无法通过 `/uploads` 静态 URL 绕过权限。
- [ ] PDF/图片预览与下载支持 Range、正确 MIME、`nosniff` 和 Content-Disposition。
- [ ] Office 转换成功、超时、失败、重试、租约过期和 scheduler 重启恢复。
- [ ] 压缩包 Zip Slip、加密包、损坏包、超条目数、超层级和解压炸弹声明值。
- [ ] SVG 脚本、事件属性、外部 URL 被移除或降级下载。
- [ ] Storage 成功/数据库失败与数据库成功/Storage 缺失的补偿路径。
- [ ] 文件替换、版本回滚和删除不会误删仍被引用的原件/预览产物。
- [ ] 公开 Wiki 只返回已发布附件版本；替换后降草稿不会泄露新文件。
- [ ] token 分享 DTO 不暴露 object key、owner_id、内部转换错误或永久存储 URL。

### 前端
- [ ] 工作空间文档行、树节点、搜索结果和键盘操作均进入只读详情。
- [ ] 仅有 `can_edit` 时显示编辑按钮；viewer 和公开访客无编辑操作。
- [ ] 保存成功回到详情；409 时保留本地内容并暂停 autosave。
- [ ] Markdown、文本、代码、JSON/YAML、PDF、图片、Office 转换中/成功/失败、压缩清单和未知类型均有正确视图。
- [ ] 上传进度与转换进度分开显示，转换期间可下载原件。
- [ ] 桌面和移动端可阅读、下载、进入编辑并返回原位置。

### 验证命令
```bash
go test ./... -count=1
go vet ./...
go build ./...

cd web/blog
pnpm lint
pnpm build
```

## 验收标准
- 从私有知识库目录点击已有 Markdown 文档时，首先看到只读渲染页，不初始化编辑器。
- 有编辑权限的用户可通过明确“编辑”进入编辑态，保存后回到最新详情。
- viewer 和公开访客无法通过 UI 或直接 API 修改内容。
- Markdown、文本、主流代码、JSON、YAML、PDF、图片、Office 和压缩包均可上传并获得对应预览或明确转换状态。
- Office 文件转换为 PDF 后可站内阅读；转换失败时详情仍可打开并可下载原件。
- 压缩包只展示安全目录清单，不执行或公开解压内容。
- 私有原件和预览产物不能通过永久公开 URL 绕过鉴权。
- 并发陈旧保存返回 409，不再发生静默后写覆盖。
- 本地存储和腾讯 COS 均通过同一 Storage 接口工作，业务层不判断 provider。
- 转换任务在 scheduler 重启后可恢复，不依赖内存队列。

## 已确认决策
- 文本与代码编辑器使用 CodeMirror 6，按文件语言动态加载扩展，不采用体积更大的 Monaco。
- 使用 100 MiB 单文件、2 MiB 在线编辑和 2 GiB 压缩包声明解压大小上限。
- LibreOffice 及压缩文件工具进入独立 `scheduler-runtime` target，不进入 HTTP 服务镜像。
- 同时提供新建 Markdown 和上传 `.md` 两个入口，数据统一归一为 Markdown Doc。
- 成员添加仅面向已有站内用户，使用用户名精确查找。
- 公开文件的搜索引擎索引行为沿用当前公开 Wiki 策略，不新增独立开关。
- 浏览器端解析器只用于可安全收敛的格式：DOCX 使用 `docx-preview`；不采用已弃用的 `vue-office`、存在未修复漏洞的 npm `xlsx` 或长期失维护的 PPTX 渲染器。旧 Office、XLSX、PPTX 继续以隔离 LibreOffice 转 PDF 为正式预览路径。

## MVP 范围
本期完整纳入：

1. 统一只读详情、显式编辑和保存后返回预览。
2. owner/admin/editor/viewer 成员角色，以及公开库所有人可读、私有库成员可读。
3. revision 乐观锁、409 冲突处理和 autosave 不建版本。
4. Markdown、纯文本、主流代码、JSON、YAML、XML、TOML、CSV 在线预览与文本类编辑。
5. PDF、图片直接安全预览。
6. Word、Excel、PowerPoint、OpenDocument 后台转 PDF 预览。
7. zip、rar、7z、tar、gz 安全目录清单预览。
8. 未知白名单文件的元数据页与受控下载降级。
9. 可插拔 Local/COS Storage、私有流式访问和文件生命周期管理。
10. MySQL 持久化转换任务、scheduler 轮询、租约、失败重试与恢复。
11. 文件型文档的版本、替换、公开发布和 token 分享适配。

后续迭代：实时协作、Office 原生在线编辑、压缩包内文件递归预览、更多专业格式、病毒扫描服务、全文内容索引和系统管理后台。
