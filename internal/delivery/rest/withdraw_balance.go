package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/utils"
)

func (h *Handler) withdrawBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt64(router.ContextUserID)

	withdrawal := &types.Withdrawal{}
	if err := c.ShouldBindJSON(&withdrawal); err != nil {
		logger.RespondError(c, multierr.Append(errcodes.ErrInvalidRequestFormat, err), http.StatusBadRequest)
		return
	}

	if withdrawal.Withdrawn.InexactFloat64() < 0 {
		logger.RespondError(c, errcodes.ErrNegativeWithdrawal, http.StatusBadRequest)
	}

	withdrawal.UserID = userID
	withdrawal.ProcessedAt = time.Now()

	orderNumber, err := strconv.ParseInt(withdrawal.OrderNumber, 10, 64)
	if err != nil {
		logger.RespondError(c, errcodes.ErrOrderNumberMustBeNumeric, http.StatusBadRequest)
		return
	}

	if !utils.IsValidLuhn(orderNumber) {
		logger.RespondError(c, errcodes.ErrInvalidOrderNumberFormat, http.StatusUnprocessableEntity)
		return
	}

	if err := h.service.WithdrawBalance(ctx, withdrawal); err != nil {
		if errors.Is(err, errcodes.ErrNotEnoughBalance) {
			logger.RespondError(c, err, http.StatusPaymentRequired)
			return
		}
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
