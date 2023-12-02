package query

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CreateURLRequest struct {
	URL string `json:"url"`
}

func Query(c *gin.Context) {
	var req CreateURLRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	url := req.URL

	log.Println("shortcode creation request for", url)

	response := Response{ShortCode: "Hello"}
	c.JSON(201, response)

}

type Response struct {
	ShortCode string `json:"short_code"`
}
