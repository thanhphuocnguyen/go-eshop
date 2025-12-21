package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

// @Summary Get list of product image by ID
// @Schemes http
// @Description get list of product image by ID
// @Tags images
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /images/product/{productId} [get]
func (s *Server) getProductImages(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	fmt.Println(userId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
}

// @Summary Remove a product by external ID
// @Schemes http
// @Description remove a product by external ID
// @Tags images
// @Accept json
// @Param publicID path int true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[bool]
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /images/{id} [delete]
func (s *Server) removeImageByPublicID(w http.ResponseWriter, r *http.Request) {
	publicID, err := GetUrlParam(r, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsoResp, _ := json.Marshal(dto.CreateErr(InvalidBodyCode, err))
		w.Write(jsoResp)
		return
	}
	res, err := s.removeImageUtil(r.Context(), publicID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, res)
}

func (s *Server) removeImageUtil(ctx context.Context, publicID string) (msg string, err error) {
	return s.uploadService.Remove(ctx, publicID)
}
