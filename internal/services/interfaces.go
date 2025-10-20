package services

import (
	"image"
	"image-processor/internal/models"
)

// ImageProcessorService defines the interface for image processing operations
type ImageProcessorService interface {
	// Basic image operations
	ConvertToGrayscale(img image.Image) image.Image
	ConvertToBinary(img image.Image) image.Image
	AdjustBrightness(img image.Image, brightness int) image.Image
	
	// Geometric transformations
	ResizeImage(src image.Image, newWidth, newHeight int) image.Image
	ZoomByFactor(src image.Image, zoomFactor float64) image.Image
	RotateByAngle(src image.Image, angle float64) image.Image
	FlipHorizontal(src image.Image) image.Image
	FlipVertical(src image.Image) image.Image
	TranslateImage(src image.Image, offsetX, offsetY int) image.Image
	
	// Advanced operations
	ArithmeticOperation(img1, img2 image.Image, operation string) (image.Image, error)
	BooleanOperation(img1, img2 image.Image, operation string) (image.Image, error)
	ApplyAllTransformations(baseImg image.Image, brightnessVal, zoomVal, rotateVal int) image.Image
}

// HistogramService defines the interface for histogram operations
type HistogramService interface {
	CalculateHistograms(img image.Image) ([]int, []int, []int, []int)
	CalculateOtsuThreshold(histogram []int) int
	ApplyThreshold(img image.Image, threshold int) image.Image
	CalculateImageStats(img image.Image) (float64, float64)
	HistogramEqualization(img image.Image) image.Image
	GenerateHistogramData(img image.Image) *models.HistogramData
}

// StateManager defines the interface for managing image processor state
type StateManager interface {
	GetProcessor() *models.ImageProcessor
	SetOriginalImage(img image.Image)
	SetSecondImage(img image.Image)
	UpdateProcessedImage(img image.Image)
	Reset()
	GetCurrentImage() image.Image
	GetProcessedImage() image.Image
}