package service

import (
	"errors"
	"fmt"
	"github.com/iceymoss/inkspace/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

func knowledgeQuotaLimit(db *gorm.DB, key string, fallback int64) (int64, error) {
	var setting models.Setting
	err := db.Where("`key` = ?", key).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fallback, nil
	}
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil || n < 0 {
		zap.L().Warn("invalid knowledge quota", zap.String("key", key))
		return 0, nil
	}
	return n, nil
}

func checkKnowledgeStorageQuota(db *gorm.DB, owner uint, added int64) error {
	limit, err := knowledgeQuotaLimit(db, QuotaMaxStorage, 1024)
	if err != nil || limit == 0 {
		return err
	}
	var used int64
	if err := db.Model(&models.Doc{}).Select("COALESCE(SUM(attachments.file_size),0)").Joins("JOIN attachments ON attachments.id=docs.attachment_id AND attachments.deleted_at IS NULL").Where("docs.owner_id=?", owner).Scan(&used).Error; err != nil {
		return err
	}
	if float64(used)+float64(added) > float64(limit)*1048576 {
		return fmt.Errorf("%w：知识库附件存储已达上限（%d MB）", ErrKnowledgeInvalid, limit)
	}
	return nil
}

func checkKnowledgeCountQuota(db *gorm.DB, key string, fallback int64, model interface{}, condition string, id uint) error {
	limit, err := knowledgeQuotaLimit(db, key, fallback)
	if err != nil || limit == 0 {
		return err
	}
	var count int64
	if err := db.Model(model).Where(condition, id).Count(&count).Error; err != nil {
		return err
	}
	if count >= limit {
		return fmt.Errorf("%w：知识库配额已达上限（%d），无法继续创建", ErrKnowledgeInvalid, limit)
	}
	return nil
}
