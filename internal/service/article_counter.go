package service

import (
	"github.com/iceymoss/inkspace/internal/models"

	"gorm.io/gorm"
)

// articleCounterChange is only for interaction updates, not content/state edits.
// UpdateColumn alone skips GORM timestamps but cannot suppress legacy MySQL
// ON UPDATE columns. An explicit self-assignment preserves editorial lastmod.
// Callers supply a trusted ArticleColumns counter and retain their own scope/tx.
func articleCounterChange(column string, delta int) map[string]any {
	return map[string]any{
		column:                          gorm.Expr(column+" + ?", delta),
		models.ArticleColumns.UpdatedAt: gorm.Expr(models.ArticleColumns.UpdatedAt),
	}
}
