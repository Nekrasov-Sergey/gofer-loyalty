package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/types"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

func (h *Handler) register(c *gin.Context) {
	ctx := c.Request.Context()

	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.RespondError(c, errors.Wrap(err, "не удалось распарсить тело запроса"), http.StatusBadRequest)
		return
	}

	if user.Login == "" {
		logger.RespondError(c, errors.New("отсутствует логин"), http.StatusBadRequest)
		return
	}

	if user.Password == "" {
		logger.RespondError(c, errors.New("отсутствует пароль"), http.StatusBadRequest)
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
