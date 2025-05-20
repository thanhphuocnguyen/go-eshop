package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func (s *Server) getHomePageHandler(ctx *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Use channels to collect results
	categoriesChan := make(chan []CategoryResponse, 1)
	collectionsChan := make(chan []CategoryResponse, 1)

	// Fetch categories
	go func() {
		defer wg.Done()
		categoryRows, err := s.repo.GetCategories(ctx, repository.GetCategoriesParams{
			Limit:      5,
			Offset:     0,
			Published:  utils.BoolPtr(true),
			Remarkable: utils.BoolPtr(true),
		})

		if err != nil {
			categoriesChan <- []CategoryResponse{}
			return
		}

		categoryModel := make([]CategoryResponse, len(categoryRows))
		for i, category := range categoryRows {
			categoryModel[i] = CategoryResponse{
				ID:          category.ID.String(),
				Name:        category.Name,
				Description: category.Description,
				ImageUrl:    category.ImageUrl,
				Slug:        category.Slug,
			}
		}
		categoriesChan <- categoryModel
	}()

	// Fetch collections
	go func() {
		defer wg.Done()
		collectionRows, err := s.repo.GetCollections(ctx, repository.GetCollectionsParams{
			Limit:      5,
			Offset:     0,
			Published:  utils.BoolPtr(true),
			Remarkable: utils.BoolPtr(true),
		})

		if err != nil {
			collectionsChan <- []CategoryResponse{}
			return
		}

		collectionModel := make([]CategoryResponse, len(collectionRows))
		for i, collection := range collectionRows {
			collectionModel[i] = CategoryResponse{
				ID:          collection.ID.String(),
				Name:        collection.Name,
				Description: collection.Description,
				ImageUrl:    collection.ImageUrl,
				Slug:        collection.Slug,
			}
		}

		collectionsChan <- collectionModel
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Read the results from channels
	categories := <-categoriesChan
	collections := <-collectionsChan

	// Create the response
	response := DashboardData{
		Categories:  categories,
		Collections: collections,
	}

	ctx.JSON(http.StatusOK, createSuccessResponse(ctx, response, "Get homepage data successfully", nil, nil))
}
