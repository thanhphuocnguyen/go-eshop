package api

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// @Summary Upload product images
// @Schemes http
// @Description upload product images by ID
// @Tags images
// @Accept json
// @Param productId path int true "Product ID"
// @Param files formData file true "Image file"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/product/{productId} [post]
func (sv *Server) uploadProductImagesHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("missing files in request")))
		return
	}

	if len(files) > 5 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("maximum 5 files allowed")))
		return
	}

	assignmentsReq := c.PostFormArray("assignments[]")
	if len(assignmentsReq) == 0 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("missing assignments in request")))
		return
	}

	roles := c.PostFormArray("roles")
	if len(roles) == 0 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("missing roles in request")))
		return
	}
	if len(roles) != len(assignmentsReq) {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("roles and assignments must be the same length")))
		return
	}
	// Check if the roles are valid
	images := make([]ImageResponse, 0)
	imagesChan := make(chan ImageResponse, len(files))
	var wg sync.WaitGroup

	wg.Wait()
	close(imagesChan)
	for img := range imagesChan {
		images = append(images, img)
	}

	c.JSON(http.StatusCreated, createDataResp(c, images, nil, nil))
}

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
func (sv *Server) getProductImagesHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, []ImageResponse{}, nil, nil))
}

// @Summary Remove a product image by ID
// @Schemes http
// @Description remove a product image by ID
// @Tags images
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/product/{productId} [delete]
func (sv *Server) removeImageHandler(c *gin.Context) {
	_, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, errors.New("missing user payload in context")))
		return
	}

	var params RemoveImageParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/{publicID} [delete]
func (sv *Server) removeImageByPublicIDHandler(c *gin.Context) {
	_, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(UnauthorizedCode, errors.New("missing user payload in context")))
		return
	}
	var params PublicIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	_, err := sv.removeImageUtil(c, params.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) removeImageUtil(c *gin.Context, publicID string) (msg string, err error) {
	return sv.uploadService.RemoveFile(c, publicID)
}
