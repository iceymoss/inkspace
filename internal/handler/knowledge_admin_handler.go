package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
	"strconv"
	"strings"
)

type KnowledgeAdminHandler struct {
	s *service.KnowledgeAdminService
}

func NewKnowledgeAdminHandler() *KnowledgeAdminHandler {
	return &KnowledgeAdminHandler{service.NewKnowledgeAdminService()}
}
func paging(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	z, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if p < 1 {
		p = 1
	}
	if z < 1 || z > 100 {
		z = 10
	}
	return p, z
}
func (h *KnowledgeAdminHandler) ListW(c *gin.Context) {
	p, z := paging(c)
	var q knowledgeAdminQuery
	if c.ShouldBindQuery(&q) != nil {
		utils.BadRequest(c, "筛选参数无效")
		return
	}
	var audit *bool
	if q.AuditStatus != nil {
		v := *q.AuditStatus == 1
		audit = &v
	}
	v, n, e := h.s.Workspaces(p, z, c.Query("keyword"), q.OwnerID, q.IsPublic, audit)
	if e != nil {
		utils.InternalServerError(c, e.Error())
		return
	}
	utils.PageResponse(c, v, n, p, z)
}
func (h *KnowledgeAdminHandler) ListD(c *gin.Context) {
	p, z := paging(c)
	var q knowledgeAdminQuery
	if c.ShouldBindQuery(&q) != nil {
		utils.BadRequest(c, "筛选参数无效")
		return
	}
	v, n, e := h.s.Docs(p, z, c.Query("keyword"), q.OwnerID, q.WorkspaceID, q.Status, q.AuditStatus, c.Query("kind"))
	if e != nil {
		utils.InternalServerError(c, e.Error())
		return
	}
	utils.PageResponse(c, v, n, p, z)
}
func (h *KnowledgeAdminHandler) Audit(c *gin.Context) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 32)
	if e != nil {
		utils.BadRequest(c, "无效ID")
		return
	}
	var r struct {
		AuditStatus *int8  `json:"audit_status" binding:"required,oneof=0 1"`
		Reason      string `json:"reason" binding:"required,min=2,max=200"`
	}
	if c.ShouldBindJSON(&r) != nil || len([]rune(strings.TrimSpace(r.Reason))) < 2 {
		utils.BadRequest(c, "处置理由不能为空")
		return
	}
	uid, _ := c.Get("user_id")
	admin, _ := uid.(uint)
	if e = h.s.Audit(c.Param("type"), uint(id), admin, *r.AuditStatus, r.Reason, c.GetString("username"), c.ClientIP()); e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.Success(c, nil)
}
func (h *KnowledgeAdminHandler) Delete(c *gin.Context) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 32)
	if e != nil {
		utils.BadRequest(c, "无效ID")
		return
	}
	var r struct {
		Reason string `json:"reason" binding:"required,min=2,max=200"`
	}
	if c.ShouldBindJSON(&r) != nil || len([]rune(strings.TrimSpace(r.Reason))) < 2 {
		utils.BadRequest(c, "删除理由不能为空")
		return
	}
	uid, _ := c.Get("user_id")
	admin, _ := uid.(uint)
	if e = h.s.Delete(c.Param("type"), uint(id), admin, r.Reason, c.GetString("username"), c.ClientIP()); e != nil {
		adminKnowledgeError(c, e)
		return
	}
	utils.Success(c, nil)
}
func (h *KnowledgeAdminHandler) Quota(c *gin.Context) { utils.Success(c, h.s.Quota()) }
func (h *KnowledgeAdminHandler) SaveQuota(c *gin.Context) {
	var v map[string]int
	if c.ShouldBindJSON(&v) != nil {
		utils.BadRequest(c, "配额格式无效")
		return
	}
	if e := h.s.SaveQuota(v); e != nil {
		utils.BadRequest(c, e.Error())
		return
	}
	utils.Success(c, h.s.Quota())
}
