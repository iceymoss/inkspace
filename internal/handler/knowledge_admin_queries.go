package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
	"gorm.io/gorm"
	"time"
)

type knowledgeAdminQuery struct {
	OwnerID     *uint      `form:"owner_id"`
	WorkspaceID *uint      `form:"workspace_id"`
	DocID       *uint      `form:"doc_id"`
	AdminID     *uint      `form:"admin_id"`
	Status      *int       `form:"status" binding:"omitempty,oneof=0 1"`
	AuditStatus *int       `form:"audit_status" binding:"omitempty,oneof=0 1"`
	IsPublic    *bool      `form:"is_public"`
	State       string     `form:"state" binding:"omitempty,oneof=active disabled expired"`
	TargetType  string     `form:"target_type" binding:"omitempty,oneof=workspace doc share_link"`
	Action      string     `form:"action" binding:"omitempty,oneof=block unblock delete"`
	Start       *time.Time `form:"start" time_format:"2006-01-02"`
	End         *time.Time `form:"end" time_format:"2006-01-02"`
}

func adminKnowledgeError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrKnowledgeInvalid) {
		utils.BadRequest(c, err.Error())
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.NotFound(c, "资源不存在")
	} else {
		utils.InternalServerError(c, "知识库操作失败")
	}
}
func (h *KnowledgeAdminHandler) Detail(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	var data interface{}
	var err error
	if c.Param("type") == "workspaces" {
		data, err = h.s.WorkspaceDetail(id)
	} else if c.Param("type") == "docs" {
		data, err = h.s.DocDetail(id)
	} else {
		utils.NotFound(c, "资源不存在")
		return
	}
	if err != nil {
		adminKnowledgeError(c, err)
		return
	}
	utils.Success(c, data)
}
func (h *KnowledgeAdminHandler) Shares(c *gin.Context) {
	var q knowledgeAdminQuery
	if c.ShouldBindQuery(&q) != nil {
		utils.BadRequest(c, "筛选参数无效")
		return
	}
	p, z := paging(c)
	v, n, e := h.s.Shares(p, z, q.OwnerID, q.DocID, q.State)
	if e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.PageResponse(c, v, n, p, z)
}
func (h *KnowledgeAdminHandler) AuditLogs(c *gin.Context) {
	var q knowledgeAdminQuery
	if c.ShouldBindQuery(&q) != nil {
		utils.BadRequest(c, "筛选参数无效")
		return
	}
	if q.End != nil {
		end := q.End.AddDate(0, 0, 1).Add(-time.Nanosecond)
		q.End = &end
	}
	p, z := paging(c)
	v, n, e := h.s.AuditLogs(p, z, q.AdminID, q.TargetType, q.Action, q.Start, q.End)
	if e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.PageResponse(c, v, n, p, z)
}
func (h *KnowledgeAdminHandler) Usage(c *gin.Context) {
	p, z := paging(c)
	v, n, e := h.s.Usage(p, z, c.Query("sort"))
	if e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.PageResponse(c, v, n, p, z)
}
func (h *KnowledgeAdminHandler) Overview(c *gin.Context) {
	v, e := h.s.Overview()
	if e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.Success(c, v)
}
