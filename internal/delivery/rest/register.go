package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) register(c *gin.Context) {
	ctx := c.Request.Context()

	user := &types.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.RespondError(c, multierr.Append(errcodes.ErrInvalidRequestFormat, err), http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		logger.RespondError(c, errcodes.ErrLoginIsMissing, http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		logger.RespondError(c, errcodes.ErrPasswordIsMissing, http.StatusBadRequest)
		return
	}

	sessionToken, err := h.service.Register(ctx, user)
	if err != nil {
		if errors.Is(err, errcodes.ErrLoginAlreadyExists) {
			logger.RespondError(c, err, http.StatusConflict)
			return
		}
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	h.setCookie(c, sessionToken)

	c.Status(http.StatusOK)
}
