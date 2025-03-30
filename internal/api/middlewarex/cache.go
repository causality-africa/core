package middlewarex

import (
	"bytes"
	"core/internal/cache"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func cacheKey(c echo.Context) string {
	hasher := sha256.New()
	hasher.Write([]byte(c.Request().URL.String()))
	return fmt.Sprintf("core:responses:%s", hex.EncodeToString(hasher.Sum(nil)))
}

func CacheMiddleware(cache *cache.Cache, ttl time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			// Try to get from cache
			key := cacheKey(c)
			cachedResp, err := cache.Get(c.Request().Context(), key)
			if err == nil && cachedResp != nil {
				cacheData, ok := cachedResp.(map[string]any)
				if ok {
					status := int(cacheData["status"].(float64))
					headers := cacheData["headers"].(map[string]any)
					body := []byte(cacheData["body"].(string))

					for k, v := range headers {
						c.Response().Header().Set(k, v.(string))
					}

					c.Response().WriteHeader(status)
					c.Response().Write(body)
					return nil
				}
			}

			// Capture the response body
			resWriter := &bodyCapturingWriter{
				ResponseWriter: c.Response().Writer,
				Buffer:         new(bytes.Buffer),
			}
			c.Response().Writer = resWriter

			// Process the request
			err = next(c)
			if err != nil {
				return err
			}

			status := c.Response().Status
			if status >= 200 && status < 300 {
				// Store response in cache
				responseData := map[string]any{
					"status":  status,
					"headers": captureHeaders(c.Response().Header()),
					"body":    resWriter.Buffer.String(),
				}
				cache.Set(c.Request().Context(), key, responseData, ttl)
			}

			return nil
		}
	}
}

type bodyCapturingWriter struct {
	http.ResponseWriter
	Buffer *bytes.Buffer
}

func (w *bodyCapturingWriter) Write(b []byte) (int, error) {
	w.Buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

func captureHeaders(headers http.Header) map[string]any {
	result := make(map[string]any)
	for k, v := range headers {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
