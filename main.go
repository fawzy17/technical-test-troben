package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/technical-test-troben/handlers"
)

func init() {

}

func main() {
	r := gin.Default()

	// http://localhost:8001/movie/detail?id=tt3896198
	r.GET("/movie/detail", handlers.GetDetailById)

	// http://localhost:8001/movie/search?title=joker&year=2024&type=(movie, series, episode)
	r.GET("/movie/search", handlers.GetMovieBySearch)



	r.Run(":8001")
	fmt.Println("Running")
}