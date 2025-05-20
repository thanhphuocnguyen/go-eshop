package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// @Summary Upload product images
// @Schemes http
// @Description upload product images by ID
// @Tags images
// @Accept json
// @Param product_id path int true "Product ID"
// @Param files formData file true "Image file"
// @Produce json
// @Success 200 {object} ApiResponse[[]ImageResponse]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [post]
func (sv *Server) uploadProductImagesHandler(c *gin.Context) {
	var param UriIDParam
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
		ID: uuid.MustParse(param.ID),
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

	images := make([]ImageResponse, 0)
	imagesChan := make(chan ImageResponse, len(files))
	var wg sync.WaitGroup
	for i, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()
			id, url, err := sv.uploadService.UploadFile(c, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[[]ImageResponse](UploadFileCode, "upload to server failed", err))
				return
			}

			mimeType := mime.TypeByExtension(file.Filename)

			img, err := sv.repo.InsertImage(c, repository.InsertImageParams{
				ExternalID: id,
				Url:        url,
				MimeType:   &mimeType,
				FileSize:   &file.Size,
			})

			createImageAssignmentReq := make([]repository.InsertBulkImageAssignmentsParams, 0)
			createImageAssignmentReq = append(createImageAssignmentReq, repository.InsertBulkImageAssignmentsParams{
				ImageID:      img.ID,
				EntityID:     existingProduct.ID,
				EntityType:   string(repository.EntityTypeProduct),
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

				createImageAssignmentReq = append(createImageAssignmentReq, repository.InsertBulkImageAssignmentsParams{
					ImageID:      img.ID,
					EntityID:     variantID,
					EntityType:   string(repository.EntityTypeProductVariant),
					DisplayOrder: int16(i) + 1,
					Role:         repository.GalleryRole,
				})
			}
			_, err = sv.repo.InsertBulkImageAssignments(c, createImageAssignmentReq)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[[]ImageResponse](InternalServerErrorCode, "save image info to db failed", err))
				return
			}

			imagesChan <- ImageResponse{
				ID:         img.ID.String(),
				ExternalID: img.ExternalID,
				Url:        img.Url,
				MimeType:   *img.MimeType,
				FileSize:   *img.FileSize,
			}
		}()
	}
	wg.Wait()
	close(imagesChan)
	for img := range imagesChan {
		images = append(images, img)
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
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [get]
func (sv *Server) getProductImagesHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ImageResponse](InvalidBodyCode, "", err))
		return
	}

	images, err := sv.repo.GetImagesByEntityID(c, uuid.MustParse(param.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ImageResponse](InternalServerErrorCode, fmt.Sprintf("server error"), err))
		return
	}

	resp := make([]ImageResponse, len(images))
	for i, image := range images {
		resp[i] = ImageResponse{
			ID:          image.ID.String(),
			ExternalID:  image.ExternalID,
			Url:         image.Url,
			MimeType:    *image.MimeType,
			FileSize:    *image.FileSize,
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
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/product/{product_id} [delete]
func (sv *Server) removeImageHandler(c *gin.Context) {
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
		ID:         uuid.MustParse(params.ImageID),
		EntityType: string(repository.EntityTypeProduct),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", errors.New("image not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	if image.EntityType != string(repository.EntityTypeProductVariant) {
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
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /images/{publicID} [delete]
func (sv *Server) removeImageByPublicIDHandler(c *gin.Context) {
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
