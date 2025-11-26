package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) getBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(router.ContextUserID)

	balance, err := h.service.GetBalance(ctx, userID)
	if err != nil {
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance)
}
