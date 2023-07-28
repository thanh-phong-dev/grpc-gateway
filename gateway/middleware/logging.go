package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			requestBody = readAndRestoreRequestBody(c.Request)
			blw         = &bodyLogWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}
		)

		c.Writer = blw
		c.Next()
		requestid, _ := c.Request.Context().Value("requestID").(string)
		// Log the request and response details
		fmt.Printf("[%s][%s][%s] Request: %s | Response: %s\n",
			c.Request.Method,
			c.Request.URL.Path,
			requestid,
			requestBody,
			blw.body.String(),
		)
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID := uuid.New().String()
		ctx := context.WithValue(c.Request.Context(), "requestID", xRequestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func readAndRestoreRequestBody(r *http.Request) string {
	// Read the request body
	body, _ := io.ReadAll(r.Body)

	// Restore the request body for further use
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}
