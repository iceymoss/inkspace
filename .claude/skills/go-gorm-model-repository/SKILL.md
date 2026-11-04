---
name: go-gorm-model-repository
description: 编写或修改 Go GORM 模型与查询时，使用显式表名和集中列名映射，减少重复硬编码。
---

# Go GORM 模型与查询约定

## 核心规则

- 表模型定义在 `objects` 包，通过 `TableName()` 显式声明表名。
- 每个模型维护 `<Model>Columns`，集中定义数据库列名。
- 查询、更新等操作引用列名映射，避免到处直接写列名字符串。
- 新增、删除或重命名字段时，同步维护模型与列名映射。
- 字段类型、GORM 标签和关联关系参考现有模型，不擅自添加约束。

## 模型示例

```go
package objects

type User struct {
	ID    string `gorm:"size:128;primarykey"`
	Name  string `gorm:"size:128"`
	Email string `gorm:"size:256"`
}

func (User) TableName() string {
	return "users"
}

var UserColumns = struct {
	ID    string
	Name  string
	Email string
}{
	ID:    "id",
	Name:  "name",
	Email: "email",
}
```

`Columns` 只包含实际落库字段，不包含关联对象字段。存在 `gorm:"column:..."` 时，以指定的列名为准。

## 查询示例

以下示例复用项目已有的 Repository 和数据库依赖：

```go
func (u *user) GetByUserID(
	ctx context.Context,
	userID string,
) (*objects.User, error) {
	db := driver.GetDB(ctx)

	var userObject objects.User

	res := db.Model(&objects.User{}).
		Where(objects.UserColumns.ID+" = ?", userID).
		First(&userObject)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}

	return &userObject, nil
}
```

## 数据访问约定

- 条件值使用参数绑定，禁止拼接用户输入。
- 组织、团队等过滤条件按接口要求保留。
- 需要关联时显式 `Preload`，参数使用模型的关联字段名。
- 单条查询默认：查到返回对象，未找到返回 `nil, nil`，执行失败返回错误；已有接口另有约定时遵循原约定。
- 动态列名和排序字段使用白名单，不能直接接收外部输入。

列名映射用于集中维护，不会自动验证数据库结构；修改后检查映射一致性，并执行 `gofmt` 和相关编译或测试。