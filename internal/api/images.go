package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// @Summary Get list of product image by ID
// @Schemes http
// @Description get list of product image by ID
// @Tags images
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[[]ImageResponse]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/product/{productId} [get]
func (sv *Server) GetProductImagesHandler(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, []models.ImageModel{}, nil, nil))
}

// @Summary Remove a product by external ID
// @Schemes http
// @Description remove a product by external ID
// @Tags images
// @Accept json
// @Param publicID path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/{publicID} [delete]
func (sv *Server) RemoveImageByPublicIDHandler(c *gin.Context) {
	_, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(UnauthorizedCode, errors.New("missing user payload in context")))
		return
	}
	var params models.PublicIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	_, err := sv.removeImageUtil(c, params.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) removeImageUtil(c *gin.Context, publicID string) (msg string, err error) {
	return sv.uploadService.Remove(c, publicID)
}
