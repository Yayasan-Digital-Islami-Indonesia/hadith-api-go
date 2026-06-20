package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/models"
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

// GetHadith godoc
// @Summary      Get a hadith by ID
// @Description  Returns a single hadith with its multilingual texts
// @Tags         hadith
// @Produce      json
// @Param        id   path  int  true  "Hadith ID"
// @Success      200 {object} models.Hadith
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /hadith/{id} [get]
func (h *HadithHandler) GetHadith(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid hadith id"})
		return
	}

	hadith, err := h.hadithService.GetHadith(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

// GetHadithByNumber godoc
// @Summary      Get a hadith by book and number
// @Description  Returns a hadith identified by its book ID and per-book hadith number
// @Tags         hadith
// @Produce      json
// @Param        id      path  int  true  "Book ID"
// @Param        number  path  int  true  "Hadith number within the book"
// @Success      200 {object} models.Hadith
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /books/{id}/hadith/{number} [get]
func (h *HadithHandler) GetHadithByNumber(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid book id"})
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid hadith number"})
		return
	}

	hadith, err := h.hadithService.GetHadithByBookAndNumber(uint(bookID), number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if hadith == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "hadith not found"})
		return
	}

	c.JSON(http.StatusOK, hadith)
}

// GetRandomHadith godoc
// @Summary      Get a random hadith
// @Description  Returns one randomly selected hadith from the collection
// @Tags         hadith
// @Produce      json
// @Success      200 {object} models.Hadith
// @Failure      500 {object} models.ErrorResponse
// @Router       /random [get]
func (h *HadithHandler) GetRandomHadith(c *gin.Context) {
	hadith, err := h.hadithService.GetRandomHadith()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, hadith)
}