package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func New(mode string, logger zerolog.Logger) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Recovery(), LoggerMiddleware(logger))
	return r
}
