# 风格、测试与模式

先遵循包内一致性。这里的风格默认值不能推动无关文件重排或批量重命名；项目适配见 [SKILL.md](../SKILL.md)。

## 命名与代码组织

- 新包名使用简短、小写、通常为单数的职责名，不用下划线；避开信息量不足的 `common` / `util` / `shared` / `lib`。已有项目包名不因此迁移。
- 函数名使用 MixedCaps；测试函数可用下划线组织场景。标识符避免遮蔽 `error`、`string`、`len` 等预声明名字。
- import 至少区分标准库与其他依赖，组内排序交给项目工具。已有项目内包分组设置时沿用，不强行压成两组。
- 导入名称与路径尾部不直观对应时显式标别名；其余仅在冲突或确有消歧需要时使用。
- 相关声明用 `const (...)` / `var (...)` / `type (...)` 成组，不把无关概念硬塞进同一个组。
- 类型后放构造函数，相关接收器方法相邻，粗略按调用关系组织；辅助函数靠后。只在当前改动有益时调整。
- 行长以可读为准，原指南约 99 字符是建议而非硬限制；避免靠手工断行对抗 `gofmt`。
- Uber 原文给未导出的包级常量/变量加 `_` 前缀，错误值 `errXxx` 例外。这是团队偏好，不是 Go 语言要求；是否采用由当前包约定决定。

## 控制流与参数

- 先处理错误和特殊分支，通过 `return` / `continue` 保持正常路径平直。
- 已返回的分支后通常无需 `else`；简单默认值可以先赋默认再条件覆盖，但不因此提前执行昂贵或有副作用的计算。
- 变量作用域尽量小；若结果后续仍需使用，宁可在外层声明，也不为缩作用域增加嵌套。
- 裸 `true` / `false` 或数字难理解时，优先用有语义的类型、常量或配置；既有调用可就近添加 `/* 参数含义 */`。
- 多行或转义密集的文本优先原始字符串；保证字符串实际内容不变。

## 声明与初始化

| 场景 | 默认写法 |
| --- | --- |
| 局部变量有明确初值 | `name := value` |
| 使用零值 | `var name T` |
| 顶层变量可推断类型 | `var name = value` |
| 需要不同的目标类型 | 显式写类型，如 `var e error = customError{}` |
| 结构体有初始化字段 | `User{Name: "alice"}` |
| 结构体指针 | `&User{Name: "alice"}` |
| 动态填充 map | `make(map[K]V, hint)`，容量未知可省略 hint |
| 固定内容 map | `map[K]V{key: value}` |

- 结构体字面量优先带字段名；省略无额外含义的零值字段。表驱动测试中的显式零值可保留以表达期望。
- 很小且清晰的测试表可以使用位置字段，跨包结构体和易变结构体仍用字段名。
- 普通空 slice 用零值；已知将追加的数量时用 `make([]T, 0, n)`。若按索引填写 n 个元素，则用 `make([]T, n)`，不要混淆长度与容量。
- 复用的格式字符串用 `const`；Printf 风格包装函数以 `f` 结尾并保持格式串与可变参数的转发关系，便于静态检查。

## 表驱动测试

- 相同执行逻辑的多个场景使用 `tests` + `t.Run`，用 `name` 标场景，输入用 `give`，期望用 `want`。行为不同的测试不强行合表。
- 断言可观察行为、失败类型和边界；公共错误契约用 `errors.Is/As` 验证。前置条件失败立即终止当前子测试。
- 并行执行前确认依赖与全局状态隔离；时间、数据库、环境变量等共享状态不能因为用了子测试就自动安全。
- **Go 版本适配**：Go 1.22+ 语言语义下，`for _, tt := range tests` 每次迭代产生独立变量，不需要机械添加 `tt := tt`。旧语义或对已声明变量使用 `=` 的循环仍需检查捕获问题。
- 使用标准库 `testing` 和项目已有辅助工具，不照搬原文示例而新增断言库。

完整测试示例（标准库、无额外依赖）：

```go
package example

import (
	"net"
	"testing"
)

func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		name     string
		give     string
		wantHost string
		wantPort string
	}{
		{name: "ip", give: "192.0.2.1:80", wantHost: "192.0.2.1", wantPort: "80"},
		{name: "empty host", give: ":8080", wantHost: "", wantPort: "8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, err := net.SplitHostPort(tt.give)
			if err != nil {
				t.Fatalf("SplitHostPort(%q): %v", tt.give, err)
			}
			if host != tt.wantHost || port != tt.wantPort {
				t.Errorf("got (%q, %q), want (%q, %q)", host, port, tt.wantHost, tt.wantPort)
			}
		})
	}
}
```

## 功能选项（Functional options）

- 公共构造函数有持续增长的可选参数、尤其已有三个或更多参数时，考虑 `New(required, opts ...Option)`；简单且稳定的 API 保持简单。
- 必选参数仍是显式参数，可选项提供有效默认值，命名为 `WithXxx`。先设默认值，再按调用顺序应用选项。
- 原指南偏好 `Option interface { apply(*options) }`，由具名类型实现，利于调试和检查选项值；并非所有 option 都天然可比较，仍取决于实现字段。
- 项目已有函数型 option 时保持一致，不仅为符合偏好替换成熟 API。必要时配置结构体比额外 option 抽象更清晰。
- 可变集合选项按所有权规则复制；验证非法值。新增选项保持旧调用兼容，别把必需依赖藏进无效默认值。

## 性能与工具

- 对原语与字符串之间的转换优先考虑 `strconv`，普通格式化和错误上下文仍用 `fmt`。
- 重复使用固定 `[]byte` 时把转换移出热点循环；只在接收方不会意外修改或长期持有共享缓冲时复用。不要引入 `unsafe` 零拷贝来满足此规则。
- 已知大小时给 slice / map 合理容量；map 大小是提示，不保证之后零分配。容量不要来自不受控输入而无上限预分配。
- 原指南的微优化面向高频路径。先 benchmark/profile，再决定是否值得增加复杂性，优先保持正确性和资源释放。
- 使用 `gofmt`、`go vet` 和仓库已有 lint；若已配置 `goimports`、`staticcheck`、`errcheck` 或 `golangci-lint` 则沿用版本与配置。
- 原文的 `golint` 已停止维护，不作为新工具要求；不因为引用了历史工具清单而改变项目工具链。
