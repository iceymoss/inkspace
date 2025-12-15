package handler

import (
	"strconv"

	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdHandler struct {
	service *service.AdService
}

func NewAdHandler() *AdHandler {
	return &AdHandler{
		service: service.NewAdService(),
	}
}

// ========== 广告位置管理 ==========

// CreatePosition 创建广告位置
// POST /admin/ad-positions
func (h *AdHandler) CreatePosition(c *gin.Context) {
	var req models.AdPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	position, err := h.service.CreatePosition(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, position.ToResponse())
}

// GetPositionList 获取广告位置列表
// GET /admin/ad-positions
func (h *AdHandler) GetPositionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.service.GetPositionList(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var responses []*models.AdPositionResponse
	for _, p := range list {
		responses = append(responses, p.ToResponse())
	}

	utils.PageResponse(c, responses, total, page, pageSize)
}

// GetPositionByID 获取广告位置详情
// GET /admin/ad-positions/:id
func (h *AdHandler) GetPositionByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	position, err := h.service.GetPositionByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "广告位置不存在")
		return
	}

	utils.Success(c, position.ToResponse())
}

// UpdatePosition 更新广告位置
// PUT /admin/ad-positions/:id
func (h *AdHandler) UpdatePosition(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.AdPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	position, err := h.service.UpdatePosition(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, position.ToResponse())
}

// DeletePosition 删除广告位置
// DELETE /admin/ad-positions/:id
func (h *AdHandler) DeletePosition(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeletePosition(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ========== 广告消息管理 ==========

// CreateAdvertisement 创建广告消息
// POST /admin/advertisements
func (h *AdHandler) CreateAdvertisement(c *gin.Context) {
	var req models.AdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ad, err := h.service.CreateAdvertisement(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, ad.ToResponse())
}

// GetAdvertisementList 获取广告消息列表
// GET /admin/advertisements
func (h *AdHandler) GetAdvertisementList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.service.GetAdvertisementList(page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var responses []*models.AdvertisementResponse
	for _, ad := range list {
		responses = append(responses, ad.ToResponse())
	}

	utils.PageResponse(c, responses, total, page, pageSize)
}

// GetAdvertisementByID 获取广告消息详情
// GET /admin/advertisements/:id
func (h *AdHandler) GetAdvertisementByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	ad, err := h.service.GetAdvertisementByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "广告消息不存在")
		return
	}

	utils.Success(c, ad.ToResponse())
}

// UpdateAdvertisement 更新广告消息
// PUT /admin/advertisements/:id
func (h *AdHandler) UpdateAdvertisement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.AdvertisementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	ad, err := h.service.UpdateAdvertisement(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, ad.ToResponse())
}

// DeleteAdvertisement 删除广告消息
// DELETE /admin/advertisements/:id
func (h *AdHandler) DeleteAdvertisement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeleteAdvertisement(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ========== 广告投放管理 ==========

// CreatePlacement 创建广告投放
// POST /admin/ad-placements
func (h *AdHandler) CreatePlacement(c *gin.Context) {
	var req models.AdPlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	placement, err := h.service.CreatePlacement(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	placement, _ = h.service.GetPlacementByID(placement.ID)
	utils.Success(c, placement.ToResponse())
}

// GetPlacementList 获取广告投放列表
// GET /admin/ad-placements
func (h *AdHandler) GetPlacementList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	positionIDStr := c.Query("position_id")

	var positionID *uint
	if positionIDStr != "" {
		id, _ := strconv.ParseUint(positionIDStr, 10, 32)
		positionID = new(uint)
		*positionID = uint(id)
	}

	list, total, err := h.service.GetPlacementList(page, pageSize, positionID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	var responses []*models.AdPlacementResponse
	for _, p := range list {
		responses = append(responses, p.ToResponse())
	}

	utils.PageResponse(c, responses, total, page, pageSize)
}

// GetPlacementByID 获取广告投放详情
// GET /admin/ad-placements/:id
func (h *AdHandler) GetPlacementByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	placement, err := h.service.GetPlacementByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "广告投放不存在")
		return
	}

	utils.Success(c, placement.ToResponse())
}

// UpdatePlacement 更新广告投放
// PUT /admin/ad-placements/:id
func (h *AdHandler) UpdatePlacement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.AdPlacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	placement, err := h.service.UpdatePlacement(uint(id), &req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	placement, _ = h.service.GetPlacementByID(placement.ID)
	utils.Success(c, placement.ToResponse())
}

// DeletePlacement 删除广告投放
// DELETE /admin/ad-placements/:id
func (h *AdHandler) DeletePlacement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.DeletePlacement(uint(id)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ========== 前端获取广告 ==========

// GetAdsByPositionCode 根据位置代码获取广告
// GET /api/ads?code=xxx
func (h *AdHandler) GetAdsByPositionCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.BadRequest(c, "位置代码不能为空")
		return
	}

	ads, err := h.service.GetAdsByPositionCode(code)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, ads)
}

// RecordAdClick 记录广告点击
// POST /api/ads/:id/click
func (h *AdHandler) RecordAdClick(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.RecordAdClick(uint(id)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "记录成功", nil)
}

// RecordAdView 记录广告展示
// POST /api/ads/:id/view
func (h *AdHandler) RecordAdView(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.RecordAdView(uint(id)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "记录成功", nil)
}
