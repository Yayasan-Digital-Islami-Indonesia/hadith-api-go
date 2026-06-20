package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// Search godoc
// @Summary      Search hadiths
// @Description  Full-text search across hadith texts (Arabic, English, Indonesian)
// @Tags         search
// @Produce      json
// @Param        q     query  string  true  "Search query"
// @Param        page  query  int     false "Page number"   default(1)
// @Param        limit query  int     false "Items per page (max 100)" default(20)
// @Success      200 {object} models.SearchListResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "missing search query"})
		return
	}

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("limit"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	results, total, err := h.searchService.Search(query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	data := make([]models.SearchResult, len(results))
	for i, r := range results {
		data[i] = models.SearchResult{Hadith: r.Hadith, Text: r.Text, Lang: r.Lang}
	}

	c.JSON(http.StatusOK, models.SearchListResponse{
		Data: data,
		Pagination: models.Pagination{
			Page:  page,
			Limit: pageSize,
			Total: int(total),
		},
	})
}