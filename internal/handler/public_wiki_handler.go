package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type PublicWikiHandler struct{ service *service.PublicWikiService }

func NewPublicWikiHandler() *PublicWikiHandler {
	return &PublicWikiHandler{service: service.NewPublicWikiService()}
}

func (h *PublicWikiHandler) Workspaces(c *gin.Context) {
	page, ok := publicWikiPageParam(c, "page", 1, 0)
	if !ok {
		return
	}
	pageSize, ok := publicWikiPageParam(c, "page_size", 20, 100)
	if !ok {
		return
	}
	workspaces, total, err := h.service.Workspaces(page, pageSize)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.PageResponse(c, workspaces, total, page, pageSize)
}

func (h *PublicWikiHandler) Tree(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	tree, err := h.service.Tree(id)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, tree)
}

func (h *PublicWikiHandler) Doc(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	doc, err := h.service.Doc(id)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc)
}

func publicWikiPageParam(c *gin.Context, name string, fallback, maximum int) (int, bool) {
	raw := c.Query(name)
	if raw == "" {
		return fallback, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 || maximum > 0 && value > maximum {
		utils.BadRequest(c, "无效的分页参数")
		return 0, false
	}
	return value, true
}
