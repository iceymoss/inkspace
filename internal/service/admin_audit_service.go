package service

import (
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

type AdminAuditService struct{}

func NewAdminAuditService() *AdminAuditService { return &AdminAuditService{} }
func (s *AdminAuditService) Log(tx *gorm.DB, entry *models.AdminAuditLog) error {
	return tx.Create(entry).Error
}
