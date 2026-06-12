package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *ChapterHandler) GetChapterHadiths(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("chapter_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chapter id"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": hadiths,
		"pagination": gin.H{
			"page":  page,
			"limit": pageSize,
			"total": total,
		},
	})
}