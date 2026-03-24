package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"popfolio/internal/routers"
	"popfolio/internal/storage"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*.html
var tmplFS embed.FS

//go:embed static/*
var staticFS embed.FS

//go:embed data/*
var dataFS embed.FS

func setupRouter() *gin.Engine {
	storage.DataFS = dataFS
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	templ := template.Must(template.ParseFS(tmplFS, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	subFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	httpFS := http.FS(subFS)
	router.StaticFS("/static", httpFS)

	routers.SetupRoutes(router)
	return router
}

// Handler is the Vercel entrypoint for Go serverless functions
func Handler() {
	router := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server running on http://localhost:" + port)
	router.Run(":" + port)
}

func main() {
	router := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server running on http://localhost:" + port)
	router.Run(":" + port)
}
