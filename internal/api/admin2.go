package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[uuid.UUID]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/orders/{orderId}/cancel [put]
func (s *Server) adminCancelOrder(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.CancelOrderModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order is already cancelled or refunded"))
		return
	}

	userRole := claims["role"].(string)
	userID := uuid.MustParse(claims["userId"].(string))

	if order.UserID != userID && userRole != "admin" {
		RespondForbidden(w, PermissionDeniedCode, errors.New("you do not have permission to access this order"))
		return
	}

	paymentRow, err := s.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// if order status is not pending or user is not admin
	if order.Status != repository.OrderStatusPending || (!errors.Is(err, repository.ErrRecordNotFound) && paymentRow.Status != repository.PaymentStatusPending) {
		RespondBadRequest(w, PermissionDeniedCode, errors.New("order cannot be cancelled"))
		return
	}

	// if order
	cancelOrderTxParams := repository.CancelOrderTxArgs{
		OrderID: uuid.MustParse(id),
		CancelPaymentFromMethod: func(paymentID string, method string) error {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        paymentRow.Amount.Int.Int64(),
			}
			_, err = s.paymentSrv.RefundPayment(c, req, *paymentRow.Gateway)
			return err
		},
	}
	ordId, err := s.repo.CancelOrderTx(c, cancelOrderTxParams)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)
	RespondSuccess(w, ordId)
}

// @Summary Refund order
// @Description Refund order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[repository.GetOrderRow]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/order/{orderId}/refund [put]
func (s *Server) adminRefundOrder(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.RefundOrderModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("order with ID %s not found", id))
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if order.Status != repository.OrderStatusDelivered {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order cannot be refunded"))
		return
	}

	err = s.repo.RefundOrderTx(c, repository.RefundOrderTxArgs{
		OrderID: uuid.MustParse(id),
		RefundPaymentFromMethod: func(paymentID string, method string) (string, error) {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        order.TotalPrice.Int.Int64(),
				Reason:        req.Reason,
			}
			rs, err := s.paymentSrv.RefundPayment(c, req, method)
			return rs.Reason, err
		},
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)

	RespondSuccess(w, order)
}

// adminGetCategories retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-admin-Categories
// @Accept json
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[[]dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/categories [get]
func (s *Server) adminGetCategories(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var query models.PaginationQuery
	query.Page = 1
	query.PageSize = 10
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = int64(p)
		}
	}
	if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			query.PageSize = int64(ps)
		}
	}
	params := repository.GetCategoriesParams{
		Limit:  10,
		Offset: 0,
	}
	params.Offset = (params.Limit) * int64(query.Page-1)
	params.Limit = int64(query.PageSize)

	categories, err := s.repo.GetCategories(c, params)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	count, err := s.repo.CountCategories(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	categoriesResp := make([]dto.AdminCategoryDetail, len(categories))

	for i, category := range categories {
		categoriesResp[i] = dto.AdminCategoryDetail{
			ID:          category.ID.String(),
			Name:        category.Name,
			Slug:        category.Slug,
			Published:   category.Published,
			CreatedAt:   category.CreatedAt.String(),
			Description: category.Description,
			ImageUrl:    category.ImageUrl,
			UpdatedAt:   category.UpdatedAt.String(),
		}
	}
	pagination := dto.CreatePagination(query.Page, query.PageSize, count)
	RespondSuccessWithPagination(w, categoriesResp, pagination)
}

// adminGetCategoryByID retrieves a Category by its ID.
// @Summary Get a Category by ID
// @Description Get a Category by ID
// @ID get-Category-by-id
// @Accept json
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/categories/{id} [get]
func (s *Server) adminGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	category, err := s.repo.GetCategoryByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, InvalidBodyCode, fmt.Errorf("category with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.CategoryDetail{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Published:   category.Published,
		CreatedAt:   category.CreatedAt.String(),
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
	}

	RespondSuccess(w, resp)
}

// adminCreateCategory creates a new Category.
// @Summary Create a new Category
// @Description Create a new Category
// @ID create-Category
// @Accept json
// @Tags Categories
// @Produce json
// @Param request body models.CreateCategoryModel true "Category request"
// @Success 201 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/categories [post]
func (s *Server) adminCreateCategory(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.CreateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	params := repository.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}

	if req.Image != nil {
		imageID, imageURL, err := s.uploadService.Upload(c, req.Image)
		if err != nil {
			RespondInternalServerError(w, UploadFileCode, err)
			return
		}
		params.ImageID = &imageID
		params.ImageUrl = &imageURL
	}

	col, err := s.repo.CreateCategory(c, params)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	resp := dto.CategoryDetail{
		ID:          col.ID.String(),
		Name:        col.Name,
		Slug:        col.Slug,
		Published:   col.Published,
		CreatedAt:   col.CreatedAt.String(),
		Description: col.Description,
		ImageUrl:    col.ImageUrl,
	}

	RespondSuccess(w, resp)
}

// adminUpdateCategory updates a Category.
// @Summary Update a Category
// @Description Update a Category
// @ID update-Category
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Category ID"
// @Param request body models.UpdateCategoryModel true "Category request"
// @Success 200 {object} dto.ApiResponse[repository.Category]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/categories/{id} [put]
func (s *Server) adminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	var req models.UpdateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	category, err := s.repo.GetCategoryByID(c, uuid.MustParse(id))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, InvalidBodyCode, fmt.Errorf("category with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	updateParam := repository.UpdateCategoryParams{
		ID: category.ID,
	}

	if req.Name != nil {
		updateParam.Name = req.Name
	}

	if req.Slug != nil {
		updateParam.Slug = req.Slug
	}

	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}
	var apiErr *dto.ApiError

	imageID, imageURL := "", ""
	if req.Image != nil {
		oldImageID := category.ImageID
		oldImageURL := category.ImageUrl
		// remove old image
		if oldImageID != nil && oldImageURL != nil {
			_, err = s.uploadService.Remove(c, *oldImageID)
			if err != nil {
				apiErr = &dto.ApiError{
					Code:    UploadFileCode,
					Details: err.Error(),
					Stack:   err}
			}
		}
		imageID, imageURL, err = s.uploadService.Upload(c, req.Image)
		if err != nil {
			RespondInternalServerError(w, UploadFileCode, err)
			return
		}
		updateParam.ImageID = &imageID
		updateParam.ImageUrl = &imageURL
	}
	col, err := s.repo.UpdateCategory(c, updateParam)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccessWithError(w, col, apiErr)
}

// adminDeleteCategory delete a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} nil
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/categories/{id} [delete]
func (s *Server) adminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	_, err = s.repo.GetCategoryByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("category with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteCategory(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

// @Summary Create a new Brand
// @Description Create a new Brand
// @Tags admin
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body models.CreateCategoryModel true "Brand request"
// @Success 201 {object} dto.ApiResponse[repository.Brand]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/brands [post]
func (s *Server) adminCreateBrand(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.CreateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	params := repository.CreateBrandParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}
	if req.Image != nil {
		publicID, imgUrl, err := s.uploadService.Upload(c, req.Image)
		if err != nil {
			RespondInternalServerError(w, UploadFileCode, err)
			return
		}
		params.ImageUrl = &imgUrl
		params.ImageID = &publicID
	}

	col, err := s.repo.CreateBrand(c, params)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, col)
}

// @Summary Get a list of brands
// @Description Get a list of brands
// @Tags admin
// @ID get-brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[[]dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/brands [get]
func (s *Server) adminGetBrands(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var queries models.PaginationQuery = GetPaginationQuery(r)

	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:  20,
		Offset: 0,
	}
	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	rows, err := s.repo.GetBrands(c, dbQueries)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cnt, err := s.repo.CountBrands(c)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	data := make([]dto.CategoryDetail, 0, len(rows))

	for _, row := range rows {
		data = append(data, dto.CategoryDetail{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, cnt)

	resp := dto.CreateDataResp(data, pagination, nil)
	RespondSuccess(w, resp)
}

// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/brands/{id} [get]
func (s *Server) adminGetBrandByID(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	result, err := s.repo.GetBrandByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("brand with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	colResp := dto.AdminCategoryDetail{
		ID:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Slug:        result.Slug,
		Published:   result.Published,
		CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
		ImageUrl:    result.ImageUrl,
		UpdatedAt:   result.UpdatedAt.String(),
	}

	RespondSuccess(w, colResp)
}

// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Tags admin
// @Param id path int true "Brand ID"
// @Param request body models.UpdateCategoryModel true "Brand request"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/brands/{id} [put]
func (s *Server) adminUpdateBrand(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.UpdateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	brand, err := s.repo.GetBrandByID(c, uuid.MustParse(id))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("brand with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	updateParam := repository.UpdateBrandWithParams{
		ID:   brand.ID,
		Name: req.Name,
	}

	if req.Image != nil {

		imgID, imgUrl, err := s.uploadService.Upload(c, req.Image)
		if err != nil {
			log.Error().Err(err).Interface("value", req.Image.Header).Msg("error when upload image")
			RespondInternalServerError(w, UploadFileCode, err)
			return
		}
		updateParam.ImageUrl = &imgUrl
		updateParam.ImageID = &imgID
		oldImageID := brand.ImageID
		if oldImageID != nil {
			_, err := s.uploadService.Remove(c, *oldImageID)
			if err != nil {
				log.Error().Err(err).Msg("error when remove old image")
				RespondInternalServerError(w, UploadFileCode, err)
				return
			}
			log.Info().Msgf("old image %s removed", *oldImageID)
		}
	}

	if req.Slug != nil {
		updateParam.Slug = req.Slug
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := s.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, col)
}

// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} dto.ApiResponse[bool]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/brands/{id} [delete]
func (s *Server) adminDeleteBrand(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	_, err = s.repo.GetBrandByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("brand with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteBrand(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Tags admin
// @Produce json
// @Param request body models.CreateCategoryModel true "Collection info"
// @Success 201 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/collections [post]
func (s *Server) adminCreateCollection(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.CreateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	createParams := repository.CreateCollectionParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = req.Description
	}

	if req.Image != nil {
		ID, url, err := s.uploadService.Upload(c, req.Image)
		if err != nil {
			RespondInternalServerError(w, UploadFileCode, err)
			return
		}

		createParams.ImageID = &ID
		createParams.ImageUrl = &url
	}

	col, err := s.repo.CreateCollection(c, createParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "collections-*")

	RespondCreated(w, col)
}

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Collections
// @Accept json
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/collections [get]
func (s *Server) adminGetCollections(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var queries models.PaginationQuery = GetPaginationQuery(r)

	dbQueries := repository.GetCollectionsParams{
		Limit:  20,
		Offset: 0,
	}

	dbQueries.Offset = int64(queries.Page-1) * int64(queries.PageSize)
	dbQueries.Limit = int64(queries.PageSize)
	collectionRows, err := s.repo.GetCollections(c, dbQueries)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cnt, err := s.repo.CountCollections(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccessWithPagination(w, collectionRows, dto.CreatePagination(cnt, queries.Page, queries.PageSize))
}

// @Summary Get a Collection by ID
// @Description Get a Collection by ID
// @ID get-Collection-by-id
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/collections/{id} [get]
func (s *Server) adminGetCollectionByID(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	collection, err := s.repo.GetCollectionByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("collection with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	colResp := dto.CategoryDetail{
		ID:          collection.ID.String(),
		Slug:        collection.Slug,
		Description: collection.Description,
		Published:   collection.Published,
		Name:        collection.Name,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	RespondSuccess(w, colResp)
}

// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body models.CreateCategoryModel true "Collection info"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/collections/{id} [put]
func (s *Server) adminUpdateCollection(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	var req models.UpdateCategoryModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	collection, err := s.repo.GetCollectionByID(c, uuid.MustParse(id))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("collection with ID %s not found", id))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	updateParam := repository.UpdateCollectionWithParams{
		ID: collection.ID,
	}
	if req.Name != nil {
		updateParam.Name = req.Name
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Image != nil {
		oldImageID := collection.ImageID
		oldImageUrl := collection.ImageUrl
		ID, url, err := s.uploadService.Upload(c, req.Image)
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}

		updateParam.ImageUrl = &url
		updateParam.ImageID = &ID

		// Delete old image
		if oldImageID != nil && oldImageUrl != nil {
			if _, err := s.uploadService.Remove(c, *oldImageID); err != nil {
				RespondInternalServerError(w, InternalServerErrorCode, err)
				return
			}
		}
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := s.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, col)
}
