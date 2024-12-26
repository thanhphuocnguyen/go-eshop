package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

// UploadProductImage godoc
// @Summary Upload a product image by ID
// @Schemes http
// @Description upload a product image by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Param file formData file true "Image file"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/upload-image [post]
func (sv *Server) uploadProductImage(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("missing user payload in context")))
		return
	}
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	var param getProductParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: param.ID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fileName := util.GetImageName(file.Filename, authPayload.UserID, product.ID)
	filePath := imageAssetsDir + fileName
	defer func() {
		os.Remove(filePath)
	}()

	err = c.SaveUploadedFile(file, util.GetImageURL(fileName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url, err := sv.uploadService.UploadFile(c, file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	img, err := sv.postgres.CreateImage(c, sqlc.CreateImageParams{
		ProductID: pgtype.Int8{
			Int64: product.ID,
			Valid: true,
		},
		ImageUrl: url,
		CloudinaryID: pgtype.Text{
			String: fileName,
			Valid:  true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, responseMapper(img, nil, nil))
}

// RemoveProductImage godoc
// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id}/remove-image [delete]
func (sv *Server) removeProductImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("missing user payload in context")))
		return
	}

	var param getProductParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req removeProductImageRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: param.ID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	img, err := sv.postgres.GetImagesByProductID(c, sqlc.GetImagesByProductIDParams{
		ProductID: pgtype.Int8{
			Int64: product.ID,
			Valid: true,
		},
		ImageID: int32(req.ID),
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if img.CloudinaryID.Valid {
		err = sv.uploadService.RemoveFile(c, img.CloudinaryID.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("image not found")))
		return
	}
}
