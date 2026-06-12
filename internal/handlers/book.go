package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	identifier := c.Param("id")
	book, err := h.bookService.GetBookOrBySlug(identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetChapters(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	chapters, err := h.chapterService.GetChaptersByBook(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": chapters})
}