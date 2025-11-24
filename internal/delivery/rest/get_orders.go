package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) getOrders(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(router.ContextUserID)

	orders, err := h.service.GetOrders(ctx, userID)
	if err != nil {
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}
