---
name: new-model
description: 新增数据库模型/表或字段（GORM + MySQL）。当用户说加表、新增 model、加数据库字段时触发。
---

新增数据库表 / 字段。本项目用 GORM struct tag + `AutoMigrate`，**只支持 MySQL**。

## 步骤

### 1. 与用户对齐 schema（必须）
按 [`database-model.md`](../../rules/database-model.md)，问清楚：表名、字段、类型、索引、关联、是否软删除。

### 2. 写 model
在 `internal/models/<name>.go` 定义：
```go
type Xxx struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`   // 需要软删除时
    Name      string         `gorm:"size:100;not null;index" json:"name" binding:"required"`
    // JSON 数据用 longtext/text 存字符串，不用 JSONB
}
```
- 关联用 `foreignKey` / `many2many:xxx` / `constraint:OnDelete:...`
- 高频查询字段加 `index`
- 同文件里加 `XxxRequest`（绑定）和 `ToResponse()`（对外 DTO，隐藏敏感字段）

### 3. 注册迁移
把 `&models.Xxx{}` 加入 `internal/database/migrate.go` 的 `AutoMigrate(...)` 列表，否则不会建表。

### 4. 验证
`go run cmd/server/main.go` 启动会自动迁移；确认表结构符合预期。

## 严格约束
- 只写 MySQL 兼容 schema，不写 SQLite/PostgreSQL 分支，不用 JSONB / AUTO_INCREMENT 手写主键
- **破坏性变更**（删列、改类型、加唯一约束）AutoMigrate 不安全：先与用户确认，并在 `scripts/` 提供迁移 SQL + 回滚方案
- 新增 NOT NULL 字段要给默认值，兼容存量行
- 绝不修改 `.env` / `config/*.yaml`，绝不 drop/truncate 现有表
