package main

import (
	"fmt"
	"log"
	"net/http"

	"image-processor/internal/handlers"
	"image-processor/internal/services"
	"image-processor/internal/templates"

	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
)

// App represents the main application structure
type App struct {
	imageHandler    *handlers.ImageHandler
	templateHandler *templates.TemplateHandler
}

// NewApp creates a new application instance with dependency injection
func NewApp() *App {
	// Initialize services
	stateManager := services.NewStateManager()
	imageProcessor := services.NewImageProcessorService()
	histogramService := services.NewHistogramService()
	edgeDetectionService := services.NewEdgeDetectionService()
	
	// Initialize handlers with dependency injection
	imageHandler := handlers.NewImageHandler(
		stateManager,
		imageProcessor,
		histogramService,
		edgeDetectionService,
	)
	
	templateHandler := templates.NewTemplateHandler()
	
	return &App{
		imageHandler:    imageHandler,
		templateHandler: templateHandler,
	}
}

// setupRoutes configures the HTTP routes
func (app *App) setupRoutes() {
	// Template routes
	http.HandleFunc("/", app.templateHandler.ServeHome)
	
	// API routes
	http.HandleFunc("/upload", app.imageHandler.HandleUpload)
	http.HandleFunc("/upload-second", app.imageHandler.HandleUploadSecond)
	http.HandleFunc("/pixel-info", app.imageHandler.HandlePixelInfo)
	http.HandleFunc("/process", app.imageHandler.HandleProcess)
	http.HandleFunc("/download", app.imageHandler.HandleDownload)
	http.HandleFunc("/histogram", app.imageHandler.HandleHistogram)
	http.HandleFunc("/edge-detection", app.imageHandler.HandleEdgeDetection)
	
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
}

// start starts the HTTP server
func (app *App) start(port string) error {
	app.setupRoutes()
	
	fmt.Println("🖼️  Image Processing Learning Tool")
	fmt.Printf("📡 Server starting on http://localhost%s\n", port)
	fmt.Println("🌐 Open your web browser and navigate to the URL above")
	fmt.Println("⚡ Press Ctrl+C to stop the server")
	
	return http.ListenAndServe(port, nil)
}

func main() {
	// Create application instance
	app := NewApp()
	
	// Start the server
	port := ":8080"
	if err := app.start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}