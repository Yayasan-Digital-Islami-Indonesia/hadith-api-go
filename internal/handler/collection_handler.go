package handler

import (
	"errors"

	"hadith-api-go/internal/domain"
	"hadith-api-go/internal/domain/collection"
	"hadith-api-go/pkg/pagination"
	"hadith-api-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type CollectionHandler struct {
	service collection.CollectionService
}

func NewCollectionHandler(service collection.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

// List godoc
// @Summary      List all hadith collections
// @Description  Returns a paginated list of all available hadith collections
// @Tags         collections
// @Produce      json
// @Param        page   query  int  false  "Page number (default: 1)"
// @Param        limit  query  int  false  "Items per page (default: 20, max: 100)"
// @Success      200  {object}  map[string]any
// @Router       /v1/collections [get]
func (h *CollectionHandler) List(c *gin.Context) {
	p := pagination.Parse(c.Query("page"), c.Query("limit"))

	collections, total, err := h.service.List(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.SuccessPaginated(c, collections, total, p)
}

// Detail godoc
// @Summary      Get a hadith collection
// @Description  Returns details of a specific hadith collection by name
// @Tags         collections
// @Produce      json
// @Param        name  path  string  true  "Collection name (e.g., bukhari)"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /v1/collections/{name} [get]
func (h *CollectionHandler) Detail(c *gin.Context) {
	name := c.Param("name")

	col, err := h.service.FindByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "collection not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, col)
}
