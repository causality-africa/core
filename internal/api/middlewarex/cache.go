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

type cachedResponse struct {
	Status  int
	Headers map[string][]string
	Body    []byte
}

func responseCacheKey(c echo.Context) string {
	hasher := sha256.New()
	hasher.Write([]byte(c.Request().URL.String()))
	return fmt.Sprintf("core:responses:%s", hex.EncodeToString(hasher.Sum(nil)))
}

func CacheMiddleware(store *cache.Cache, ttl time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			// Try to get from cache
			key := responseCacheKey(c)
			ctx := c.Request().Context()
			cachedResp, err := cache.Get[cachedResponse](store, ctx, key)
			if err == nil {
				for k, values := range cachedResp.Headers {
					for _, v := range values {
						c.Response().Header().Add(k, v)
					}
				}

				c.Response().WriteHeader(cachedResp.Status)
				c.Response().Write(cachedResp.Body)
				return nil
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
				cachedResp.Status = status
				cachedResp.Headers = captureHeaders(c.Response().Header())
				cachedResp.Body = resWriter.Buffer.Bytes()
				cache.Set(store, ctx, key, cachedResp, ttl)
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

func captureHeaders(headers http.Header) map[string][]string {
	result := make(map[string][]string)
	for k, v := range headers {
		result[k] = v
	}

	return result
}
