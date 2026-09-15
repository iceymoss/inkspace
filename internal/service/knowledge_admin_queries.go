package service

import (
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"time"
)

// AdminDocMetadata deliberately has no content or attachment storage fields.
type AdminDocMetadata struct {
	OwnerName   string    `json:"owner_name"`
	ID          uint      `json:"id"`
	WorkspaceID uint      `json:"workspace_id"`
	OwnerID     uint      `json:"owner_id"`
	Title       string    `json:"title"`
	Kind        string    `json:"kind"`
	Status      int       `json:"status"`
	AuditStatus int8      `json:"audit_status"`
	AuditReason string    `json:"audit_reason"`
	WordCount   int       `json:"word_count"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPublic    bool      `json:"is_public"`
}
type AdminDocDetail struct {
	AdminDocMetadata
	Content     *string `json:"content,omitempty"`
	ContentHTML *string `json:"content_html,omitempty"`
}

const adminDocColumns = "docs.id, docs.workspace_id, docs.owner_id, docs.title, docs.kind, docs.status, docs.audit_status, docs.audit_reason, docs.word_count, docs.view_count, docs.created_at, docs.updated_at, workspaces.is_public"

func (s *KnowledgeAdminService) DocDetail(id uint) (*AdminDocDetail, error) {
	var result AdminDocDetail
	q := database.DB.Model(&models.Doc{}).Joins("JOIN workspaces ON workspaces.id = docs.workspace_id AND workspaces.deleted_at IS NULL")
	if err := q.Select(adminDocColumns).Where("docs.id = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}
	// Read body only through a public-only query, including ownership and moderation.
	if result.IsPublic && result.Status == models.DocStatusPublished && result.AuditStatus == 0 {
		var body struct {
			Content     string
			ContentHTML string
		}
		r := database.DB.Model(&models.Doc{}).Select("docs.content, docs.content_html").Joins("JOIN workspaces ON workspaces.id = docs.workspace_id AND workspaces.owner_id = docs.owner_id AND workspaces.deleted_at IS NULL").Where("docs.id = ? AND docs.status = 1 AND docs.audit_status = 0 AND workspaces.is_public = true AND workspaces.audit_status = 0", id).Scan(&body)
		if r.Error != nil {
			return nil, r.Error
		}
		if r.RowsAffected > 0 {
			result.Content = &body.Content
			html := sanitizePublicWikiHTML(body.ContentHTML)
			result.ContentHTML = &html
		}
	}
	return &result, nil
}
func (s *KnowledgeAdminService) WorkspaceDetail(id uint) (interface{}, error) {
	var w models.Workspace
	if err := database.DB.First(&w, id).Error; err != nil {
		return nil, err
	}
	var docs, catalogs int64
	if err := database.DB.Model(&models.Doc{}).Where("workspace_id=?", id).Count(&docs).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Catalog{}).Where("workspace_id=?", id).Count(&catalogs).Error; err != nil {
		return nil, err
	}
	return struct {
		models.Workspace
		DocumentCount int64 `json:"document_count"`
		CatalogCount  int64 `json:"catalog_count"`
	}{w, docs, catalogs}, nil
}

type AdminShareLink struct {
	models.ShareLink
	OwnerName string `json:"owner_name"`
}

// table is supplied only by the fixed internal callers, never request input.
func adminOwnerNameSQL(table string) string {
	return "COALESCE((SELECT COALESCE(NULLIF(users.nickname, ''), users.username) FROM users WHERE users.id = " + table + ".owner_id AND users.deleted_at IS NULL), '用户已删除') AS owner_name"
}

func (s *KnowledgeAdminService) Shares(page, size int, owner, doc *uint, state string) ([]AdminShareLink, int64, error) {
	rows := make([]AdminShareLink, 0)
	var total int64
	q := database.DB.Model(&models.ShareLink{})
	if owner != nil {
		q = q.Where("owner_id=?", *owner)
	}
	if doc != nil {
		q = q.Where("doc_id=?", *doc)
	}
	switch state {
	case "active":
		q = q.Where("enabled = true AND (expires_at IS NULL OR expires_at > ?)", time.Now())
	case "disabled":
		q = q.Where("enabled = false")
	case "expired":
		q = q.Where("enabled = true AND expires_at <= ?", time.Now())
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := pageQuery(q.Select("share_links.*, "+adminOwnerNameSQL("share_links")), page, size).Find(&rows).Error
	return rows, total, err
}
func (s *KnowledgeAdminService) AuditLogs(page, size int, admin *uint, target, action string, start, end *time.Time) ([]models.AdminAuditLog, int64, error) {
	rows := make([]models.AdminAuditLog, 0)
	var total int64
	q := database.DB.Model(&models.AdminAuditLog{})
	if admin != nil {
		q = q.Where("admin_id=?", *admin)
	}
	if target != "" {
		q = q.Where("target_type=?", target)
	}
	if action != "" {
		q = q.Where("action=?", action)
	}
	if start != nil {
		q = q.Where("created_at >= ?", *start)
	}
	if end != nil {
		q = q.Where("created_at <= ?", *end)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := pageQuery(q, page, size).Find(&rows).Error
	return rows, total, err
}

type KnowledgeUsage struct {
	ID               uint   `json:"id"`
	Username         string `json:"username"`
	WorkspaceCount   int64  `json:"workspace_count"`
	DocCount         int64  `json:"doc_count"`
	StorageBytes     int64  `json:"storage_bytes"`
	MaxWorkspaceDocs int64  `json:"max_workspace_docs"`
	OverQuota        bool   `json:"over_quota"`
}

func (s *KnowledgeAdminService) Usage(page, size int, sort string) ([]KnowledgeUsage, int64, error) {
	rows := make([]KnowledgeUsage, 0)
	var total int64
	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	w := database.DB.Model(&models.Workspace{}).Select("owner_id, COUNT(*) AS workspace_count").Group("owner_id")
	d := database.DB.Model(&models.Doc{}).Select("owner_id, COUNT(*) AS doc_count").Group("owner_id")
	storage := database.DB.Model(&models.Doc{}).Select("docs.owner_id, COALESCE(SUM(attachments.file_size),0) AS storage_bytes").Joins("JOIN attachments ON attachments.id = docs.attachment_id AND attachments.deleted_at IS NULL").Group("docs.owner_id")
	perWorkspace := database.DB.Model(&models.Doc{}).Select("owner_id, workspace_id, COUNT(*) AS count").Group("owner_id, workspace_id")
	maximum := database.DB.Table("(?) AS counts", perWorkspace).Select("owner_id, MAX(count) AS max_workspace_docs").Group("owner_id")
	orders := map[string]string{"workspace_count": "workspace_count DESC", "doc_count": "doc_count DESC", "storage_bytes": "storage_bytes DESC"}
	order, ok := orders[sort]
	if !ok {
		order = "storage_bytes DESC"
	}
	err := database.DB.Model(&models.User{}).Select("users.id, users.username, COALESCE(w.workspace_count,0) AS workspace_count, COALESCE(d.doc_count,0) AS doc_count, COALESCE(a.storage_bytes,0) AS storage_bytes, COALESCE(m.max_workspace_docs,0) AS max_workspace_docs").Joins("LEFT JOIN (?) AS w ON w.owner_id=users.id", w).Joins("LEFT JOIN (?) AS d ON d.owner_id=users.id", d).Joins("LEFT JOIN (?) AS a ON a.owner_id=users.id", storage).Joins("LEFT JOIN (?) AS m ON m.owner_id=users.id", maximum).Order(order).Order("users.id DESC").Offset((page - 1) * size).Limit(size).Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	quota := s.Quota()
	for i := range rows {
		r := &rows[i]
		r.OverQuota = (quota[QuotaMaxWorkspaces] > 0 && r.WorkspaceCount > int64(quota[QuotaMaxWorkspaces])) || (quota[QuotaMaxDocs] > 0 && r.MaxWorkspaceDocs > int64(quota[QuotaMaxDocs])) || (quota[QuotaMaxStorage] > 0 && float64(r.StorageBytes) > float64(quota[QuotaMaxStorage])*1048576)
	}
	return rows, total, nil
}

type KnowledgeTrend struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

func fillKnowledgeTrend(now time.Time, rows []KnowledgeTrend) []KnowledgeTrend {
	counts := map[string]int64{}
	for _, r := range rows {
		counts[r.Day] = r.Count
	}
	out := make([]KnowledgeTrend, 7)
	for i := range out {
		day := now.AddDate(0, 0, i-6).Format("2006-01-02")
		out[i] = KnowledgeTrend{day, counts[day]}
	}
	return out
}
func (s *KnowledgeAdminService) Overview() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for _, item := range []struct {
		key       string
		model     interface{}
		condition string
	}{{"workspaces", &models.Workspace{}, "1=1"}, {"public_workspaces", &models.Workspace{}, "is_public=true"}, {"private_workspaces", &models.Workspace{}, "is_public=false"}, {"docs", &models.Doc{}, "1=1"}, {"draft_docs", &models.Doc{}, "status=0 AND audit_status=0"}, {"published_docs", &models.Doc{}, "status=1 AND audit_status=0"}, {"blocked_docs", &models.Doc{}, "audit_status=1"}, {"shares", &models.ShareLink{}, "1=1"}} {
		var n int64
		if err := database.DB.Model(item.model).Where(item.condition).Count(&n).Error; err != nil {
			return nil, err
		}
		result[item.key] = n
	}
	var bytes int64
	if err := database.DB.Model(&models.Doc{}).Select("COALESCE(SUM(attachments.file_size),0)").Joins("JOIN attachments ON attachments.id=docs.attachment_id AND attachments.deleted_at IS NULL").Scan(&bytes).Error; err != nil {
		return nil, err
	}
	result["storage_bytes"] = bytes
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)
	var trend []KnowledgeTrend
	if err := database.DB.Model(&models.Doc{}).Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS day, COUNT(*) AS count").Where("created_at>=?", start).Group("day").Order("day").Scan(&trend).Error; err != nil {
		return nil, err
	}
	result["trend"] = fillKnowledgeTrend(now, trend)
	top := make([]AdminDocMetadata, 0)
	if err := database.DB.Model(&models.Doc{}).Select(adminDocColumns).Joins("JOIN workspaces ON workspaces.id=docs.workspace_id AND workspaces.owner_id=docs.owner_id AND workspaces.deleted_at IS NULL").Where("docs.status=1 AND docs.audit_status=0 AND workspaces.audit_status=0 AND workspaces.is_public=true").Order("docs.view_count DESC, docs.id DESC").Limit(10).Scan(&top).Error; err != nil {
		return nil, err
	}
	result["top_docs"] = top
	return result, nil
}
