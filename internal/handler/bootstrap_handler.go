package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/service"
	"github.com/iceymoss/inkspace/internal/utils"
)

type BootstrapHandler struct {
	service *service.BootstrapService
}

func NewBootstrapHandler() *BootstrapHandler {
	return &BootstrapHandler{service: service.NewBootstrapService()}
}

func (h *BootstrapHandler) Status(c *gin.Context) {
	status, err := h.service.Status()
	if err != nil {
		utils.InternalServerError(c, "无法读取系统初始化状态，请联系维护者")
		return
	}
	utils.Success(c, status)
}

func (h *BootstrapHandler) Create(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16<<10)
	var req models.AdminBootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请检查密钥、用户名、邮箱、密码及确认密码")
		return
	}
	if err := h.service.Create(&req); err != nil {
		switch {
		case errors.Is(err, service.ErrBootstrapSecret):
			utils.Forbidden(c, err.Error())
		case errors.Is(err, service.ErrBootstrapConflict), errors.Is(err, service.ErrBootstrapCompleted):
			utils.BadRequest(c, err.Error())
		case errors.Is(err, service.ErrBootstrapUnavailable):
			utils.InternalServerError(c, err.Error())
		default:
			utils.BadRequest(c, "请检查初始化信息")
		}
		return
	}
	utils.SuccessWithMessage(c, "首位管理员创建成功，请登录", nil)
}
