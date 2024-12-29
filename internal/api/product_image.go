package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

// ------------------------------------------ Request and Response ------------------------------------------
type getProductImageParams struct {
	ID      int64 `uri:"id" binding:"required"`
	ImageID int32 `uri:"image_id" binding:"required"`
}

// UploadProductImage godoc
// @Summary Upload a product image by ID
// @Schemes http
// @Description upload a product image by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Param file formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]sqlc.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/upload-image [post]
func (sv *Server) uploadProductImage(c *gin.Context) {
	var param getProductParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	productDetailRows, err := sv.postgres.GetProductDetail(c, sqlc.GetProductDetailParams{
		ID: param.ID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	if len(productDetailRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
		return
	}

	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("missing file in request")))
		return
	}
	// file name is public id
	fileName := util.GetImageName(file.Filename)
	url, err := sv.uploadService.UploadFile(c, file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	img, err := sv.postgres.CreateImage(c, sqlc.CreateImageParams{
		ProductID: pgtype.Int8{
			Int64: productDetailRows[0].Product.ID,
			Valid: true,
		},
		ImageUrl: url,
		ExternalID: pgtype.Text{
			String: fileName,
			Valid:  true,
		},
	})

	if !productDetailRows[0].ImageUrl.Valid {
		err := sv.postgres.SetPrimaryImageTx(c, postgres.SetPrimaryImageTxParams{
			NewPrimaryID: img.ID,
			ProductID:    productDetailRows[0].Product.ID,
		})
		if err == nil {
			img.IsPrimary = pgtype.Bool{
				Bool:  true,
				Valid: true,
			}
		}
	}

	if err != nil {
		sv.uploadService.RemoveFile(c, fileName)
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[sqlc.Image]{&img, nil, nil})
}

// GetProductImages godoc
// @Summary Get list of product images by ID
// @Schemes http
// @Description get list of product images by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[[]sqlc.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/images [get]
func (sv *Server) getProductImages(c *gin.Context) {
	var param getProductParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	images, err := sv.postgres.GetImagesByProductID(c, pgtype.Int8{
		Int64: param.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[[]sqlc.Image]{&images, nil, nil})
}

// SetImagesPrimary godoc
// @Summary Set a product image as primary by ID
// @Schemes http
// @Description set a product image as primary by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/images/{image_id}/primary [put]
func (sv *Server) setImagesPrimary(c *gin.Context) {

	var params getProductImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: params.ID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
	}

	img, err := sv.postgres.GetImageByID(c, params.ImageID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if img.ProductID.Int64 != product.ID {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("image not found")))
		return
	}

	err = sv.postgres.SetPrimaryImageTx(c, postgres.SetPrimaryImageTxParams{
		NewPrimaryID: img.ID,
		ProductID:    product.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	success := true
	message := "Set primary image successfully"
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// RemoveProductImage godoc
// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/remove-image [delete]
func (sv *Server) removeProductImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(errors.New("missing user payload in context")))
		return
	}

	var param getProductParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var params getProductImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: param.ID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	img, err := sv.postgres.GetImageByID(c, params.ImageID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if img.ProductID.Int64 != product.ID {
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
	err = sv.postgres.DeleteImage(c, img.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &result, nil})
}
