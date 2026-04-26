// Package response provides shared HTTP response helpers for all handlers.
// Never call c.JSON directly in handlers — always use these helpers.
//
// Response formats:
//
//	Success:    { "data": any, "timestamp": string }
//	Paginated:  { "data": any, "total": int, "page": int, "limit": int, "timestamp": string }
//	Error:      { "error": string, "code": string, "timestamp": string }
//
// Usage:
//
//	response.Success(c, data)
//	response.SuccessPaginated(c, data, total, p)
//	response.NotFound(c, "hadith not found")
//	response.BadRequest(c, "invalid parameter")
//	response.InternalError(c)
package response

import (
	"net/http"
	"time"

	"hadith-api-go/pkg/pagination"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func SuccessPaginated(c *gin.Context, data any, total int, p pagination.Params) {
	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"total":     total,
		"page":      p.Page,
		"limit":     p.Limit,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":     message,
		"code":      "not found",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":     message,
		"code":      "bad request",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":      "internal server error",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
