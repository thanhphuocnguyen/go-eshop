package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ------------------------------------------ Request and Response ------------------------------------------
type getProductImageParams struct {
	ID      int64 `uri:"id" binding:"required"`
	ImageID int32 `uri:"image_id" binding:"required"`
}
type productImageParam struct {
	ID int64 `uri:"product_id" binding:"required"`
}
type productImageModel struct {
	ID        int32  `json:"id"`
	ProductID int32  `json:"product_id"`
	ImageUrl  string `json:"image_url"`
}

// ------------------------------------------ Handlers ------------------------------------------

// @Summary Upload a product image by ID
// @Schemes http
// @Description upload a product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param file formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product [post]
func (sv *Server) uploadProductImage(c *gin.Context) {
	var param productImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: param.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("missing file in request")))
		return
	}
	// file name is public id
	fileName := GetImageName(file.Filename)
	url, err := sv.uploadService.UploadFile(c, file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	img, err := sv.repo.CreateImage(c, repository.CreateImageParams{
		ProductID:  util.GetPgTypeInt8(existingProduct.ProductID),
		ImageUrl:   url,
		ExternalID: util.GetPgTypeText(fileName),
	})

	if err != nil {
		sv.uploadService.RemoveFile(c, fileName)
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Image]{&img, nil, nil})
}

// @Summary Get list of product image by ID
// @Schemes http
// @Description get list of product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [get]
func (sv *Server) getProductImage(c *gin.Context) {
	var param productImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	images, err := sv.repo.GetImagesByProductID(c, pgtype.Int8{
		Int64: param.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[[]repository.Image]{&images, nil, nil})
}

// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id}/remove [delete]
func (sv *Server) removeProductImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(errors.New("missing user payload in context")))
		return
	}

	var param productImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var params getProductImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: param.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	img, err := sv.repo.GetImageByID(c, params.ImageID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if img.ProductID.Int64 != existingProduct.ProductID {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("image not found")))
		return
	}

	var result string
	if img.ExternalID.Valid {
		result, err = sv.uploadService.RemoveFile(c, img.ExternalID.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}
	err = sv.repo.DeleteImage(c, img.ImageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[bool]{nil, &result, nil})
}
