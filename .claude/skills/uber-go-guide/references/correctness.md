# 正确性与生命周期

规则以中文指南为基础；标为“适配补充”的内容用于消除原文歧义或处理当前 Go 行为。项目例外见 [SKILL.md](../SKILL.md)。

## 错误处理

| 调用方需要什么 | 选择 |
| --- | --- |
| 仅需静态错误消息 | `errors.New("message")` |
| 仅需动态上下文 | `fmt.Errorf("read %q", name)` |
| 匹配固定错误 | 包级 sentinel，配合 `errors.Is` |
| 提取结构化错误信息 | 自定义 `XxxError`，配合 `errors.As` |

- 导出 sentinel 使用 `ErrXxx`，未导出使用 `errXxx`；错误类型以 `Error` 结尾。
- 没有新增上下文时原样返回。补充上下文时默认用 `%w` 保留错误链；明确要隐藏底层错误身份时用 `%v`，并检查是否破坏既有契约。
- 上下文描述动作与对象，例如 `get user %q: %w`，避免每层重复添加 `failed to`。
- 暴露的 sentinel、错误类型和可解包的底层错误都是 API 的一部分，应记录并测试匹配行为。使用 `errors.Is/As`，不按错误字符串分支。
- 每个错误在合适的边界处理一次：恢复并继续，或包装向上传递。普通内部调用链避免每层都记录后再返回相同错误。
- 可恢复错误通过 `error` 返回；`panic/recover` 不作为业务分支。初始化中对固定、确定性资源使用 `Must` 需有启动失败即中止的理由。
- 测试前置条件失败用 `t.Fatal` / `t.Fatalf`。适配补充：它们必须在所属测试 goroutine 中调用；工作 goroutine 将结果传回测试。

示例片段（假设 `load` 已定义）：

```go
user, err := load(id)
if err != nil {
	return nil, fmt.Errorf("load user %d: %w", id, err)
}
return user, nil
```

## 接口、接收器与封装

- 通常传接口值 `io.Reader`，而不是 `*io.Reader`；需要修改具体对象时，将实现接口的对象指针存入接口值。
- 适配补充：接口不等同于普通指针；需要覆写接口变量本身的少数 API 可以合法接收接口指针，如 `errors.As` 的目标参数。审查时先确认用途。
- 对有明确实现承诺的类型添加编译期检查，如 `var _ http.Handler = (*Handler)(nil)`。值类型用 `T{}`，指针类型用 `(*T)(nil)`。
- `T` 的方法集包含值接收器方法；`*T` 的方法集还包含指针接收器方法。可寻址值的调用语法糖不改变接口满足关系；map 元素不可寻址。
- 修改对象、避免大对象复制或包含锁时使用指针接收器，保持同一类型接收器选择一致。适配补充：含锁/原子状态的对象首次使用后不能复制。
- 公共类型以具名字段组合、显式方法委托隐藏实现；嵌入只用于确实希望暴露整套字段/方法且零值语义合理的场景。
- 有意嵌入的字段置顶，与普通字段空行分组。mutex 始终使用具名字段，例如 `mu sync.Mutex`，而非匿名嵌入。
- 对可能失败的类型断言用 `value, ok := input.(T)` 并处理失败。适配补充：接口中装入 typed nil 后接口本身不为 nil，返回错误时尤其要检查。

## slice / map 的所有权

- 保存调用者传入的可变集合、或返回内部集合快照时，明确所有权；需要隔离修改时在边界复制，而不是只复制 slice 头或 map 引用。
- 有意共享时文档说明可变性和同步责任，不能把所有函数参数机械地复制一遍。
- 快照在锁保护范围内完成复制；解锁后返回内部 map 仍会暴露竞态。
- 适配补充：`copy`、`slices.Clone`、`maps.Clone` 都是浅拷贝；嵌套 map、slice、指针按实际隔离需求进一步复制。
- nil slice 可 `append`、`range`、`len`，检查空集合用 `len(s) == 0`；nil map 可读不可写。初始化 map 不代表并发读写安全。
- 返回 nil 还是空集合由 API 契约决定，特别检查 JSON `null` 与 `[]` / `{}` 的差异。

## 锁、原子状态与资源

- 使用零值可用的 `sync.Mutex` / `sync.RWMutex` 具名值字段，通常不需要 `new(sync.Mutex)`。
- 成功获取锁、文件、rows 等资源后就近安排释放；优先 `defer`，让多个返回路径共享清理逻辑。
- 适配补充：循环内 `defer` 在外层函数退出时才执行；资源需每轮及时释放时提取小函数或显式关闭。影响持久化结果的关闭/刷新错误也要处理。
- 原子共享变量的所有读写都通过一致的原子 API；当前 Go 优先使用标准库 `sync/atomic` 的 `atomic.Bool`、`atomic.Int64` 等类型化包装。
- 多个字段的一致性约束不能仅靠各自原子化保证，必要时用同一把锁保护整体状态。

## channel 与 goroutine

- channel 默认无缓冲，或为单次结果/信号使用容量 1。容量大于 1 需要解释队列界限、吞吐与背压、满队列时的阻塞或丢弃语义。
- 每个 goroutine 都有可预测的结束条件，或可触发的取消信号，并有等待其退出的方式：单个可用 `done` channel，多个可用 `sync.WaitGroup`。
- 包持有后台任务时提供 `Close` / `Stop` / `Shutdown` 等生命周期入口；该入口既发出停止信号，也等待任务退出。
- 在显式生命周期中启动 goroutine，而不是在 `init()` 中启动。
- 适配补充：检查所有阻塞点，包括发送、接收和 I/O；只在循环顶部检查取消不足以打断后续永久阻塞。底层调用需支持取消/超时。
- 明确 channel 关闭方，避免重复关闭、向已关闭 channel 发送。取消后停止 ticker、释放资源，确保等待路径实际可达。
- `WaitGroup.Add` 在启动待等待的工作之前调用，工作中 `defer Done`。若使用较新便捷 API，先确认项目 Go 版本。
- 测试正常完成、取消和关闭路径；仓库已有泄漏检查工具时沿用，不为此自动新增依赖。

生命周期示例，`flush` 自身也必须遵守传入的 context：

```go
func runFlusher(ctx context.Context, interval time.Duration, flush func(context.Context)) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval) // interval 由调用者保证为正。
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				flush(ctx)
			}
		}
	}()
	return done // 调用者取消 context 后，接收 done 等待退出。
}
```

## 初始化、全局状态与退出

- 一般通过显式初始化函数组织可能失败的工作。`init()` 限于确定性注册或预计算，避免环境读取、I/O 和初始化顺序耦合。
- 新依赖优先显式表达，减少可变全局状态；已有架构规定的单例按项目约定使用。测试替换全局状态时必须隔离并恢复，不能与并行测试共享修改。
- `os.Exit`、`log.Fatal` 及等效的库调用集中在 `main()`，理想上只有一个退出点；其余函数返回错误或退出码。
- 可把逻辑提取到可测试的 `run(...) error` 中。`os.Exit` / `log.Fatal` 不执行 defer，资源清理应在返回 `main` 之前完成。

## 时间、枚举与序列化

- 瞬时时间用 `time.Time`，时长用 `time.Duration`；比较用时间方法。自然日运算用 `AddDate`，经过固定时长用 `Add`。
- 外部协议不能直接使用时间类型时，时间戳优先采用已约定的 RFC 3339，数值时长在字段名写明单位，例如 `IntervalMillis`。
- 适配补充：`encoding/json` 默认将 `time.Duration` 编码为底层整数（纳秒），不会自动生成 `"5s"`；文本时长需明确的编码/解析方式。
- 枚举通常从 `iota + 1` 开始保留零值无效；零值确为合理默认时可以从 0 开始。已有持久化或协议枚举保持数值兼容。
- 对 JSON/YAML 等序列化字段写显式 tag，把外部字段名与 Go 标识符解耦；重命名时同时核对 tag、省略行为和兼容性。
