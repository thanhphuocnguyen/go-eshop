package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"golang.org/x/sync/errgroup"
)

// ------------------------------ API Models ------------------------------
type Attribute struct {
	ID              int32    `json:"id"`
	Name            string   `json:"name"`
	AttributeValues []string `json:"values,omitempty"`
}

type AttributeRequest struct {
	Name string `json:"name" binding:"required"`
}

type AttributeParam struct {
	ID int32 `uri:"id" binding:"required"`
}

type ProductAttributeDetail struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ------------------------------ API Handlers ------------------------------

// @Summary Create an attribute
// @Description Create an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param params body AttributeRequest true "Attribute name"
// @Success 201 {object} GenericResponse[Attribute]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes [post]
func (sv *Server) createAttribute(c *gin.Context) {
	var params AttributeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := "Attribute created successfully"
	c.JSON(http.StatusCreated, GenericResponse[repository.Attribute]{&attribute, &message, nil})
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} GenericResponse[Attribute]
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes/{id} [get]
func (sv *Server) getAttributeByID(c *gin.Context) {
	var attributeParam AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	attributeRows, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Attribute]{&attributeRows, nil, nil})
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} GenericListResponse[Attribute]
// @Failure 500 {object} errorResponse
// @Router /attributes [get]
func (sv *Server) getAttributes(c *gin.Context) {
	errgroup, ctx := errgroup.WithContext(c)
	attributeRowsChan := make(chan []repository.Attribute, 1)
	attributeCntChan := make(chan int64, 1)
	defer close(attributeRowsChan)
	defer close(attributeCntChan)

	errgroup.Go(func() error {
		attributeRows, err := sv.repo.GetAttributes(ctx)
		if err != nil {
			return err
		}
		attributeRowsChan <- attributeRows
		return nil
	})
	errgroup.Go(func() error {
		cnt, err := sv.repo.CountAttributes(ctx)
		if err != nil {
			return err
		}
		attributeCntChan <- cnt
		return nil
	})

	if err := errgroup.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	resp := <-attributeRowsChan

	c.JSON(http.StatusOK, GenericListResponse[repository.Attribute]{&resp, <-attributeCntChan, nil, nil})
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeRequest true "Attribute name"
// @Success 200 {object} GenericResponse[Attribute]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes/{id} [put]
func (sv *Server) updateAttribute(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var req AttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	existingAttribute, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("attribute not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if existingAttribute.Name == req.Name {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("attribute name is the same as the existing one")))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		AttributeID: existingAttribute.AttributeID,
		Name:        req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := "Attribute updated successfully"
	c.JSON(http.StatusOK, GenericResponse[repository.Attribute]{&attribute, &message, nil})
}

// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 500 {object} errorResponse
// @Router /attributes/{id} [delete]
func (sv *Server) deleteAttribute(c *gin.Context) {
	var params AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	attribute, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attribute.AttributeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
