package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type AssignmentRequest struct {
	VariantIDs []string `json:"variant_ids" binding:"required"`
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
	ID          int32    `json:"id"`
	ExternalID  string   `json:"external_id"`
	Url         string   `json:"url"`
	MimeType    string   `json:"mime_type,omitempty"`
	FileSize    int64    `json:"file_size,omitzero"`
	Assignments []string `json:"assignments,omitempty"`
}

// ------------------------------------------ Handlers ------------------------------------------

// @Summary Upload product images
// @Schemes http
// @Description upload product images by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param files formData file true "Image file"
// @Produce json
// @Success 200 {object} ApiResponse[[]ImageResponse]
// @Failure 404 {object} ApiResponse[[]ImageResponse]
// @Failure 500 {object} ApiResponse[[]ImageResponse]
// @Router /images/product/{product_id} [post]
func (sv *Server) uploadProductImages(c *gin.Context) {
	var param EntityIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", err))
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("missing files in request")))
		return
	}

	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("maximum 5 files allowed")))
		return
	}

	assignmentsReq := c.PostFormArray("assignments[]")
	if len(assignmentsReq) == 0 {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("missing assignments in request")))
		return
	}

	roles := c.PostFormArray("roles")
	if len(roles) == 0 {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("missing roles in request")))
		return
	}
	if len(roles) != len(assignmentsReq) {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("roles and assignments must be the same length")))
		return
	}
	// Check if the roles are valid

	// Parse the assignments JSON
	// log.Debug().Msgf("Assignments: %v", assignmentsReq)
	var assignmentsList [][]string
	for _, assignment := range assignmentsReq {
		var assignments []string
		if err := json.Unmarshal([]byte(assignment), &assignments); err != nil {
			c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InvalidBodyCode, "", errors.New("invalid assignments format")))
			return
		}
		assignmentsList = append(assignmentsList, assignments)
	}

	existingProduct, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ID: uuid.MustParse(param.EntityID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[[]ImageResponse](NotFoundCode, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "", err))
		return
	}

	productImages, err := sv.repo.GetImagesByEntityID(c, existingProduct.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "", err))
		return
	}

	if len(productImages)+len(files) > 10 {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "", errors.New("maximum 10 files allowed please remove old files")))
		return
	}

	images := make([]ImageResponse, len(files))

	for i, file := range files {
		id, url, err := sv.uploadService.UploadFile(c, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[[]ImageResponse](UploadFileCode, "upload to server failed", err))
			return
		}

		img, err := sv.repo.CreateImage(c, repository.CreateImageParams{
			ExternalID: id,
			Url:        url,
			MimeType:   utils.GetPgTypeText(file.Header.Get("Content-Type")),
			FileSize:   utils.GetPgTypeInt8(file.Size),
		})

		createImageAssignmentReq := make([]repository.CreateBulkImageAssignmentsParams, 0)
		createImageAssignmentReq = append(createImageAssignmentReq, repository.CreateBulkImageAssignmentsParams{
			ImageID:      img.ID,
			EntityID:     existingProduct.ID,
			EntityType:   repository.ProductEntityType,
			DisplayOrder: int16(i) + 1,
			Role:         strings.ToLower(roles[i]),
		})

		for _, assignment := range assignmentsList[i] {
			if assignment == "" {
				continue
			}

			variantID, err := uuid.Parse(assignment)
			if err != nil {
				c.JSON(http.StatusBadRequest, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "", errors.New("invalid variant id")))

				return
			}

			createImageAssignmentReq = append(createImageAssignmentReq, repository.CreateBulkImageAssignmentsParams{
				ImageID:      img.ID,
				EntityID:     variantID,
				EntityType:   repository.VariantEntityType,
				DisplayOrder: int16(i) + 1,
				Role:         repository.GalleryRole,
			})
		}
		_, err = sv.repo.CreateBulkImageAssignments(c, createImageAssignmentReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "save image info to db failed", err))
			return
		}

		images[i] = ImageResponse{
			ID:         img.ID,
			ExternalID: img.ExternalID,
			Url:        img.Url,
			MimeType:   img.MimeType.String,
			FileSize:   img.FileSize.Int64,
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
// @Success 200 {object} ApiResponse[[]ImageResponse]
// @Failure 404 {object} ApiResponse[ImageResponse]
// @Failure 500 {object} ApiResponse[ImageResponse]
// @Router /images/product/{product_id} [get]
func (sv *Server) getImages(c *gin.Context) {
	var param EntityIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ImageResponse](InvalidBodyCode, "", err))
		return
	}

	images, err := sv.repo.GetImagesByEntityID(c, uuid.MustParse(param.EntityID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ImageResponse](InternalServerErrorCode, fmt.Sprintf("server error"), err))
		return
	}

	resp := make([]ImageResponse, len(images))
	for i, image := range images {
		resp[i] = ImageResponse{
			ID:          image.ID,
			ExternalID:  image.ExternalID,
			Url:         image.Url,
			MimeType:    image.MimeType.String,
			FileSize:    image.FileSize.Int64,
			Assignments: nil,
		}

	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /images/product/{product_id} [delete]
func (sv *Server) removeImage(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse[bool](UnauthorizedCode, "", errors.New("missing user payload in context")))
		return
	}

	var params RemoveImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	image, err := sv.repo.GetImageFromID(c, repository.GetImageFromIDParams{
		ID:         params.ImageID,
		EntityType: repository.ProductEntityType,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", errors.New("image not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	if image.EntityType != repository.ProductEntityType {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", errors.New("image not found")))
		return
	}

	msg, err := sv.removeImageUtil(c, image.ExternalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, msg, err))
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
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /images/{publicID} [delete]
func (sv *Server) removeImageByPublicID(c *gin.Context) {
	_, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](UnauthorizedCode, "", errors.New("missing user payload in context")))
		return
	}
	var params PublicIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}
	msg, err := sv.removeImageUtil(c, params.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, msg, err))
		return
	}

	image, err := sv.repo.GetImageFromExternalID(c, params.PublicID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteProductImage(c, image.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) removeImageUtil(c *gin.Context, publicID string) (msg string, err error) {
	return sv.uploadService.RemoveFile(c, publicID)
}
