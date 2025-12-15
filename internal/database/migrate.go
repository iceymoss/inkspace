package database

import (
	"fmt"
	"log"

	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

// MigrateDB 执行数据库迁移
func MigrateDB() error {
	log.Println("开始数据库迁移...")

	// 自动迁移所有表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	// 创建索引
	if err := createIndexes(); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	// 创建外键约束（可选）
	if err := createForeignKeys(); err != nil {
		log.Printf("警告: 创建外键失败: %v", err)
		// 不返回错误，因为有些环境可能不支持外键
	}

	log.Println("数据库迁移完成！")
	return nil
}

// createIndexes 创建额外的索引
func createIndexes() error {
	log.Println("创建索引...")

	// 检查并创建文章表的组合索引
	var count int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'articles' AND index_name = 'idx_articles_top_status_created'").Scan(&count)
	if count == 0 {
		if err := DB.Exec(`CREATE INDEX idx_articles_top_status_created ON articles(is_top DESC, status, created_at DESC)`).Error; err != nil {
			log.Printf("警告: 创建组合索引失败: %v", err)
		}
	}

	// 检查并创建评论表的状态和创建时间索引（已由GORM创建，此处跳过）
	// idx_status_created 已在模型中定义

	// 检查并创建全文搜索索引（MySQL 5.7+）
	DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'articles' AND index_name = 'idx_articles_title_content'").Scan(&count)
	if count == 0 {
		if err := DB.Exec(`CREATE FULLTEXT INDEX idx_articles_title_content ON articles(title, content)`).Error; err != nil {
			log.Printf("警告: 创建全文索引失败: %v (可以忽略)", err)
		}
	}

	log.Println("索引创建完成")
	return nil
}

// createForeignKeys 创建外键约束
func createForeignKeys() error {
	log.Println("创建外键约束...")

	// 文章表外键
	DB.Exec(`
		ALTER TABLE articles 
		ADD CONSTRAINT fk_articles_author 
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	`)

	DB.Exec(`
		ALTER TABLE articles 
		ADD CONSTRAINT fk_articles_category 
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	`)

	// 评论表外键（注意：article_id 和 work_id 可以为 NULL）
	DB.Exec(`
		ALTER TABLE comments 
		ADD CONSTRAINT fk_comments_article 
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
	`)

	DB.Exec(`
		ALTER TABLE comments 
		ADD CONSTRAINT fk_comments_work 
		FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE
	`)

	DB.Exec(`
		ALTER TABLE comments 
		ADD CONSTRAINT fk_comments_user 
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	`)

	DB.Exec(`
		ALTER TABLE comments 
		ADD CONSTRAINT fk_comments_parent 
		FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
	`)

	// 作品表外键
	DB.Exec(`
		ALTER TABLE works 
		ADD CONSTRAINT fk_works_author 
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	`)

	// 文章标签关联表外键
	DB.Exec(`
		ALTER TABLE article_tags 
		ADD CONSTRAINT fk_article_tags_article 
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
	`)

	DB.Exec(`
		ALTER TABLE article_tags 
		ADD CONSTRAINT fk_article_tags_tag 
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	`)

	log.Println("外键约束创建完成")
	return nil
}

// DropAllTables 删除所有表（危险操作，仅用于开发环境）
func DropAllTables() error {
	log.Println("警告: 正在删除所有表...")

	// 禁用外键检查
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 删除所有表
	tables := []string{
		"article_tags",
		"comments",
		"articles",
		"works",
		"tags",
		"categories",
		"users",
	}

	for _, table := range tables {
		if err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			log.Printf("删除表 %s 失败: %v", table, err)
		}
	}

	// 启用外键检查
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("所有表已删除")
	return nil
}

// SyncCounters 同步所有计数器字段
func SyncCounters() error {
	log.Println("开始同步计数器...")

	// 同步用户的文章数
	if err := DB.Exec(`
		UPDATE users u 
		SET article_count = (
			SELECT COUNT(*) FROM articles a 
			WHERE a.author_id = u.id AND a.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步用户的评论数
	if err := DB.Exec(`
		UPDATE users u 
		SET comment_count = (
			SELECT COUNT(*) FROM comments c 
			WHERE c.user_id = u.id AND c.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步文章的评论数
	if err := DB.Exec(`
		UPDATE articles a 
		SET comment_count = (
			SELECT COUNT(*) FROM comments c 
			WHERE c.article_id = a.id AND c.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步分类的文章数
	if err := DB.Exec(`
		UPDATE categories cat 
		SET article_count = (
			SELECT COUNT(*) FROM articles a 
			WHERE a.category_id = cat.id AND a.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步标签的文章数
	if err := DB.Exec(`
		UPDATE tags t 
		SET article_count = (
			SELECT COUNT(*) FROM article_tags at 
			INNER JOIN articles a ON at.article_id = a.id 
			WHERE at.tag_id = t.id AND a.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步评论的回复数
	if err := DB.Exec(`
		UPDATE comments c 
		SET reply_count = (
			SELECT COUNT(*) FROM comments r 
			WHERE r.parent_id = c.id AND r.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步作品的评论数
	if err := DB.Exec(`
		UPDATE works w 
		SET comment_count = (
			SELECT COUNT(*) FROM comments c 
			WHERE c.work_id = w.id AND c.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步用户的关注数
	if err := DB.Exec(`
		UPDATE users u
		SET following_count = (
			SELECT COUNT(*) FROM user_follows uf
			WHERE uf.follower_id = u.id AND uf.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步用户的粉丝数
	if err := DB.Exec(`
		UPDATE users u
		SET follower_count = (
			SELECT COUNT(*) FROM user_follows uf
			WHERE uf.following_id = u.id AND uf.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步用户的收藏数
	if err := DB.Exec(`
		UPDATE users u
		SET favorite_count = (
			SELECT COUNT(*) FROM article_favorites af
			WHERE af.user_id = u.id AND af.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	// 同步文章的收藏数
	if err := DB.Exec(`
		UPDATE articles a
		SET favorite_count = (
			SELECT COUNT(*) FROM article_favorites af
			WHERE af.article_id = a.id AND af.deleted_at IS NULL
		)
	`).Error; err != nil {
		return err
	}

	log.Println("计数器同步完成")
	return nil
}

// GetTableInfo 获取表信息
func GetTableInfo(tableName string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := DB.Raw(fmt.Sprintf("DESCRIBE %s", tableName)).Scan(&result).Error
	return result, err
}

// GetAllTables 获取所有表名
func GetAllTables() ([]string, error) {
	var tables []string
	err := DB.Raw(`
		SELECT TABLE_NAME 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME
	`).Scan(&tables).Error
	return tables, err
}

// CheckDatabaseHealth 检查数据库健康状态
func CheckDatabaseHealth() error {
	// 检查数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库ping失败: %w", err)
	}

	// 检查表是否存在
	requiredTables := []string{"users", "articles", "categories", "tags", "comments", "works"}
	tables, err := GetAllTables()
	if err != nil {
		return fmt.Errorf("获取表列表失败: %w", err)
	}

	tableMap := make(map[string]bool)
	for _, table := range tables {
		tableMap[table] = true
	}

	for _, required := range requiredTables {
		if !tableMap[required] {
			return fmt.Errorf("缺少必需的表: %s", required)
		}
	}

	log.Println("数据库健康检查通过")
	return nil
}

// BackupDatabase 备份数据库（仅结构，不包含数据）
func BackupDatabaseSchema(outputFile string) error {
	// 这里只是一个示例，实际应该使用 mysqldump
	log.Printf("数据库schema备份功能需要使用 mysqldump 命令")
	log.Printf("建议使用: mysqldump -u username -p database_name --no-data > %s", outputFile)
	return nil
}

// Hooks for maintaining counters

// AfterCreateArticle 文章创建后的钩子
func AfterCreateArticle(db *gorm.DB, article *models.Article) error {
	// 更新用户文章数
	return db.Model(&models.User{}).
		Where("id = ?", article.AuthorID).
		UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error
}

// AfterDeleteArticle 文章删除后的钩子
func AfterDeleteArticle(db *gorm.DB, article *models.Article) error {
	// 更新用户文章数
	return db.Model(&models.User{}).
		Where("id = ?", article.AuthorID).
		UpdateColumn("article_count", gorm.Expr("article_count - ?", 1)).Error
}

// AfterCreateComment 评论创建后的钩子
func AfterCreateComment(db *gorm.DB, comment *models.Comment) error {
	// 更新文章评论数
	if err := db.Model(&models.Article{}).
		Where("id = ?", comment.ArticleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		return err
	}

	// 如果是登录用户评论，更新用户评论数
	if comment.UserID > 0 {
		if err := db.Model(&models.User{}).
			Where("id = ?", comment.UserID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
	}

	// 如果是回复，更新父评论的回复数
	if comment.ParentID != nil {
		if err := db.Model(&models.Comment{}).
			Where("id = ?", *comment.ParentID).
			UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}

// AfterDeleteComment 评论删除后的钩子
func AfterDeleteComment(db *gorm.DB, comment *models.Comment) error {
	// 更新文章评论数
	if err := db.Model(&models.Article{}).
		Where("id = ?", comment.ArticleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
		return err
	}

	// 如果是登录用户评论，更新用户评论数
	if comment.UserID > 0 {
		if err := db.Model(&models.User{}).
			Where("id = ?", comment.UserID).
			UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			return err
		}
	}

	// 如果是回复，更新父评论的回复数
	if comment.ParentID != nil {
		if err := db.Model(&models.Comment{}).
			Where("id = ?", *comment.ParentID).
			UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}
