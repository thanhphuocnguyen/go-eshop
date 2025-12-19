package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

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
func (sv *Server) getProductImages(w http.ResponseWriter, r *http.Request) {
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
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /images/{publicID} [delete]
func (sv *Server) removeImageByPublicID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		err := dto.CreateErr(UnauthorizedCode, errors.New("unauthorized"))
		w.WriteHeader(http.StatusUnauthorized)
		jsoResp, _ := json.Marshal(err)
		w.Write(jsoResp)
		return
	}
	var params models.PublicIDParam
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsoResp, _ := json.Marshal(dto.CreateErr(InvalidBodyCode, err))
		w.Write(jsoResp)
		return
	}
	res, err := sv.removeImageUtil(r.Context(), params.PublicID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsoResp, _ := json.Marshal(dto.CreateErr(InternalServerErrorCode, err))
		w.Write(jsoResp)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	jsoResp, _ := json.Marshal(dto.CreateDataResp(r.Context(), res, nil, nil))
	w.Write(jsoResp)
}

func (sv *Server) removeImageUtil(ctx context.Context, publicID string) (msg string, err error) {
	return sv.uploadService.Remove(ctx, publicID)
}
