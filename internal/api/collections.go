package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// --- Public API ---

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Shop-Collection-by-slug
// @Accept json
// @Tags Collections
// @Produce json
// @Param slug path string true "Collection slug"
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /collections/{slug} [get]
func (sv *Server) GetCollectionBySlugHandler(c *gin.Context) {
	var param models.URISlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionBySlug(c, param.Slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode,
				fmt.Errorf("category with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	rows, err := sv.repo.GetProductList(c, repository.GetProductListParams{
		CollectionIds: []uuid.UUID{collection.ID},
		Limit:         query.PageSize,
		Offset:        (query.PageSize) * int64(query.Page-1),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}
	collectionResp := dto.CategoryDetail{
		ID:          collection.ID.String(),
		Name:        collection.Name,
		Description: collection.Description,
		Slug:        collection.Slug,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.String(),
		Products:    make([]dto.ProductSummary, len(rows)),
	}
	for i, row := range rows {
		collectionResp.Products[i] = dto.MapToShopProductResponse(row)
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, collectionResp, nil, nil))
}

// Setup collection-related routes
func (sv *Server) addCollectionRoutes(rg *gin.RouterGroup) {
	collections := rg.Group("collections")
	{
		collections.GET("", sv.AdminGetCollectionsHandler)
		collections.GET(":slug", sv.GetCollectionBySlugHandler)
	}
}
