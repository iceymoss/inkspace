---
name: go-imports
description: zhiyan Go 文件 import 分组与排序规范。写或改 Go 文件、整理 import、新增依赖时按此分组。顺序固定为 内置包 → 项目包 → 第三方包,三组之间用空行分隔。当用户说 import 顺序、导入分组、整理 import、依赖顺序、包怎么排时触发。
---

zhiyan 的 Go 文件 import 必须分成三组,**按此顺序**,组间用一个空行分隔:

1. **内置包**(标准库)—— 无域名前缀,如 `context`、`fmt`、`net/http`、`strings`。
2. **项目包**—— 本仓库,前缀 `github.com/iceymoss/zhiyan/...`。
3. **第三方包**—— 其余带域名的外部依赖,如 `github.com/gin-gonic/gin`、`gorm.io/gorm`、`go.uber.org/zap`。

组内按 import 路径字母序排列。

## 标准写法

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/iceymoss/zhiyan/common"
	"github.com/iceymoss/zhiyan/internal/dto"
	"github.com/iceymoss/zhiyan/internal/error_codes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
```

## 规则细节

- 三组都存在时,严格 内置 → 项目 → 第三方;某组为空则跳过该组,但不改变相对顺序。
- 组与组之间**只用一个空行**;组内**不留空行**。
- 只有一个 import 时不必分组,直接单行或单组即可。
- 带别名/点导入/下划线导入的,按其所属包归到对应组内,不单独成组。
  例:`_ "github.com/lib/pq"` 归第三方组,`_ "github.com/iceymoss/zhiyan/internal/xxx"` 归项目组。

## 工具注意(重要)

这个顺序**不是** Go 官方默认,靠工具需注意:

- `gofmt`(项目 PostToolUse hook 用的)只在**已分好的空行组内**排序,**不会重排或合并组**。所以三组的划分和顺序要**手动摆对**,gofmt 只帮你把组内排字母序 —— 摆错组它不会纠正。
- `goimports -local github.com/iceymoss/zhiyan` 的固定顺序是 内置 → 第三方 → **项目**(项目垫底),与本规范**相反**,不要用它来自动整理,会把项目包排到第三方后面。
- 需要自动化就用 `gci`,它支持自定义分组顺序:

  ```bash
  gci write --section standard \
            --section "prefix(github.com/iceymoss/zhiyan)" \
            --section default \
            <file-or-dir>
  ```

  `standard` = 内置,`prefix(...)` = 项目,`default` = 第三方,顺序即本规范。

## 写代码时

新增或修改 Go 文件后,自查 import 是否为 内置 → 项目 → 第三方 三组、组间单空行、组内字母序。gofmt 保存钩子不保证组顺序,以本规范为准。
