package handler

import (
	"errors"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/hadith"
	"hadith-api-go/pkg/pagination"
	"hadith-api-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type HadithHandler struct {
	service hadith.HadithService
}

func NewHadithHandler(service hadith.HadithService) *HadithHandler {
	return &HadithHandler{service: service}
}

// ListByBook godoc
// @Summary      List hadiths in a book
// @Description  Returns a paginated list of hadiths within a specific book of a collection
// @Tags         hadiths
// @Produce      json
// @Param        name         path   string  true   "Collection name (e.g., bukhari)"
// @Param        book_number  path   string  true   "Book number"
// @Param        page         query  int     false  "Page number (default: 1)"
// @Param        limit        query  int     false  "Items per page (default: 20, max: 100)"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/collections/{name}/books/{book_number}/hadiths [get]
func (h *HadithHandler) ListByBook(c *gin.Context) {
	collectionName := c.Param("name")
	bookNumber := c.Param("bookNumber")
	p := pagination.Parse(c.Query("page"), c.Query("limit"))

	hadiths, total, err := h.service.ListByBook(c.Request.Context(), collectionName, bookNumber, p.Limit, p.Offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "book not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.SuccessPaginated(c, hadiths, total, p)
}

// FindByCollection godoc
// @Summary      Get a specific hadith
// @Description  Returns a specific hadith by its number within a collection
// @Tags         hadiths
// @Produce      json
// @Param        name          path  string  true  "Collection name (e.g., bukhari)"
// @Param        hadith_number path  string  true  "Hadith number"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/collections/{name}/hadiths/{hadith_number} [get]
func (h *HadithHandler) FindByCollection(c *gin.Context) {
	collectionName := c.Param("name")
	hadithNumber := c.Param("hadithNumber")

	had, err := h.service.FindByCollection(c.Request.Context(), collectionName, hadithNumber)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "hadith not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, had)
}

// Random godoc
// @Summary      Get a random hadith
// @Description  Returns a random hadith from the database
// @Tags         hadiths
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /v1/hadiths/random [get]
func (h *HadithHandler) Random(c *gin.Context) {
	had, err := h.service.Random(c.Request.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "no hadiths available")
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, had)
}
