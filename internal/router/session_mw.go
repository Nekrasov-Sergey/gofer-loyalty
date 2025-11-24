package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/gofer-loyalty/internal/service"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/errcodes"
	"github.com/Nekrasov-Sergey/gofer-loyalty/pkg/logger"
)

const ContextUserID = "user_id"

// SessionMiddleware проверяет сессию пользователя
func SessionMiddleware(repo service.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session")
		if err != nil {
			logger.RespondError(c, errcodes.ErrUserUnauthorized, http.StatusUnauthorized)
			return
		}

		userID, err := repo.GetUserIDByTokenSession(c.Request.Context(), sessionToken)
		if err != nil {
			if errors.Is(err, errcodes.ErrUserUnauthorized) {
				logger.RespondError(c, err, http.StatusUnauthorized)
				return
			}
			logger.RespondError(c, err, http.StatusInternalServerError)
			return
		}

		c.Set(ContextUserID, userID)

		c.Next()
	}
}
