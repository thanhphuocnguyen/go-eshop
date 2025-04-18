package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

type ProductVariantImageParam struct {
	VariantID string `uri:"variant_id" binding:"required,uuid"`
}

type ProductVariantImageModel struct {
	ID        string    `json:"id"`
	VariantID uuid.UUID `json:"variant_id,omitempty"`
	ImageUrl  string    `json:"image_url"`
	ImageID   string    `json:"image_id"`
}

type ImageResponse struct {
	ID           int32  `json:"id"`
	ExternalID   string `json:"external_id"`
	Url          string `json:"url"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitzero"`
	EntityID     string `json:"entity_id,omitempty"`
	EntityType   string `json:"entity_type,omitempty"`
	DisplayOrder int16  `json:"display_order"`
	Role         string `json:"role,omitempty"`
}

// ------------------------------------------ Handlers ------------------------------------------

// @Summary Upload a product images
// @Schemes http
// @Description upload a product images by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param files formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Images]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [post]
func (sv *Server) uploadProductImages(c *gin.Context) {
	var param ProductImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ID: uuid.MustParse(param.ProductID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("missing files in request")))
		return
	}
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("maximum 5 files allowed")))
		return
	}

	productImages, err := sv.repo.GetProductImagesProductID(c, uuid.MustParse(param.ProductID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	if len(productImages)+len(files) > 3 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("maximum 3 files allowed please remove old files")))
		return
	}

	images := make([]ImageResponse, 0)
	for i, file := range files {
		id, url, err := sv.uploadService.UploadFile(c, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "upload to server failed", err))
			return
		}
		img, err := sv.repo.CreateImage(c, repository.CreateImageParams{
			ExternalID: id,
			Url:        url,
			MimeType:   utils.GetPgTypeText(file.Header.Get("Content-Type")),
			FileSize:   utils.GetPgTypeInt8(file.Size),
		})
		assignment, err := sv.repo.CreateImageAssignment(c, repository.CreateImageAssignmentParams{
			ImageID:      img.ID,
			EntityID:     existingProduct.ID,
			EntityType:   repository.ProductImageType,
			DisplayOrder: int16(i),
			Role:         repository.ProductRole,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "save image info to db failed", err))
			return
		}
		images = append(images, ImageResponse{
			ID:           img.ID,
			ExternalID:   img.ExternalID,
			Url:          img.Url,
			MimeType:     img.MimeType.String,
			FileSize:     img.FileSize.Int64,
			EntityID:     assignment.EntityID.String(),
			EntityType:   assignment.EntityType,
			DisplayOrder: assignment.DisplayOrder,
			Role:         assignment.Role,
		})
	}
	c.JSON(http.StatusCreated, createSuccessResponse(c, images, "", nil, nil))

}

// @Summary Upload a product variant image
// @Schemes http
// @Description upload a product variant image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param file formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Image]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{variant_id} [post]
func (sv *Server) uploadProductVariantImage(c *gin.Context) {
	var param ProductVariantImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	existingVariant, err := sv.repo.GetProductVariantByID(c, uuid.MustParse(param.VariantID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to bind file: %v", err), err))
		return
	}

	if file == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, fmt.Sprintf("missing file in request"), errors.New("missing file in request")))
		return
	}
	// file name is public id
	publicID, url, err := sv.uploadService.UploadFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	img, err := sv.repo.UpdateProductVariant(c, repository.UpdateProductVariantParams{
		ID:       existingVariant.ID,
		ImageUrl: utils.GetPgTypeText(url),
		ImageID:  utils.GetPgTypeText(publicID),
	})

	if err != nil {
		sv.uploadService.RemoveFile(c, publicID)
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, fmt.Sprintf("failed to remove file: %v", publicID), err))
		return
	}

	if existingVariant.ImageID.Valid {
		errRemoveFile := sv.removeProductVariantImageUtil(c, existingVariant.ImageID.String)
		if errRemoveFile != nil {
			log.Error().Err(err).Str("external_id", existingVariant.ImageID.String).Msg("remove product variant image error")
			// TODO: push the error to the queue
		}
	}

	resp := ProductVariantImageModel{
		ID:        img.ImageID.String,
		VariantID: existingVariant.ID,
		ImageUrl:  img.ImageUrl.String,
		ImageID:   img.ImageID.String,
	}
	c.JSON(http.StatusCreated, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Get list of product image by ID
// @Schemes http
// @Description get list of product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [get]
func (sv *Server) getProductImages(c *gin.Context) {
	var param ProductImageParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	images, err := sv.repo.GetProductImagesProductID(c, uuid.MustParse(param.ProductID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	if len(images) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", errors.New("product not found")))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, images, "", nil, nil))
}

// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [delete]
func (sv *Server) removeProductImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", errors.New("missing user payload in context")))
		return
	}

	var params RemoveImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	err := sv.removeImageUtil(c, params.ImageID)
	if err != nil {
		createErrorResponse(http.StatusBadRequest, "", err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) removeProductVariantImageUtil(c *gin.Context, imgID string) (err error) {
	_, err = sv.uploadService.RemoveFile(c, imgID)
	return
}

func (sv *Server) removeImageUtil(c *gin.Context, imageID int32) (err error) {
	img, err := sv.repo.GetImageFromID(c, imageID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	msg, err := sv.uploadService.RemoveFile(c, img.ExternalID)
	err = sv.repo.DeleteProductImage(c, img.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, msg, err))
		return
	}
	return
}
