package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service *service.TagService
}

func NewTagHandler() *TagHandler {
	return &TagHandler{
		service: service.NewTagService(),
	}
}

func (h *TagHandler) Create(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 如果是用户创建标签，自动设置user_id
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		req.UserID = &uid
	}

	tag, err := h.service.Create(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, tag.ToResponse())
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tag, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, tag.ToResponse())
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *TagHandler) GetList(c *gin.Context) {
	tags, err := h.service.GetList()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	tagResponses := make([]*models.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = tag.ToResponse()
	}

	utils.Success(c, tagResponses)
}
