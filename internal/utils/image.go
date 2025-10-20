package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
	"math"
)

// ImageToBase64 converts an image to base64 string
func ImageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// GeneratePixelMatrix creates a pixel matrix from an image
func GeneratePixelMatrix(img image.Image) string {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	
	// Optimized matrix generation with intelligent size limits
	// For UI display, we don't need massive matrices that slow down rendering
	maxDisplaySize := 100 // Much more reasonable for UI performance
	
	var matrixWidth, matrixHeight int
	var sampledImg image.Image = img
	
	if width > maxDisplaySize || height > maxDisplaySize {
		// Calculate optimal sampling ratio to maintain aspect ratio
		ratio := math.Min(float64(maxDisplaySize)/float64(width), float64(maxDisplaySize)/float64(height))
		matrixWidth = int(float64(width) * ratio)
		matrixHeight = int(float64(height) * ratio)
		
		// Use intelligent sampling instead of full resize for better performance
		sampledImg = CreateSampledMatrix(img, matrixWidth, matrixHeight)
	} else {
		matrixWidth = width
		matrixHeight = height
	}
	
	// Pre-allocate matrix for better memory efficiency
	matrix := make([][]int, matrixHeight)
	for i := range matrix {
		matrix[i] = make([]int, matrixWidth)
	}
	
	// Fast matrix generation with optimized color conversion
	sampledBounds := sampledImg.Bounds()
	for y := 0; y < matrixHeight; y++ {
		for x := 0; x < matrixWidth; x++ {
			c := sampledImg.At(sampledBounds.Min.X+x, sampledBounds.Min.Y+y)
			
			// Faster grayscale conversion using integer arithmetic
			r, g, b, _ := c.RGBA()
			// Use bit shifting for faster conversion and integer weights
			grayscale := int((r*77 + g*151 + b*28) >> 16) // Equivalent to 0.299, 0.587, 0.114 but faster
			matrix[y][x] = grayscale
		}
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(matrix)
	if err != nil {
		return "[]"
	}
	
	return string(jsonData)
}

// CreateSampledMatrix creates a representative matrix without full resize
func CreateSampledMatrix(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	
	// Calculate step size for sampling
	stepX := float64(srcWidth) / float64(targetWidth)
	stepY := float64(srcHeight) / float64(targetHeight)
	
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Sample from original image using nearest neighbor (fastest method)
			srcX := int(float64(x) * stepX)
			srcY := int(float64(y) * stepY)
			
			// Clamp to bounds
			if srcX >= srcWidth {
				srcX = srcWidth - 1
			}
			if srcY >= srcHeight {
				srcY = srcHeight - 1
			}
			
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}
	
	return dst
}

// ClampInt clamps an integer value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampUint8 clamps a uint8 value
func ClampUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}