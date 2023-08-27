package main

import (
	"github.com/gin-gonic/gin"

	"example/url-shortener/utils"
)

func main() {
	db, err := utils.ConnectDb()

	if err != nil {
		return
	}

	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	router.Use(utils.AttachDbCollection(db))

	router.GET("/:path", utils.RedirectToLongUrl)
	router.POST("/shorten", utils.GetShortUrl)

	router.Run(utils.InternalAddress)
}
