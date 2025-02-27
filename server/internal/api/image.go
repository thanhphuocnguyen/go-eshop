package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------------------ Request and Response ------------------------------------------
type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"id" binding:"omitempty,uuid"`
	ImageID   int32   `uri:"image_id" binding:"required"`
}
type ProductImageParam struct {
	ProductID string `uri:"product_id" binding:"required,uuid"`
}
type VariantImageParam struct {
	ProductImageParam
	VariantID string `uri:"variant_id" binding:"required,uuid"`
}
type ImageModel struct {
	ID        int32   `json:"id"`
	ProductID *string `json:"product_id"`
	VariantID *string `json:"variant_id"`
	ImageUrl  string  `json:"image_url"`
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
	var param ProductImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: uuid.MustParse(param.ProductID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingImg, err := sv.repo.GetImageByProductID(c, pgtype.UUID{
		Bytes: existingProduct.ProductID,
	})
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	sv.uploadService.RemoveFile(c, existingImg.ExternalID.String)

	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("missing file in request")))
		return
	}
	// file name is public id
	publicID, url, err := sv.uploadService.UploadFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	img, err := sv.repo.CreateImage(c, repository.CreateImageParams{
		ProductID:  utils.GetPgTypeUUID(existingProduct.ProductID),
		ImageUrl:   url,
		ExternalID: utils.GetPgTypeText(publicID),
	})

	if err != nil {
		sv.uploadService.RemoveFile(c, publicID)
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
// @Success 200 {object} GenericResponse[repository.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [get]
func (sv *Server) getProductImage(c *gin.Context) {
	var param ProductImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	images, err := sv.repo.GetImageByProductID(c, utils.GetPgTypeUUID(uuid.MustParse(param.ProductID)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[repository.Image]{&images, nil, nil})
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

	var params RemoveImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: uuid.MustParse(params.ID),
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

	if img.ProductID.Bytes != existingProduct.ProductID {
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

// @Summary Upload a variant image by ID
// @Schemes http
// @Description upload a variant image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param variant_id path int true "Variant ID"
// @Param file formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Image]
func (sv *Server) uploadVariantImage(c *gin.Context) {
	var param VariantImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingVariant, err := sv.repo.GetVariantByID(c, repository.GetVariantByIDParams{
		ProductID: uuid.MustParse(param.ProductID),
		VariantID: uuid.MustParse(param.VariantID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("variant not found")))
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
	publicID, url, err := sv.uploadService.UploadFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	img, err := sv.repo.CreateImage(c, repository.CreateImageParams{
		ProductID:  utils.GetPgTypeUUID(existingVariant.ProductID),
		VariantID:  utils.GetPgTypeUUID(existingVariant.VariantID),
		ImageUrl:   url,
		ExternalID: utils.GetPgTypeText(publicID),
	})

	if err != nil {
		sv.uploadService.RemoveFile(c, publicID)
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Image]{&img, nil, nil})
}

// @Summary Get list of variant image by ID
// @Schemes http
// @Description get list of variant image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param variant_id path int true "Variant ID"
// @Produce json
// @Success 200 {object} GenericResponse[repository.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id}/variant/{variant_id} [get]
func (sv *Server) getVariantImage(c *gin.Context) {
	var param VariantImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	images, err := sv.repo.GetImageByVariantID(c, utils.GetPgTypeUUID(uuid.MustParse(param.VariantID)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[repository.Image]{&images, nil, nil})
}

// @Summary Remove a variant image by ID
// @Schemes http
// @Description remove a variant image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param variant_id path int true "Variant ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id}/variant/{variant_id}/remove [delete]
func (sv *Server) removeVariantImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(errors.New("missing user payload in context")))
		return
	}

	var params RemoveImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	if params.VariantID == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("missing variant_id in request")))
		return
	}

	existingProduct, err := sv.repo.GetVariantByID(c, repository.GetVariantByIDParams{
		ProductID: uuid.MustParse(params.ID),
		VariantID: uuid.MustParse(*params.VariantID),
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

	if uuid.UUID(img.VariantID.Bytes).String() != existingProduct.VariantID.String() {
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
