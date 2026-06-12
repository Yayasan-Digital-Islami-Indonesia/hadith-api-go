package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/services"
)

type HadithHandler struct {
	hadithService *services.HadithService
	bookService   *services.BookService
}

func NewHadithHandler(hadithService *services.HadithService, bookService *services.BookService) *HadithHandler {
	return &HadithHandler{
		hadithService: hadithService,
		bookService:   bookService,
	}
}

func (h *HadithHandler) GetHadith(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hadith id"})
		return
	}

	hadith, err := h.hadithService.GetHadith(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

func (h *HadithHandler) GetHadithByNumber(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hadith number"})
		return
	}

	hadith, err := h.hadithService.GetHadithByBookAndNumber(uint(bookID), number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

func (h *HadithHandler) GetRandomHadith(c *gin.Context) {
	hadith, err := h.hadithService.GetRandomHadith()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hadith)
}