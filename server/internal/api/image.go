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
type PublicIDParam struct {
	PublicID string `uri:"public_id" binding:"required"`
}
type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"id" binding:"omitempty,uuid"`
	ImageID   int32   `uri:"image_id" binding:"required"`
}

type EntityIDParam struct {
	EntityID string `uri:"entity_id" binding:"required,uuid"`
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

// @Summary Upload images
// @Schemes http
// @Description upload images
// @Tags images
// @Accept json
// @Param files formData file true "Image file"
// @Produce json
// @Success 200 {object} GenericResponse[[]repository.Images]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images [post]
func (sv *Server) uploadImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	variantIDs := c.PostFormArray("variant_ids")
	log.Debug().Msgf("variant_ids: %v", variantIDs)
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("missing files in request")))
		return
	}
	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("maximum 5 files allowed")))
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
			DisplayOrder: int16(i),
		})
	}
	c.JSON(http.StatusCreated, createSuccessResponse(c, images, "", nil, nil))
}

// @Summary Upload product images
// @Schemes http
// @Description upload product images by ID
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
	var param EntityIDParam
	if err := c.ShouldBindUri(&param); err != nil {
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

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ID: uuid.MustParse(param.EntityID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	productImages, err := sv.repo.GetImagesByEntityID(c, existingProduct.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	if len(productImages)+len(files) > 3 {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("maximum 3 files allowed please remove old files")))
		return
	}

	images := make([]ImageResponse, len(files))

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
			DisplayOrder: int16(i) + 1,
			Role:         repository.ProductRole,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "save image info to db failed", err))
			return
		}
		images[i] = ImageResponse{
			ID:           img.ID,
			ExternalID:   img.ExternalID,
			Url:          img.Url,
			MimeType:     img.MimeType.String,
			FileSize:     img.FileSize.Int64,
			EntityID:     assignment.EntityID.String(),
			EntityType:   assignment.EntityType,
			DisplayOrder: assignment.DisplayOrder,
			Role:         assignment.Role,
		}
	}
	c.JSON(http.StatusCreated, createSuccessResponse(c, images, "", nil, nil))

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
func (sv *Server) getImages(c *gin.Context) {
	var param EntityIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	images, err := sv.repo.GetImagesByEntityID(c, uuid.MustParse(param.EntityID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, fmt.Sprintf("server error"), err))
		return
	}
	if len(images) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, fmt.Sprintf("not images found"), errors.New("product not found")))
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
func (sv *Server) removeImage(c *gin.Context) {
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

	image, err := sv.repo.GetImageFromID(c, repository.GetImageFromIDParams{
		ID:         params.ImageID,
		EntityType: repository.ProductImageType,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", errors.New("image not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	if image.EntityType != repository.ProductImageType {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("image not found")))
		return
	}

	msg, err := sv.removeImageUtil(c, image.ExternalID)
	if err != nil {
		createErrorResponse(http.StatusBadRequest, msg, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Remove a product by external ID
// @Schemes http
// @Description remove a product by external ID
// @Tags images
// @Accept json
// @Param publicID path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/{publicID} [delete]
func (sv *Server) removeImageByPublicID(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", errors.New("missing user payload in context")))
		return
	}
	var params PublicIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	msg, err := sv.removeImageUtil(c, params.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, msg, err))
		return
	}

	image, err := sv.repo.GetImageFromExternalID(c, params.PublicID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	err = sv.repo.DeleteProductImage(c, image.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) removeImageUtil(c *gin.Context, publicID string) (msg string, err error) {
	return sv.uploadService.RemoveFile(c, publicID)
}
