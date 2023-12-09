package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type ReqBody struct {
	Question string `json:"question"`
}

func main() {
	client := openai.NewClient("your-token-api")	
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	r.Use(cors.New(config))


	r.POST("/ask", func(c *gin.Context) {
		var req ReqBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	
		resp, err := client.CreateChatCompletion(c,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Temperature: 0.2,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: req.Question,
					},
				},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"answer": resp.Choices[0].Message.Content})
	})

	r.Run(":8080")
}