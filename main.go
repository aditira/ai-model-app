package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ai-model-app/handler"
	"github.com/ai-model-app/service"
)

func main() {
	fmt.Println("AI MODEL APP")
	godotenv.Load()

	// Read token in .env file
	token := os.Getenv("HUGGINGFACE_TOKEN")

	service := service.NewService(token)
	handler := handler.NewHandler(service)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/mask", handler.MaskModel)
	r.POST("/translate", handler.TranslateModel)
	r.Run()
}
