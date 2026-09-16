package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"

	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/gorm"
)

const (
	QuotaMaxWorkspaces = "knowledge.max_workspaces_per_user"
	QuotaMaxDocs       = "knowledge.max_docs_per_workspace"
	QuotaMaxStorage    = "knowledge.max_storage_per_user_mb"
)

type KnowledgeAdminService struct{}

func NewKnowledgeAdminService() *KnowledgeAdminService { return &KnowledgeAdminService{} }

func pageQuery(q *gorm.DB, page, size int) *gorm.DB {
	return q.Order("id DESC").Offset((page - 1) * size).Limit(size)
}

type AdminWorkspace struct {
	models.Workspace
	OwnerName string `json:"owner_name"`
}

func (s *KnowledgeAdminService) Workspaces(page, size int, keyword string, owner *uint, public, audit *bool) ([]AdminWorkspace, int64, error) {
	v := make([]AdminWorkspace, 0)
	var n int64
	q := database.DB.Model(&models.Workspace{})
	if keyword != "" {
		q = q.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if owner != nil {
		q = q.Where("owner_id = ?", *owner)
	}
	if public != nil {
		q = q.Where("is_public = ?", *public)
	}
	if audit != nil {
		q = q.Where("audit_status = ?", *audit)
	}
	if err := q.Count(&n).Error; err != nil {
		return nil, 0, err
	}
	err := pageQuery(q.Select("workspaces.*, "+adminOwnerNameSQL("workspaces")), page, size).Find(&v).Error
	return v, n, err
}
func (s *KnowledgeAdminService) Docs(page, size int, keyword string, owner, workspace *uint, status, audit *int, kind string) ([]AdminDocMetadata, int64, error) {
	v := make([]AdminDocMetadata, 0)
	var n int64
	q := database.DB.Model(&models.Doc{})
	if keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	if owner != nil {
		q = q.Where("owner_id=?", *owner)
	}
	if workspace != nil {
		q = q.Where("workspace_id=?", *workspace)
	}
	if status != nil {
		q = q.Where("status=?", *status)
	}
	if audit != nil {
		q = q.Where("audit_status=?", *audit)
	}
	if kind != "" {
		q = q.Where("kind=?", kind)
	}
	if err := q.Count(&n).Error; err != nil {
		return nil, 0, err
	}
	err := pageQuery(q.Select("id,workspace_id,owner_id,title,kind,status,audit_status,audit_reason,view_count,created_at,updated_at,"+adminOwnerNameSQL("docs")), page, size).Find(&v).Error
	return v, n, err
}
func (s *KnowledgeAdminService) Audit(kind string, id, admin uint, status int8, reason, username, ip string) error {
	kind = strings.TrimSuffix(kind, "s")
	if kind != "workspace" && kind != "doc" {
		return ErrKnowledgeInvalid
	}
	if status != 0 && status != 1 {
		return errors.New("invalid audit status")
	}
	var title string
	var owner uint
	changed := false
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if kind == "workspace" {
			var w models.Workspace
			if err := tx.First(&w, id).Error; err != nil {
				return err
			}
			title = w.Name
			owner = w.OwnerID
			r := tx.Model(&models.Workspace{}).Where("id=? AND audit_status <> ?", id, status).Updates(map[string]interface{}{"audit_status": status, "audit_reason": reason, "audited_by": admin, "audited_at": now})
			if r.Error != nil {
				return r.Error
			}
			if r.RowsAffected == 0 {
				return nil
			}
		} else {
			var d models.Doc
			if err := tx.First(&d, id).Error; err != nil {
				return err
			}
			title = d.Title
			owner = d.OwnerID
			r := tx.Model(&models.Doc{}).Where("id=? AND audit_status <> ?", id, status).Updates(map[string]interface{}{"audit_status": status, "audit_reason": reason, "audited_by": admin, "audited_at": now})
			if r.Error != nil {
				return r.Error
			}
			if r.RowsAffected == 0 {
				return nil
			}
		}
		action := "unblock"
		changed = true
		if status == 1 {
			action = "block"
		}
		return tx.Create(&models.AdminAuditLog{AdminID: admin, AdminUsername: username, Action: action, TargetType: kind, TargetID: id, TargetTitle: title, OwnerID: owner, Reason: reason, IP: ip}).Error
	})
	if err == nil && changed {
		deleteWorkspaceCache(id)
		go func() {
			var notificationErr error
			if kind == "doc" {
				notificationErr = NewNotificationService().CreateDocAuditNotification(id, status == 1, reason)
			} else {
				notificationErr = NewNotificationService().CreateWorkspaceAuditNotification(id, status == 1, reason)
			}
			if notificationErr != nil {
				zap.L().Error("knowledge audit notification failed", zap.Error(notificationErr))
			}
		}()
	}
	return err
}
func (s *KnowledgeAdminService) Delete(kind string, id, admin uint, reason, username, ip string) error {
	kind = strings.TrimSuffix(kind, "s")
	if kind != "workspace" && kind != "doc" && kind != "share" {
		return ErrKnowledgeInvalid
	}
	var workspaceID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var title string
		var owner uint
		if kind == "share" {
			var link models.ShareLink
			if err := tx.First(&link, id).Error; err != nil {
				return err
			}
			if err := tx.Delete(&link).Error; err != nil {
				return err
			}
			return NewAdminAuditService().Log(tx, &models.AdminAuditLog{AdminID: admin, AdminUsername: username, Action: "delete", TargetType: "share_link", TargetID: id, TargetTitle: link.Token, OwnerID: link.OwnerID, Reason: reason, IP: ip})
		}
		if kind == "workspace" {
			var w models.Workspace
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&w, id).Error; err != nil {
				return err
			}
			title = w.Name
			workspaceID = w.ID
			owner = w.OwnerID
			var docs []uint
			if err := tx.Model(&models.Doc{}).Where("workspace_id=?", id).Pluck("id", &docs).Error; err != nil {
				return err
			}
			if len(docs) > 0 {
				if err := tx.Where("doc_id IN ?", docs).Delete(&models.DocVersion{}).Error; err != nil {
					return err
				}
				if err := tx.Where("doc_id IN ?", docs).Delete(&models.ShareLink{}).Error; err != nil {
					return err
				}
				if err := tx.Where("workspace_id=?", id).Delete(&models.Doc{}).Error; err != nil {
					return err
				}
			}
			if err := tx.Where("workspace_id=?", id).Delete(&models.Catalog{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&w).Error; err != nil {
				return err
			}
			if err := tx.Where("workspace_id=?", id).Delete(&models.WorkspaceMember{}).Error; err != nil {
				return err
			}
		} else {
			var d models.Doc
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&d, id).Error; err != nil {
				return err
			}
			title = d.Title
			workspaceID = d.WorkspaceID
			owner = d.OwnerID
			if err := tx.Where("doc_id=?", id).Delete(&models.DocVersion{}).Error; err != nil {
				return err
			}
			if err := tx.Where("doc_id=?", id).Delete(&models.ShareLink{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&d).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Workspace{}).Where("id=?", d.WorkspaceID).UpdateColumn("doc_count", gorm.Expr("GREATEST(doc_count-1,0)")).Error; err != nil {
				return err
			}
		}
		return tx.Create(&models.AdminAuditLog{AdminID: admin, AdminUsername: username, Action: "delete", TargetType: kind, TargetID: id, TargetTitle: title, OwnerID: owner, Reason: reason, IP: ip}).Error
	})
	if err == nil && workspaceID != 0 {
		deleteWorkspaceCache(workspaceID)
	}
	return err
}
func (s *KnowledgeAdminService) Quota() map[string]int {
	out := map[string]int{QuotaMaxWorkspaces: 10, QuotaMaxDocs: 500, QuotaMaxStorage: 1024}
	var rows []models.Setting
	if err := database.DB.Where("`key` IN ?", []string{QuotaMaxWorkspaces, QuotaMaxDocs, QuotaMaxStorage}).Find(&rows).Error; err != nil {
		zap.L().Error("read knowledge quotas failed", zap.Error(err))
	}
	for _, r := range rows {
		if n, e := strconv.Atoi(r.Value); e == nil && n >= 0 {
			out[r.Key] = n
		} else {
			out[r.Key] = 0
			zap.L().Warn("invalid knowledge quota", zap.String("key", r.Key))
		}
	}
	return out
}
func (s *KnowledgeAdminService) SaveQuota(v map[string]int) error {
	for k, n := range v {
		if (k != QuotaMaxWorkspaces && k != QuotaMaxDocs && k != QuotaMaxStorage) || n < 0 {
			return ErrKnowledgeInvalid
		}
	}
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for k, n := range v {
			if n < 0 {
				return errors.New("配额不能为负数")
			}
			if err := tx.Unscoped().Where("`key`=?", k).Assign(map[string]interface{}{"key": k, "value": strconv.Itoa(n), "type": "int", "group": "knowledge", "is_public": false, "deleted_at": nil}).FirstOrCreate(&models.Setting{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
