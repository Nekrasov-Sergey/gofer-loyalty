package router

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type compressWriter struct {
	gin.ResponseWriter
	zw           *gzip.Writer
	isCompressed bool
}

func newCompressWriter(rw gin.ResponseWriter) *compressWriter {
	return &compressWriter{ResponseWriter: rw}
}

func (w *compressWriter) Write(data []byte) (int, error) {
	ct := w.ResponseWriter.Header().Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "text/html") {
		if !w.isCompressed {
			w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
			w.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
			w.ResponseWriter.Header().Del("Content-Length")

			w.zw = gzip.NewWriter(w.ResponseWriter)
			w.isCompressed = true
		}
		return w.zw.Write(data)
	}
	return w.ResponseWriter.Write(data)
}

func (w *compressWriter) Close() error {
	if w.zw != nil {
		return w.zw.Close()
	}
	return nil
}

// CompressMiddleware добавляет сжатие и декомпрессию данных в gin
func CompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			cw := newCompressWriter(c.Writer)
			defer func(cw *compressWriter) {
				if err := cw.Close(); err != nil {
					log.Error().Err(err).Msg("Не удалось закрыть compressWriter")
				}
			}(cw)
			c.Writer = cw
		}

		c.Next()
	}
}
