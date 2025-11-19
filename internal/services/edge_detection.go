package services

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"image-processor/internal/models"
	"image-processor/internal/utils"
)

// edgeDetectionService implements the EdgeDetectionService interface
type edgeDetectionService struct{}

// NewEdgeDetectionService creates a new edge detection service
func NewEdgeDetectionService() EdgeDetectionService {
	return &edgeDetectionService{}
}

// GetSobelKernels returns the Sobel X and Y kernels
func (eds *edgeDetectionService) GetSobelKernels() ([][]int, [][]int) {
	sobelX := [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	
	sobelY := [][]int{
		{-1, -2, -1},
		{0,   0,  0},
		{1,   2,  1},
	}
	
	return sobelX, sobelY
}

// ApplyConvolution applies a convolution kernel to an image
func (eds *edgeDetectionService) ApplyConvolution(img image.Image, kernel [][]int) (*models.ConvolutionResult, error) {
	if img == nil {
		return nil, fmt.Errorf("input image is nil")
	}
	
	// Convert to grayscale first
	grayImg := eds.convertToGrayscale(img)
	bounds := grayImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Create result image
	resultImg := image.NewGray(bounds)
	
	// Initialize convolution data matrix
	convolutionData := make([][]float64, height)
	for i := range convolutionData {
		convolutionData[i] = make([]float64, width)
	}
	
	kernelSize := len(kernel)
	kernelRadius := kernelSize / 2
	
	// Apply convolution
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sum float64 = 0
			
			// Apply kernel
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					// Calculate pixel coordinates
					pixelY := y - kernelRadius + ky
					pixelX := x - kernelRadius + kx
					
					// Handle edge cases with padding (replicate border)
					if pixelX < 0 {
						pixelX = 0
					}
					if pixelX >= width {
						pixelX = width - 1
					}
					if pixelY < 0 {
						pixelY = 0
					}
					if pixelY >= height {
						pixelY = height - 1
					}
					
					// Get pixel value
					grayColor := grayImg.GrayAt(pixelX+bounds.Min.X, pixelY+bounds.Min.Y)
					pixelValue := float64(grayColor.Y)
					
					// Apply kernel weight
					kernelValue := float64(kernel[ky][kx])
					sum += pixelValue * kernelValue
				}
			}
			
			// Store convolution result
			convolutionData[y][x] = sum
			
			// Normalize and clamp for display
			normalizedValue := math.Abs(sum)
			if normalizedValue > 255 {
				normalizedValue = 255
			}
			
			resultImg.SetGray(x+bounds.Min.X, y+bounds.Min.Y, color.Gray{uint8(normalizedValue)})
		}
	}
	
	// Convert result image to base64
	imageData, err := utils.ImageToBase64(resultImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode result image: %v", err)
	}
	
	// Determine kernel name
	kernelName := "Custom"
	if eds.isSobelX(kernel) {
		kernelName = "Sobel X"
	} else if eds.isSobelY(kernel) {
		kernelName = "Sobel Y"
	}
	
	return &models.ConvolutionResult{
		ResultImage:     imageData,
		Kernel:          kernel,
		KernelName:      kernelName,
		ConvolutionData: convolutionData,
		ImageWidth:      width,
		ImageHeight:     height,
	}, nil
}

// ApplySobelX applies Sobel X kernel for vertical edge detection
func (eds *edgeDetectionService) ApplySobelX(img image.Image) (*models.EdgeDetectionResult, error) {
	sobelX, _ := eds.GetSobelKernels()
	
	result, err := eds.ApplyConvolution(img, sobelX)
	if err != nil {
		return nil, fmt.Errorf("failed to apply Sobel X: %v", err)
	}
	
	return &models.EdgeDetectionResult{
		EdgeImage:    result.ResultImage,
		SobelXResult: result,
		SobelXKernel: sobelX,
	}, nil
}

// ApplySobelY applies Sobel Y kernel for horizontal edge detection
func (eds *edgeDetectionService) ApplySobelY(img image.Image) (*models.EdgeDetectionResult, error) {
	_, sobelY := eds.GetSobelKernels()
	
	result, err := eds.ApplyConvolution(img, sobelY)
	if err != nil {
		return nil, fmt.Errorf("failed to apply Sobel Y: %v", err)
	}
	
	return &models.EdgeDetectionResult{
		EdgeImage:    result.ResultImage,
		SobelYResult: result,
		SobelYKernel: sobelY,
	}, nil
}

// ApplySobelEdgeDetection applies both Sobel X and Y kernels and combines them
func (eds *edgeDetectionService) ApplySobelEdgeDetection(img image.Image) (*models.EdgeDetectionResult, error) {
	if img == nil {
		return nil, fmt.Errorf("input image is nil")
	}
	
	sobelX, sobelY := eds.GetSobelKernels()
	
	// Apply both Sobel kernels
	resultX, err := eds.ApplyConvolution(img, sobelX)
	if err != nil {
		return nil, fmt.Errorf("failed to apply Sobel X: %v", err)
	}
	
	resultY, err := eds.ApplyConvolution(img, sobelY)
	if err != nil {
		return nil, fmt.Errorf("failed to apply Sobel Y: %v", err)
	}
	
	// Calculate magnitude and gradient
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	magnitudeImg := image.NewGray(bounds)
	magnitudeData := make([][]float64, height)
	gradientData := make([][]float64, height)
	
	for i := range magnitudeData {
		magnitudeData[i] = make([]float64, width)
		gradientData[i] = make([]float64, width)
	}
	
	// Combine X and Y gradients using magnitude calculation
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gx := resultX.ConvolutionData[y][x]
			gy := resultY.ConvolutionData[y][x]
			
			// Calculate magnitude: sqrt(gx² + gy²)
			magnitude := math.Sqrt(gx*gx + gy*gy)
			magnitudeData[y][x] = magnitude
			
			// Calculate gradient direction (angle in degrees)
			gradient := math.Atan2(gy, gx) * 180 / math.Pi
			gradientData[y][x] = gradient
			
			// Normalize magnitude for display
			normalizedMagnitude := magnitude
			if normalizedMagnitude > 255 {
				normalizedMagnitude = 255
			}
			
			magnitudeImg.SetGray(x+bounds.Min.X, y+bounds.Min.Y, color.Gray{uint8(normalizedMagnitude)})
		}
	}
	
	// Convert magnitude image to base64
	magnitudeImageData, err := utils.ImageToBase64(magnitudeImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode magnitude image: %v", err)
	}
	
	return &models.EdgeDetectionResult{
		EdgeImage:       magnitudeImageData,
		MagnitudeImage:  magnitudeImageData,
		SobelXResult:    resultX,
		SobelYResult:    resultY,
		SobelXKernel:    sobelX,
		SobelYKernel:    sobelY,
		MagnitudeData:   magnitudeData,
		GradientData:    gradientData,
	}, nil
}

// Helper function to convert image to grayscale
func (eds *edgeDetectionService) convertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			
			// Luminance formula with integer arithmetic for speed
			gray := uint8((r*77 + g*151 + b*28) >> 16)
			grayImg.Set(x, y, color.Gray{gray})
		}
	}
	
	return grayImg
}

// Helper function to check if kernel is Sobel X
func (eds *edgeDetectionService) isSobelX(kernel [][]int) bool {
	sobelX, _ := eds.GetSobelKernels()
	if len(kernel) != len(sobelX) || len(kernel[0]) != len(sobelX[0]) {
		return false
	}
	
	for i := range kernel {
		for j := range kernel[i] {
			if kernel[i][j] != sobelX[i][j] {
				return false
			}
		}
	}
	return true
}

// Helper function to check if kernel is Sobel Y
func (eds *edgeDetectionService) isSobelY(kernel [][]int) bool {
	_, sobelY := eds.GetSobelKernels()
	if len(kernel) != len(sobelY) || len(kernel[0]) != len(sobelY[0]) {
		return false
	}
	
	for i := range kernel {
		for j := range kernel[i] {
			if kernel[i][j] != sobelY[i][j] {
				return false
			}
		}
	}
	return true
}