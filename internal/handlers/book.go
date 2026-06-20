package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydgi/hadith-api-go/internal/models"
	"github.com/ydgi/hadith-api-go/internal/services"
)

type BookHandler struct {
	bookService    *services.BookService
	chapterService *services.ChapterService
}

func NewBookHandler(bookService *services.BookService, chapterService *services.ChapterService) *BookHandler {
	return &BookHandler{
		bookService:    bookService,
		chapterService: chapterService,
	}
}

// GetBooks godoc
// @Summary      List all hadith books
// @Description  Returns the 6 canonical hadith books (Kutub al-Sittah)
// @Tags         books
// @Produce      json
// @Success      200 {object} models.BookListResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BookListResponse{Data: books})
}

// GetBook godoc
// @Summary      Get a hadith book
// @Description  Returns a single book by numeric ID or slug
// @Tags         books
// @Produce      json
// @Param        id   path  string  true  "Book ID or slug (e.g. 1 or bukhari)"
// @Success      200 {object} models.Book
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	identifier := c.Param("id")
	book, err := h.bookService.GetBookOrBySlug(identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// GetChapters godoc
// @Summary      List chapters of a book
// @Description  Returns all chapters for the given book
// @Tags         books
// @Produce      json
// @Param        id   path  int  true  "Book ID"
// @Success      200 {object} models.ChapterListResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /books/{id}/chapters [get]
func (h *BookHandler) GetChapters(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid book id"})
		return
	}

	chapters, err := h.chapterService.GetChaptersByBook(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ChapterListResponse{Data: chapters})
}