# 来源与设计方法提炼

## 阅读边界

来源为用户提供的三份 RFC 文本。本次没有对应代码仓库、实验报告或 ADR 正文，因此以下引用代表文档中的主张与证据入口，未独立核实实现。

| 标识 | 文档 | 在这组设计中的作用 |
|---|---|---|
| A | RFC 20260926 Dynamic Routing in Traffic Segment | 跨 dns-server、cert-mgr、ingest 与 proxy 的完整动态配置链路 |
| B | Proxy Object Store Reader | 运行时读取、TLS 集成、缓存、WAL、生命周期与运维设计；文本标注 Implemented，并含后续实现修订 |
| C | Object Store backend for certificate and merchant data | 共享存储协议、发布容器、证书工作流与 ingest 契约；reader 为其 step 4 |

## 八项能力及原文依据

| 能力 | 原文入口 | 可复用的方法 |
|---|---|---|
| 事实取证 | C：Facts established from the codebase；B：References and established facts、Implementation notes | 事实表将路径/API/依赖版本与设计后果连接；区分原始提案、实验和明确修订 |
| 范围与取舍 | B：Scope and compatibility、Alternatives considered；C：Scope, non-goals and known breakage | 说清放弃的能力、已知破坏和未提供的恢复机制，不用新功能掩盖边界 |
| 跨层契约 | C：Key Encodings、Value encoding、Entry format；A：Storage record、DNS Server APIs | 给边界所有者、精确格式、校验层次和错误语义，使用共享模型减少解释分歧 |
| 一致性推理 | A：Revision and publication、Snapshot and routing publication、Type 6 ingest；C：Version stamps、Ingest contract | 写明版本与数据的绑定、每键排序、原子性范围，以及重试/乱序后仍成立的不变量 |
| 故障恢复 | A：Trigger failure and recovery、Atomic publication and stale-result protection、Failure fallback and recovery；B：Startup, outages and recovery | 用状态转换覆盖部分成功、陈旧异步结果、恢复无新事件和故障中重启 |
| 预算与运维 | B：Cold-read budget、Resource ownership and lifecycle、Configurable load gating、Health and observability | 推导超时预算、并发和容量；将实验阈值、初始默认与生产保证分开 |
| 实施拆分 | C：Implementation approach steps 1–4；B：Implementation boundaries；A：Workflow and worker registration | 契约先行，明确组件依赖、注册/配置/资源注入和阶段间破坏 |
| 行为验证 | A：各节 Verification；B：Verification；C：golden tests、Testing | 从失败窗口与不变量生成可观察验收场景，验证真实系统行为 |

## 案例 1：确认成功不等于配置生效

A 中的边界依次为：

```text
OLTP 提交
  → Temporal 确认工作流启动
  → buffer 对象上传且位置持久进入 manifest
  → ingest 持久写入记录与 Applied Version
  → reader 观察到更新
  → proxy 原子激活路由
```

DNS 返回成功仅确认前两步。工作流成功仅证明 buffer 已发布有效快照，可能是合并后的更高 revision。每层 ACK 都只证明本层承诺。

若提交成功但启动响应丢失，A 返回带 `committed: true` 和已提交 revision 的部分成功错误；这不证明工作流未启动。其恢复是重新触发发布当前 revision，不是重放旧 routing PUT。后者会增加 revision，并可能覆盖并发更新。没有 outbox 的提交到触发空窗是接受的限制。

**提炼**：先列阶段，再给每个返回值、错误和恢复操作限定意义。

## 案例 2：不同数据不一定共享同一种版本规则

A 的 routing 要求一次一致读取得到 routing 与 routing_version。触发版本为 12、当前快照为 14 时，可发布 14 的内容并标记 14，不能标记 12。同一个 routing revision 对应同一配置。

C 的证书 buffer 活动则从前一个 save 活动接收版本，在执行时重读数据库；其过滤器重试可以在相同版本读到更大快照，ingest 对相同版本采用 first-accepted-wins。该例外依赖明确的过滤器不变量。

C 还承认：晚读的 `F@1={a,b,c}` 可能被先捕获的 `F@2={a,b}` 替代，若 c 在版本 3 才加入。存储版本单调不等于过滤器集合单调，恢复依赖之后包含 c 的成功发布。

**提炼**：先决定每种版本承诺什么，再证明规则；不能把 routing 的“同版本同内容”错误推广到所有数据，也不能把过滤器例外推广到 routing。

## 案例 3：有序发布仍然可能看到缺失依赖

C 顺序发布 chain → certificate → filter，每条都等待 durable。并行 ingest 与晚读 filter 仍可能令读者先看到 filter 或 certificate。

读者契约因此规定 filter 命中但记录缺失要按 miss 处理；`chain_hash: Some(hash)` 但链缺失时证书不可用，不能当作无链证书；之后允许再次查找。链每次签发/续期都重发，依靠内容寻址与版本规则消重，避免“仅首次创建时发布”形成永久丢链。

**提炼**：将发布顺序、持久化顺序、应用顺序和可见性分别讨论；依赖补发和缺失处理都是契约。

## 案例 4：缺失、空值、损坏可能呈现相同默认行为，但语义不同

A 中 SQL `NULL` 和恰好 `{}` 由 publisher 归一化为 `{"routing":[]}`；API reset 只接受显式空数组；proxy 读到 `{}` 属于错误，不能当作合法 reset。

proxy 对 absent、empty 和 failure 都可能使用默认路由，但原因、degraded 状态与恢复动作不同。合法非空 custom routing 整体替换默认规则，未匹配路径按其明确的服务兜底，不与默认特殊路由合并。

**提炼**：归一化只发生在指定边界。区分“对外结果相似”与“内部语义相同”，并据此设计指标与测试。

## 案例 5：缓存支持持续服务不代表支持离线启动

B 放弃 proxy 自建完整 snapshot，使用 SlateDB 自有缓存。缓存不持久保存 manifest；reader open 必须访问对象存储。因此已打开进程可在缓存覆盖范围内维持服务，故障中重启则需等存储恢复才 ready。

缓存还受容量淘汰和 best-effort 写入影响，预热完成也不是永久完整副本保证。ready 不等于 WAL 已追平或 warming complete。实验的本地延迟不构成生产新鲜度承诺。

**提炼**：逐阶段描述能力与证据边界，给复杂度减少写清对应的能力损失。

## 案例 6：同一有效序号也可能需要重新激活

A 指出 active custom 因读失败切 defaults 后，恢复时可能读到以前验证过的相同 row sequence。若仅凭“序号未变”提前返回，路由会永久停在 failure defaults。

此外，旧 read error 不可覆盖新 custom；旧 success 也不可覆盖较新的 failure transition。absent/error 不带 row seq，需其他 generation 或串行化机制保护激活，并在过期结果被丢弃后保留重试义务。

**提炼**：分开记录观察历史、编译缓存和当前 serving 状态；测试状态转换，不能只测试数值递增。

## 不能固化为通用规则的原文选择

- Rust、SlateDB 0.17.0、matchit 0.9.1、Type 6、七字节键是这组协议的约束。
- 一条 Record Change 对应一次 produce 是这组 buffer 契约；其他设计可选择批次原子操作，但需定义边界。
- 故障时路由切默认、证书某些 resident replacement 保留旧值，是不同领域/类型的策略，不能统一套用。
- 不做迁移、无 outbox、无自动 reconciliation 是已接受的局部范围；已有部署不能默认照搬。
- 默认 5 秒、10 次尝试、900 秒等数字需要在新系统重新推导。不要沿用原文“参数看似表示什么”而忽略 SDK 的实际有效设置。
- `Simplicity First` 与 `Surgical Changes` 应落实为范围与取舍，不能用固定代码行数替代工程判断。

## 明确的修订关系

B 的 Implementation notes 声明在差异处实现为准，并明确替代下文部分描述：forked rustls 已移除；SNI Resolution/ReqInfo 已携带 Traffic Segment；warm_sst 失败有 timer 重试；调用者提供的 split cache 需单独关闭；Type 1/2 损坏替换保留之前解码值。

使用这些文本设计新方案时按修订关系理解，不能复制旧段落中的“从 SNI 推导 segment”或“必须保留 rustls fork”。若要判断当前仓库的事实，还需查验实际代码。
