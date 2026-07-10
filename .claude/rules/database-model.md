---
description: 数据模型与数据库迁移规则（GORM + MySQL）
globs: ["internal/models/**/*.go", "internal/database/**/*.go"]
---

- schema 用 GORM struct tag 定义，`internal/database/migrate.go` 的 `AutoMigrate` 自动同步；改字段=改 model 的 `gorm:` tag
- 只支持 MySQL（`utf8mb4`）。不要写 SQLite/PostgreSQL 兼容代码，不要用 JSONB（用 `longtext`/`text` 存 JSON 字符串）
- 主键用 `gorm:"primarykey"` 交给 GORM；软删除内嵌 `gorm.DeletedAt`（`gorm:"index" json:"-"`）
- 新增 model 后，必须把 `&models.Xxx{}` 加进 `migrate.go` 的 AutoMigrate 列表，否则不会建表
- **破坏性变更**（删列、改类型、加唯一约束/外键）AutoMigrate 不安全处理：先与用户确认，并在 `scripts/` 提供可重复执行的迁移 SQL
- 请求体（`XxxRequest`，带 `binding:"required"`）和响应转换（`ToResponse()`）与 model 放在同一文件；对外 DTO 隐藏密码/敏感字段
- 关联用 GORM tag（`foreignKey` / `many2many` / `constraint:OnDelete:...`），并为高频查询字段加 `index`
- 不要修改 `.env` 与 `config/*.yaml` 里的连接配置
