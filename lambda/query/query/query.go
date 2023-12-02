package query

import (
	"github.com/gin-gonic/gin"
	bedrock "github.com/megaproaktiv/go-rag-kendra-bedrock/bedrock"
	"log"
	"net/http"
)

type QueryRequest struct {
	Question string `json:"question"`
}

type Response struct {
	Answer string `json:"answer"`
}

func Query(c *gin.Context) {
	var req QueryRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	question := req.Question

	log.Println("Question received", question)

	answer := bedrock.Chat(question)

	response := Response{Answer: answer}
	c.JSON(201, response)

}
