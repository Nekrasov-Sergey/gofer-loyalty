package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/router"
	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) login(c *gin.Context) {
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

	sessionToken, err := h.service.Login(ctx, user)
	if err != nil {
		if errors.Is(err, errcodes.ErrInvalidCredentials) {
			logger.RespondError(c, err, http.StatusUnauthorized)
			return
		}
		logger.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	h.setCookie(c, sessionToken)

	c.Status(http.StatusOK)
}

func (h *Handler) setCookie(c *gin.Context, sessionToken string) {
	c.SetCookie(
		router.CookieSession,
		sessionToken,
		int(time.Duration(h.config.SessionTTL).Seconds()),
		"/api/user",
		"",
		false,
		true,
	)
}
