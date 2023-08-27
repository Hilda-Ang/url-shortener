package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShortenUrlRequest struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type ShortenUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
	Error    string `json:"error"`
}

func RedirectToLongUrl(c *gin.Context) {
	path := c.Param("path")

	url := GetLongUrl(c.MustGet("db").(*mongo.Collection), path)

	if url == "" {
		c.IndentedJSON(http.StatusNotFound, "No mapping found.")
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GetShortUrl(c *gin.Context) {
	var request ShortenUrlRequest

	if err := c.BindJSON(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	shortUrl := request.ShortUrl

	if request.LongUrl == "" {
		c.IndentedJSON(http.StatusBadRequest, ShortenUrlResponse{Error: "Missing long URL in request."})
		return
	}

	if shortUrl == "" {
		shortUrl = GenerateHash(request.LongUrl)
	}

	if CheckUrl(c.MustGet("db").(*mongo.Collection), shortUrl) {
		c.IndentedJSON(http.StatusNotFound, ShortenUrlResponse{Error: "Short URL has been used before."})
		return
	}

	err := StoreUrl(c.MustGet("db").(*mongo.Collection), shortUrl, request.LongUrl)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, ShortenUrlResponse{Error: "Failed to store URL."})
		return
	}

	c.IndentedJSON(http.StatusOK, ShortenUrlResponse{ShortUrl: fmt.Sprintf("%s/%s", ExternalAddress, shortUrl)})
}
