package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/technical-test-troben/models"
)

var (
	movieCache = struct {
		sync.Mutex
		cache map[string]interface{}
	}{
		cache: make(map[string]interface{}),
	}
)

func GetDetailById(ctx *gin.Context) {
	id := ctx.Query("id")

	movieCache.Lock()
	cachedMovie, found := movieCache.cache[id]
	movieCache.Unlock()

	if found {
		ctx.JSON(http.StatusOK, models.Response{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       cachedMovie,
		})
		return
	}

	baseURL := "http://www.omdbapi.com/"

	u, err := url.Parse(baseURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error parsing base URL",
		})
		return
	}

	query := u.Query()
	query.Set("apikey", "d0d5937d")

	if id != "" {
		query.Set("i", id)
	}

	u.RawQuery = query.Encode()

	fmt.Println("OMDB API: ", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    http.StatusText(http.StatusNotFound),
			Data:       nil,
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    http.StatusText(http.StatusInternalServerError),
			Data:       nil,
		})
		return
	}

	var movie models.MovieDetailModel
	if err := json.Unmarshal(body, &movie); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNoContent, models.Response{
			StatusCode: http.StatusNoContent,
			Message:    http.StatusText(http.StatusNoContent),
			Data:       nil,
		})
		return
	}

	movieCache.Lock()
	movieCache.cache[id] = movie
	movieCache.Unlock()

	ctx.JSON(http.StatusOK, models.Response{

		StatusCode: http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Data:       movie,
	})
}

func GetMovieBySearch(ctx *gin.Context) {
	title := ctx.Query("title")
	year := ctx.Query("year")
	movieType := ctx.Query("type")

	cacheKey := fmt.Sprintf("title=%s&year=%s&type=%s", title, year, movieType)

	movieCache.Lock()
	cachedMovie, found := movieCache.cache[cacheKey]
	movieCache.Unlock()

	if found {
		ctx.JSON(http.StatusOK, models.Response{
			StatusCode: http.StatusOK,
			Message:    http.StatusText(http.StatusOK),
			Data:       cachedMovie,
		})
		return
	}

	baseURL := "http://www.omdbapi.com/"

	u, err := url.Parse(baseURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error parsing base URL",
		})
		return
	}

	query := u.Query()
	query.Set("apikey", "d0d5937d")

	if title != "" {
		query.Set("s", title)
	}
	if year != "" {
		query.Set("y", year)
	}
	if movieType != "" {
		query.Set("type", movieType)
	}

	u.RawQuery = query.Encode()
	fmt.Println("OMDB API: ", u.String())

	resp, err := http.Get(u.String())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.Response{
			StatusCode: http.StatusNotFound,
			Message:    http.StatusText(http.StatusNotFound),
			Data:       nil,
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    http.StatusText(http.StatusInternalServerError),
			Data:       nil,
		})
		return
	}
	type SearchResult struct {
		Search []models.MovieModel `json:"Search"`
	}
	var movies SearchResult
	if err := json.Unmarshal(body, &movies); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNoContent, models.Response{
			StatusCode: http.StatusNoContent,
			Message:    http.StatusText(http.StatusNoContent),
			Data:       nil,
		})
		return
	}
	
	movieCache.Lock()
	movieCache.cache[cacheKey] = movies.Search
	movieCache.Unlock()

	ctx.JSON(http.StatusOK, models.Response{

		StatusCode: http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Data:       movies.Search,
	})
}
