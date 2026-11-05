# 文章 SEO：迁移、部署、验收与回滚

范围：稳定 `/blog/<slug>`、数字旧链接 301、初始 HTML head、公开 slug API、
RSS/Atom、sitemap（含 `/sitemaps/<n>.xml`）、robots。正文仍由 Vue 渲染，不是 SSR。
实现验证记录见 [SEO_VALIDATION.md](SEO_VALIDATION.md)；测试通过不代表已经部署。

## 1. 唯一规范站点地址（所有运行入口）

`PUBLIC_SITE_URL` 是**博客对外 origin**，不是容器地址、API 内网地址或后台域名。
server、admin、scheduler 均在数据库连接前校验，debug/test/release 无例外。
只接受带 host 的 `http`/`https` origin；禁止凭证、非根路径、query、fragment
（包括空 `?` / `#`）。尾斜杠被去掉。没有 Host、X-Forwarded-Host、
`window.location.origin` 或编译期地址兜底。缺失/非法配置拒绝启动。

优先级：进程环境 > `.env` > `config/config.yaml` 或 `config/admin.yaml`。
admin 文件缺失时由共享 loader 回退 config 文件；**非法 admin 配置不会再重试绕过校验**。
迁移专用 `cmd/article-slugs` 不加载这些业务配置，只接受显式离线迁移 DSN。

| 场景 | 显式配置 | 说明 |
| --- | --- | --- |
| Go 内嵌页面本地开发 | `http://localhost:8081` | 两个 YAML 示例和 env.example 的开发默认 |
| Vite 浏览器开发 | `http://localhost:3001` | **在 Go 进程设置**；Vite 代理后端仍是 :8081 |
| Compose / HTTPS | 如 `https://blog.example.org` | 改 `.env`，两个 Compose 的三个 blog 副本 + admin + scheduler 全部转发 |
| K8s development | `http://dev.blog.test.com` | 与现有 development ingress 一致；若启用 TLS，三个 overlay 一起改 https |
| K8s production | 环境拥有的 ConfigMap | base 仅引用必需 key，不内置开发 origin |

不要把示例 localhost 或 `dev.blog.test.com` 用于生产。规范 origin 不能从请求头推导，
但不能自动判断操作者的“合法但选错的域名”；发布前必须核对 TLS/DNS/重定向。
无需新增 Docker `ENV` 或前端 `VITE_PUBLIC_SITE_URL`；同一镜像可用于不同环境。
`/api/settings/public` 返回校验过的 `data.public_site_url`，覆盖同名数据库设置；
公开 slug API 的 `data.seo` 是 SPA metadata 的权威投影，异步站点设置不会覆盖文章 head。

### Compose

两个基础文件均要求非空 `PUBLIC_SITE_URL`，否则 `docker compose config/up` 即报错，
并将值传入 backend-1/2/3、admin-backend、scheduler，不能只靠宿主机 `.env` 隐式注入。
HTTPS overlay 继承基础文件的应用环境，不需要复制配置。以下只渲染配置，不部署：

```sh
PUBLIC_SITE_URL=https://blog.example.org docker compose -f docker-compose.yml config --quiet
PUBLIC_SITE_URL=https://blog.example.org docker compose -f docker-compose.external-db.yml config --quiet
PUBLIC_SITE_URL=https://blog.example.org docker compose -f docker-compose.yml -f docker-compose.https.yml config --quiet
```

生产环境不能直接照搬示例数据库密码；连接与 secret 仍按原部署文档提供。
不要把渲染后的全部配置粘贴到日志/工单（可能包含 secret）。

### Kubernetes

base 的三个 Deployment 各引用必需 ConfigMap key `PUBLIC_SITE_URL`：
`blog-api-public-site`、`admin-api-public-site`、`scheduler-task-public-site`。
未提供时 Pod 无法建立容器配置；空/非法值由应用拒绝。
development overlay 各自生成带 hash 的 ConfigMap，Kustomize 自动重写引用，配置变更会更新 Pod template。
生产 overlay 应各生成同名逻辑 ConfigMap、写入**同一个生产博客 origin**，不能复用 development overlay。

例如在生产 blog overlay 的 `kustomization.yaml` 使用（admin/scheduler 换对应名称）：

```yaml
namespace: inkspace-backend
resources:
  - ../../base
configMapGenerator:
  - name: blog-api-public-site
    literals:
      - PUBLIC_SITE_URL=https://blog.example.org
```

若直接部署 base，先由环境配置系统在**相同 namespace**创建这三个 ConfigMap；
更新无 hash 的 ConfigMap 后必须显式 rollout restart，不能期待 env 自动热更新。
镜像、数据库/Redis secret、上传存储等现有基础设施仍需另行配置；base 不是完整生产部署。
可离线运行 `kubectl kustomize deploy/apps/backend/<app>/overlays/development` 审查结果；
不要误执行 apply 到现有集群。Ingress/nginx 的 `/` 前缀已经转发全部根发现端点和分片。
不增加 SPA rewrite、隐藏 5xx 的 fallback 或跨认证公共缓存。

## 2. 发布顺序：迁移/回填先于应用（已有数据库）

1. 在**隔离 MySQL**演练备份恢复、迁移和下述验收，记录版本与耗时。禁止使用开发库测试。
2. 为应用、数据、上传文件做可恢复备份。记录文章总数/公开数，保存维护者受控的
   slug 映射审计（不要公开私密文章标题/slug）。准备维护窗口与旧镜像。
3. 停止**全部写入者**：旧/新 blog 副本、admin、scheduler、知识库发布入口、后台任务、
   ad-hoc SQL 写入。网关可返回维护页/503，不能让旧进程继续创建无 slug 文章。
   MySQL DDL 隐式提交，advisory lock 不能替代停止旧写入者。
4. 审查分阶段、可重复执行的 [schema.sql](../scripts/article-slugs/schema.sql)、
   [finalize.sql](../scripts/article-slugs/finalize.sql)（只在修复后验证/建约束，completion marker 最后安装）、
   [迁移说明](../scripts/article-slugs/README.md)、`internal/database/article_slugs.go`。
   SQL 文件必须按迁移说明**分阶段执行**，不能拼接运行；每个 SQL 会话须显式确认离线，
   确定性拼音/冲突处理和完整日期修复仍由中间的 Go 命令完成。推荐用 Go 命令编排完整流程。
5. 从受保护的 secret 注入 `INKSPACE_ARTICLE_SLUG_MIGRATION_DSN`（必须显式选择目标库），
   设置 `INKSPACE_ARTICLE_SLUG_MIGRATION_OFFLINE=1`，在可信构建环境运行：

   ```sh
   go run ./cmd/article-slugs
   ```

   不要在命令行/工单打印 DSN，不要使用 `set -x`。命令不读取开发配置兜底。
   **即使 ticket01 的 slug 迁移已经执行，也必须重跑 ticket02 扩展后的命令**，修复
   NULL/零/部分零/不可能的历史 `publish_at` 并安装 completion marker。
   所有草稿/私密/软删除记录一起处理；有效历史日期和已有非空 slug 保留。
6. 核对 slug `scanned/filled` 与 publication `scanned/filled/preserved`，检查错误/异常 ID。
   非法/重复已有 slug、无法回退的 created_at 必须停止并人工确认，不删记录，不随意改已公开 slug。
   重跑命令，预期两种 `filled=0`。中断后保持离线并重跑，不手工伪造完成状态。
7. 用受保护 MySQL client defaults file 执行只读
   [verify.sql](../scripts/article-slugs/verify.sql)，所有 missing/invalid/duplicate 计数为 0，
   `ux_articles_slug` 必须是全局 unique（包含删除行保留），
   `idx_articles_first_publication_v1(status,deleted_at,publish_at DESC,id DESC)` 必须存在。
   该非 unique index 是**已完成 publication 修复的启动标记**，不能提前手建。
   SQL 的简易零日期计数不能替代迁移命令的完整文本日历校验。
8. 才允许部署依赖 slug 的应用。启动只检查索引元数据，不批量回填；所有副本共享同一数据库，
   不在每个启动容器中运行迁移。新空库由初始化建立就绪 schema；有历史文章的库不自动升级。
9. 核对所有入口同一个 `PUBLIC_SITE_URL`、数据库、镜像版本；先 canary，再恢复全部流量/写入。
   TLS 启用后使用 https origin；重建博客前端产物再构建 Go，以免旧 dist 缺少 head markers。

## 3. HTTP / SPA 上线验收

先在隔离/预发布环境创建维护者自有的公开和草稿/私密/删除 fixture；**不要猜测真实私密 slug**。
CI 自动测试中，最高 seam 是真实 Gin 路由；可注入公开 reader 验证协议，实际 SQL 过滤由
显式 `INKSPACE_SEO_TEST_MYSQL_DSN` 门禁 fixture 验证。没有 DSN 的 SKIP 不能算持久化验收。

只读 curl smoke（Python 3 标准库编排 curl，默认不输出正文、URL 或鉴权信息）：

```sh
python3 scripts/seo-smoke.py --origin https://blog.example.org \
  --slug public-fixture --id 123 --hidden-slug unpublished-fixture
# 直连 canary Go 服务，验证 Host 独立性；不要对基于域名路由的 ingress 发送错误 Host：
python3 scripts/seo-smoke.py --origin https://blog.example.org \
  --target http://127.0.0.1:8081 --slug public-fixture --id 123 \
  --hidden-slug unpublished-fixture --host-check
```

脚本禁用用户 curlrc、不发 cookie/token、不跟随重定向、不写生产数据，只把响应临时保存在自动清理的私有目录。
应使用本来就可公开的 fixture；私有 fixture slug 仍可能进入服务访问日志，需遵循日志权限策略。
脚本验证指定隐藏 fixture，不声称扫描或证明所有历史隐私内容；完整过滤依赖 SQL 门禁测试。
高频发文时 GET/HEAD 间 ETag 合法变化可能导致 smoke 失败，应在安静窗口重跑，不放宽契约。

人工/自动清单：

- [ ] slug 页面 200，数字 URL GET/HEAD 301，Location 指向绝对 canonical；大小写/尾斜杠归一。
- [ ] 草稿、私密、软删除和不存在页面均 404；作者/admin 登录也不改变公开 HTML/slug API 语义。
      `status/show_all/author_id` 不可绕过公开集合。
- [ ] **直接请求构建后的 Go 服务初始 HTML**，解析 title/description/canonical/OG/Twitter/
      BlogPosting JSON-LD；中文、引号、`<`、`&` 不变成代码；不含邮箱/权限、私密设置或正文。
      非法/私有/签名图片 URL 被省略；无封面无摘要仍合法；RSS/Atom alternate 指向同一 origin。
- [ ] `/api/settings/public` 的 `data.public_site_url` 与 `/api/articles/slug/<slug>` 的
      `data.seo.canonical` 一致，不能从数据库同名设置或 Host 改写。
- [ ] RSS/Atom XML 可解析，最近最多 50 篇首次发布时间排序、ID tie-break、仅纯文本摘要；
      sitemap 全部公开页面，无 drafts/private/deleted，lastmod 不因互动计数改变。
- [ ] `/feed.xml`、`/atom.xml`、`/sitemap.xml`、`/robots.txt` MIME/UTF-8 正确，
      `/sitemaps/1.xml` 可请求；根超过 50,000 URL 时为 index，跟随每个 shard 做完整并集校验。
- [ ] 发现 GET/HEAD 头一致、HEAD 无正文；ETag 条件请求 304 无正文；成功 public max-age=300；
      HTML no-cache。删除/撤下文章后新源站请求立即移除（公共代理最多有 5 分钟陈旧窗口）。
- [ ] 仅在隔离测试库演练 DB 故障：文章/feeds/sitemap 500 + no-store，不返回空 XML 或 304，
      不泄露 SQL；robots 无 DB 依赖仍可用。**不要在生产删表做故障测试**。
- [ ] Vite :3001 代理四端点、`/sitemaps/` 和 `/api`；SPA A→B metadata 更新/去重，
      404/离开后清理，晚到设置不覆盖文章。Vite 自己返回的原始 HTML **不是** Go 初始 metadata 验收。
- [ ] `/uploads/knowledge` 继续鉴权/阻断；robots 不替代权限控制。admin SPA 不受博客 metadata 注入。

## 4. 运维监控与容量边界

上线前建立基线，canary 期间/上线后观察：数字 URL 301/404，slug URL 200/404/5xx，
各发现端点（含 shard）4xx/5xx、p50/p95/p99 延迟、响应大小、XML 解析/生成失败、
slug 冲突/迁移拒绝、数据库查询延迟/连接池和进程内存。
按路由模板聚合，避免将敏感 slug、query、作者邮箱或原始 SQL 放进指标标签/外部告警。
既有 Gin 访问日志可能带 URI/query，按生产日志策略限制权限、保留期和脱敏，不扩展成公开日志面板。

ETag 304 比率可作源站条件命中指标；代理 HIT/MISS/Age 需从网关/CDN单独采集。
当前实现是**公共缓存头 + 每次新查询/渲染 + content-derived ETag**，不是应用内结果缓存，
304 也会查询 DB，不能把它当成减少 DB 查询的缓存命中。先测负载，不凭空承诺性能。
告警示例（按本部署流量调整）：canary 新增持续 5xx、源站发现延迟显著高于基线、
XML smoke 失败、发布后 metadata origin 不一致，立即暂停扩容。

sitemap 每请求读取所有公开文章的紧凑 ID/slug/date 投影，内存 O(公开篇数)，不是全文。
分片按稳定 ID 顺序，但**不同 HTTP 请求不是同一持久快照**；发布/删除可能移动 offset 归属。
高写入期遍历可重试；若站点规模超出容量，应另立持久快照/流式生成设计，不在本次冒充已解决。

## 5. 回滚

1. 停止扩大流量，保持维护窗口，保存脱敏状态码/版本/迁移计数；先修复错误配置或回滚 canary。
2. **不得删除 slug/publish_at 列、unique index 或 publication completion marker**。
   不修改已有 slug，不释放软删除地址。回滚应用不会自动回滚 MySQL DDL。
3. 只有已验证仍维护 slug 和首次发布时间契约的上一应用版本可恢复写入。
   SEO 之前的老版本可能创建 NULL slug 或重置发布时间，不能直接恢复它的写流量；
   必要时回滚到只读维护页，修复后重新离线迁移/验证。
4. origin 误配时所有 server/admin/scheduler 同步恢复正确值并重启；清除网关/CDN
   发现文档缓存，考虑 301/搜索引擎缓存无法立即撤回，不临时重用旧地址。
5. 数据恢复只用批准的备份恢复流程，评估备份之后已对外发布的永久链接保留；
   不把“删索引让应用启动”作为回滚手段。完成同一份 HTTP/SPA smoke 后再恢复写入。
