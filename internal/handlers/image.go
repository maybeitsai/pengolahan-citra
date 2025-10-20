package handlers

import (
	"encoding/json"
	"image"
	"image-processor/internal/models"
	"image-processor/internal/services"
	"image-processor/internal/utils"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	"image/png"
	"net/http"
	"strconv"
)

// ImageHandler handles image-related HTTP requests
type ImageHandler struct {
	stateManager      services.StateManager
	imageProcessor    services.ImageProcessorService
	histogramService  services.HistogramService
}

// NewImageHandler creates a new image handler
func NewImageHandler(
	stateManager services.StateManager,
	imageProcessor services.ImageProcessorService,
	histogramService services.HistogramService,
) *ImageHandler {
	return &ImageHandler{
		stateManager:     stateManager,
		imageProcessor:   imageProcessor,
		histogramService: histogramService,
	}
}

// HandleUpload handles the primary image upload
func (h *ImageHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.SendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	file, _, err := r.FormFile("image")
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to read uploaded file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to decode image: " + err.Error(),
		})
		return
	}
	
	h.stateManager.SetOriginalImage(img)
	
	// Convert image to base64 for display
	imageData, err := utils.ImageToBase64(img)
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to encode image: " + err.Error(),
		})
		return
	}
	
	// Generate pixel matrix
	pixelMatrix := utils.GeneratePixelMatrix(img)
	
	utils.SendJSONResponse(w, models.ProcessResponse{
		Success:     true,
		Message:     "Image uploaded successfully",
		ImageData:   imageData,
		PixelMatrix: pixelMatrix,
	})
}

// HandleUploadSecond handles the second image upload
func (h *ImageHandler) HandleUploadSecond(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.SendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	file, _, err := r.FormFile("image")
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to read uploaded file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to decode image: " + err.Error(),
		})
		return
	}
	
	h.stateManager.SetSecondImage(img)
	
	utils.SendJSONResponse(w, models.ProcessResponse{
		Success: true,
		Message: "Second image uploaded successfully",
	})
}

// HandlePixelInfo handles pixel information requests
func (h *ImageHandler) HandlePixelInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		utils.SendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse query parameters
	x := 0
	y := 0
	var err error
	
	if xStr := r.URL.Query().Get("x"); xStr != "" {
		x, err = strconv.Atoi(xStr)
		if err != nil {
			utils.SendJSONResponse(w, models.PixelInfo{
				Success: false,
				Message: "Invalid x coordinate",
			})
			return
		}
	}
	
	if yStr := r.URL.Query().Get("y"); yStr != "" {
		y, err = strconv.Atoi(yStr)
		if err != nil {
			utils.SendJSONResponse(w, models.PixelInfo{
				Success: false,
				Message: "Invalid y coordinate",
			})
			return
		}
	}
	
	// Use current image if processed image is nil
	img := h.stateManager.GetProcessedImage()
	if img == nil {
		img = h.stateManager.GetCurrentImage()
	}
	
	if img == nil {
		utils.SendJSONResponse(w, models.PixelInfo{
			Success: false,
			Message: "No image loaded",
		})
		return
	}
	
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X ||
		y < bounds.Min.Y || y >= bounds.Max.Y {
		utils.SendJSONResponse(w, models.PixelInfo{
			Success: false,
			Message: "Coordinates out of bounds",
		})
		return
	}
	
	c := img.At(x, y)
	red, green, blue, alpha := c.RGBA()
	
	// Convert to 8-bit values
	rVal := uint8(red >> 8)
	gVal := uint8(green >> 8)
	bVal := uint8(blue >> 8)
	
	// Calculate grayscale value
	grayscale := uint8((0.299*float64(rVal) + 0.587*float64(gVal) + 0.114*float64(bVal)))
	
	response := models.PixelInfo{
		Success:   true,
		X:         x,
		Y:         y,
		R:         rVal,
		G:         gVal,
		B:         bVal,
		A:         uint8(alpha >> 8),
		Grayscale: grayscale,
	}
	
	utils.SendJSONResponse(w, response)
}

// HandleProcess handles image processing operations
func (h *ImageHandler) HandleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.SendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request models.ProcessRequest
	
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}
	
	currentImg := h.stateManager.GetCurrentImage()
	if currentImg == nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "No image loaded",
		})
		return
	}
	
	var resultImage image.Image
	var err2 error
	var skipMatrix bool = false // Skip matrix generation for performance
	
	processor := h.stateManager.GetProcessor()
	
	switch request.Operation {
	case "grayscale":
		grayImg := h.imageProcessor.ConvertToGrayscale(currentImg)
		h.stateManager.GetProcessor().UpdateCurrentImage(grayImg) // Update base image
		processor.SetBrightness(0)
		processor.SetZoom(0)
		processor.SetRotation(0)
		resultImage = grayImg
	case "binary":
		binaryImg := h.imageProcessor.ConvertToBinary(currentImg)
		processor.UpdateCurrentImage(binaryImg) // Update base image
		processor.SetBrightness(0)
		processor.SetZoom(0)
		processor.SetRotation(0)
		resultImage = binaryImg
	case "brightness":
		processor.SetBrightness(request.Value)
		resultImage = h.imageProcessor.ApplyAllTransformations(currentImg, processor.BrightnessValue, processor.ZoomValue, processor.RotateValue)
		// Skip matrix generation for real-time operations
		skipMatrix = true
	case "add", "subtract", "multiply":
		secondImg := processor.SecondImage
		if secondImg == nil {
			utils.SendJSONResponse(w, models.ProcessResponse{
				Success: false,
				Message: "Second image required for arithmetic operations",
			})
			return
		}
		resultImage, err2 = h.imageProcessor.ArithmeticOperation(currentImg, secondImg, request.Operation)
	case "and", "or", "xor":
		secondImg := processor.SecondImage
		if secondImg == nil {
			utils.SendJSONResponse(w, models.ProcessResponse{
				Success: false,
				Message: "Second image required for boolean operations",
			})
			return
		}
		resultImage, err2 = h.imageProcessor.BooleanOperation(currentImg, secondImg, request.Operation)
	case "flipH":
		// Apply flip as a transformation, don't permanently modify currentImage
		baseImg := h.imageProcessor.ApplyAllTransformations(currentImg, processor.BrightnessValue, processor.ZoomValue, processor.RotateValue)
		resultImage = h.imageProcessor.FlipHorizontal(baseImg)
	case "flipV":
		// Apply flip as a transformation, don't permanently modify currentImage
		baseImg := h.imageProcessor.ApplyAllTransformations(currentImg, processor.BrightnessValue, processor.ZoomValue, processor.RotateValue)
		resultImage = h.imageProcessor.FlipVertical(baseImg)
	case "zoomSlider":
		processor.SetZoom(request.Value)
		resultImage = h.imageProcessor.ApplyAllTransformations(currentImg, processor.BrightnessValue, processor.ZoomValue, processor.RotateValue)
		// Skip matrix generation for real-time operations
		skipMatrix = true
	case "rotateSlider":
		processor.SetRotation(request.Value)
		resultImage = h.imageProcessor.ApplyAllTransformations(currentImg, processor.BrightnessValue, processor.ZoomValue, processor.RotateValue)
		// Skip matrix generation for real-time operations
		skipMatrix = true
	case "translate":
		sourceImg := h.stateManager.GetProcessedImage()
		if sourceImg == nil {
			sourceImg = currentImg
		}
		resultImage = h.imageProcessor.TranslateImage(sourceImg, 50, 50)
	case "reset":
		// Check if we have an original image
		if processor.OriginalImage == nil {
			utils.SendJSONResponse(w, models.ProcessResponse{
				Success: false,
				Message: "No original image available to reset to",
			})
			return
		}
		// Reset all transformation values
		h.stateManager.Reset()
		resultImage = processor.OriginalImage
	default:
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Unknown operation: " + request.Operation,
		})
		return
	}
	
	if err2 != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Processing failed: " + err2.Error(),
		})
		return
	}
	
	h.stateManager.UpdateProcessedImage(resultImage)
	
	// Convert image to base64
	imageData, err := utils.ImageToBase64(resultImage)
	if err != nil {
		utils.SendJSONResponse(w, models.ProcessResponse{
			Success: false,
			Message: "Failed to encode processed image: " + err.Error(),
		})
		return
	}
	
	// Generate pixel matrix only if needed (skip for real-time operations)
	var pixelMatrix string
	if !skipMatrix {
		pixelMatrix = utils.GeneratePixelMatrix(resultImage)
	}
	
	utils.SendJSONResponse(w, models.ProcessResponse{
		Success:     true,
		Message:     "Image processed successfully",
		ImageData:   imageData,
		PixelMatrix: pixelMatrix,
	})
}

// HandleDownload serves the processed image for download
func (h *ImageHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	processedImg := h.stateManager.GetProcessedImage()
	if processedImg == nil {
		http.Error(w, "No processed image available", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=processed_image.png")
	
	err := png.Encode(w, processedImg)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
	}
}

// HandleHistogram handles histogram calculation and analysis
func (h *ImageHandler) HandleHistogram(w http.ResponseWriter, r *http.Request) {
	currentImg := h.stateManager.GetCurrentImage()
	if currentImg == nil {
		utils.SendJSONResponse(w, models.HistogramResponse{
			Success: false,
			Message: "No image uploaded",
		})
		return
	}
	
	histogramData := h.histogramService.GenerateHistogramData(currentImg)
	
	utils.SendJSONResponse(w, models.HistogramResponse{
		Success: true,
		Data:    *histogramData,
	})
}