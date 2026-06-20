package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/services"
)

type ChapterHandler struct {
	hadithService  *services.HadithService
	chapterService *services.ChapterService
}

func NewChapterHandler(hadithService *services.HadithService, chapterService *services.ChapterService) *ChapterHandler {
	return &ChapterHandler{
		hadithService:  hadithService,
		chapterService: chapterService,
	}
}

// GetChapterHadiths godoc
// @Summary      List hadiths in a chapter
// @Description  Returns paginated hadiths for the given chapter
// @Tags         chapters
// @Produce      json
// @Param        chapter_id  path  int  true  "Chapter ID"
// @Param        page        query int  false "Page number"   default(1)
// @Param        limit       query int  false "Items per page (max 100)" default(20)
// @Success      200 {object} models.HadithListResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /books/{id}/chapters/{chapter_id} [get]
func (h *ChapterHandler) GetChapterHadiths(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid chapter id"})
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

	hadiths, total, err := h.hadithService.GetHadithsByChapter(uint(chapterID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.HadithListResponse{
		Data: hadiths,
		Pagination: models.Pagination{
			Page:  page,
			Limit: pageSize,
			Total: int(total),
		},
	})
}