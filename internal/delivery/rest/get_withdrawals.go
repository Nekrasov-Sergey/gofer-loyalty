package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) getWithdrawals(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(router.ContextUserID)

	withdrawals, err := h.service.GetWithdrawals(ctx, userID)
	if err != nil {
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
