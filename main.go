package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ai-model-app/handler"
	"github.com/ai-model-app/service"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	fmt.Println("AI MODEL APP")
	godotenv.Load()
	token := os.Getenv("HUGGINGFACE_TOKEN")

	service := service.NewService(token)
	handler := handler.NewHandler(t, service)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		t.ExecuteTemplate(c.Writer, "index.html", nil)
	})

	r.POST("/mask", handler.MaskModel)
	r.GET("/translate", handler.TranslateModel)
	r.GET("/chat", handler.ChatModel)
	r.Run()
}
