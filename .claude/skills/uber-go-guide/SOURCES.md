# 来源与适配说明

## 来源记录

- 中文来源：[xxjwxc/uber_go_guide_cn](https://github.com/xxjwxc/uber_go_guide_cn)。
- 本次整理基准：[README.md @ 0b729562805e94aa4d7f5e2b37991156e0adfc18](https://github.com/xxjwxc/uber_go_guide_cn/blob/0b729562805e94aa4d7f5e2b37991156e0adfc18/README.md)。
- 中文 README 自述对应版本：2024-08-10，链接到 [uber-go/guide @ 97314412dfcb60c3f4f0372bc311f9ddcdac54a9](https://github.com/uber-go/guide/commit/97314412dfcb60c3f4f0372bc311f9ddcdac54a9)；这是来源自述，不表示本 skill 与英文最新版逐条同步。
- 英文原始指南：[Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)。原指南由 Prashant Varanasi 和 Simon Newton 发起，后由 Uber 及社区贡献者维护；中文版本由 xxjwxc 及翻译贡献者整理。
- 整理日期：2026-11-04。
- 中文仓库以 Apache-2.0 许可证发布，许可证副本见 [LICENSE](LICENSE)。

本目录是对原指南的**摘要、重组与改写**，并新增了执行步骤、审查输出要求和项目适配；不是 Uber 官方 skill，也不是原文全文副本。参考文档中的代码例子经过重新编写。

## 覆盖范围

- `references/correctness.md`：接口与接收器、错误、类型断言、集合所有权、mutex/atomic、资源释放、channel、goroutine、全局状态、初始化/退出、时间、枚举、序列化。
- `references/style-and-patterns.md`：命名与代码组织、import、控制流、声明和初始化、字符串格式、表驱动测试、功能选项、容器容量与热点转换优化、lint。
- 省略与执行无关的背景介绍、历史 benchmark 数字、贡献说明、星标图等内容。

## 有意的适配与澄清

| 原文建议或表述 | 本 skill 的处理 |
| --- | --- |
| 推荐 `go.uber.org/atomic` 包装原始类型 | 保留类型安全和一致原子访问原则，当前 Go 优先标准库类型化 atomic |
| 并行子测试循环中添加 `tt := tt` | 按模块/文件实际 Go 语言版本区分，不能忽略 Go 1.22 的循环变量变化 |
| 使用 `golint` | 保留 lint 原则，沿用项目工具，不要求已停止维护的工具 |
| 接口“本身就是一个指针”、接口指针“永远无意义” | 明确接口值不等同普通指针，保留接口变量覆写等合法例外 |
| JSON “不支持 time.Duration” | 澄清默认是纳秒整数，文本时长须显式转换，而非无法序列化 |
| 避免可变全局，优先依赖注入 | 遵循 InkSpace 明确的数据库单例约定，不擅自进行架构迁移 |
| 未导出顶层变量/常量使用 `_` 前缀 | 标为 Uber 团队风格，项目包内一致性优先 |
| 优先 nil slice、限制公共类型嵌入 | 保留序列化契约、框架模型嵌入与有意 API 设计的例外 |
| 一次处理错误 | 检查实际调用链与项目日志规则，区别无意义重复和有意的边界观测 |

附加的 typed nil、浅拷贝深度、循环内 defer、阻塞点取消、测试 goroutine 限制等以“适配补充”标识。更新时先核对固定版本与目标版本差异，再核对 `go.mod`、项目规则和 Go 官方文档；保留来源记录，不能静默把项目偏好说成 Uber 原文硬性要求。
