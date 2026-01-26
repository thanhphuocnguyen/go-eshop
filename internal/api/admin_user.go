package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

// adminGetUsers godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.ApiResponse[[]dto.UserDetail]
// @Failure 500 {object} ErrorResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /admin/users [get]
func (s *Server) adminGetUsers(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	// Parse query parameters
	queries := ParsePaginationQuery(r)

	users, err := s.repo.GetUsers(c, repository.GetUsersParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	total, err := s.repo.CountUsers(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	userResp := make([]dto.UserDetail, 0)
	roleCode := claims["role_code"].(string)
	for _, user := range users {
		userResp = append(userResp, dto.MapToUserResponse(user, roleCode))
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, total)
	RespondSuccessWithPagination(w, userResp, pagination)
}

// adminGetUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} dto.ApiResponse[dto.UserDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/users/{id} [get]
func (s *Server) adminGetUser(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	user, err := s.repo.GetUserByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	roleCode := claims["role_code"].(string)
	userResp := dto.MapToUserResponse(user, roleCode)
	RespondSuccess(w, userResp)
}
