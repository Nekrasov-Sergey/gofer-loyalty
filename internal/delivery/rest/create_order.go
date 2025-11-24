package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/utils"
)

func (h *Handler) createOrder(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(router.ContextUserID)

	var orderNumber int64
	if err := c.ShouldBindJSON(&orderNumber); err != nil {
		logger.RespondError(c, errors.Wrap(err, "неверный формат запроса"), http.StatusBadRequest)
		return
	}

	if !utils.IsValidLuhn(orderNumber) {
		logger.RespondError(c, errors.New("неверный формат номера заказа"), http.StatusUnprocessableEntity)
		return
	}

	err := h.service.CreateOrder(ctx, types.Order{
		Number:     strconv.FormatInt(orderNumber, 10),
		UserID:     userID,
		UploadedAt: time.Now(),
	})
	if err != nil {
		if errors.Is(err, errcodes.ErrOrderAlreadyUploadedByUser) {
			c.Status(http.StatusOK)
			return
		}
		if errors.Is(err, errcodes.ErrOrderAlreadyUploadedByAnotherUser) {
			logger.RespondError(c, err, http.StatusConflict)
			return
		}
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
