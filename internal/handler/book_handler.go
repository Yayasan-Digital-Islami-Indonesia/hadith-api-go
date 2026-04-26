package handler

import (
	"errors"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/book"
	"hadith-api-go/pkg/pagination"
	"hadith-api-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service book.BookService
}

func NewBookHandler(service book.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// List godoc
// @Summary      List books in a collection
// @Description  Returns a paginated list of books within a hadith collection
// @Tags         books
// @Produce      json
// @Param        name   path   string  true   "Collection name (e.g., bukhari)"
// @Param        page   query  int     false  "Page number (default: 1)"
// @Param        limit  query  int     false  "Items per page (default: 20, max: 100)"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/collections/{name}/books [get]
func (h *BookHandler) List(c *gin.Context) {
	collectionName := c.Param("name")
	p := pagination.Parse(c.Query("page"), c.Query("limit"))

	books, total, err := h.service.List(c.Request.Context(), collectionName, p.Limit, p.Offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "collection not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.SuccessPaginated(c, books, total, p)
}

// Detail godoc
// @Summary      Get a specific book
// @Description  Returns details of a specific book within a hadith collection
// @Tags         books
// @Produce      json
// @Param        name         path  string  true  "Collection name (e.g., bukhari)"
// @Param        book_number  path  string  true  "Book number"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/collections/{name}/books/{book_number} [get]
func (h *BookHandler) Detail(c *gin.Context) {
	collectionName := c.Param("name")
	bookNumber := c.Param("bookNumber")

	b, err := h.service.FindByNumber(c.Request.Context(), collectionName, bookNumber)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "book not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, b)
}
