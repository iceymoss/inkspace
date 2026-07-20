package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type ShareHandler struct{ service *service.ShareService }

func NewShareHandler() *ShareHandler { return &ShareHandler{service: service.NewShareService()} }

func (h *ShareHandler) Create(c *gin.Context) {
	docID, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	var req models.ShareLinkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	link, err := h.service.Create(docID, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, link.ToResponse())
}

func (h *ShareHandler) List(c *gin.Context) {
	docID, userID, ok := docAndUser(c)
	if !ok {
		return
	}
	links, err := h.service.List(docID, userID)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	responses := make([]*models.ShareLinkResponse, len(links))
	for i := range links {
		responses[i] = links[i].ToResponse()
	}
	utils.Success(c, responses)
}

func (h *ShareHandler) Update(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req models.ShareLinkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	link, err := h.service.Update(id, userID, &req)
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, link.ToResponse())
}

func (h *ShareHandler) Delete(c *gin.Context) {
	id, ok := pathUint(c, "id")
	if !ok {
		return
	}
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		knowledgeError(c, err)
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ShareHandler) Public(c *gin.Context) {
	doc, err := h.service.Public(c.Param("token"), time.Now())
	if err != nil {
		knowledgeError(c, err)
		return
	}
	utils.Success(c, doc.ToResponse())
}
