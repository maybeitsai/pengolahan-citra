package services

import (
	"image"
	"image-processor/internal/models"
	"image-processor/internal/utils"
	"image/color"
	"math"
)

type histogramService struct{}

// NewHistogramService creates a new histogram service
func NewHistogramService() HistogramService {
	return &histogramService{}
}

// CalculateHistograms computes RGB and grayscale histograms
func (hs *histogramService) CalculateHistograms(img image.Image) ([]int, []int, []int, []int) {
	bounds := img.Bounds()
	
	// Initialize histogram arrays for 256 intensity levels (0-255)
	redHist := make([]int, 256)
	greenHist := make([]int, 256)
	blueHist := make([]int, 256)
	grayHist := make([]int, 256)
	
	// Iterate through all pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Convert from 16-bit to 8-bit values
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			// Increment histogram bins
			redHist[r8]++
			greenHist[g8]++
			blueHist[b8]++
			
			// Calculate grayscale value using standard luminance formula
			gray := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))
			grayHist[gray]++
		}
	}
	
	return redHist, greenHist, blueHist, grayHist
}

// CalculateOtsuThreshold implements Otsu's method for automatic threshold selection
func (hs *histogramService) CalculateOtsuThreshold(histogram []int) int {
	total := 0
	for _, count := range histogram {
		total += count
	}
	
	if total == 0 {
		return 128 // Default threshold if no pixels
	}
	
	sum := 0.0
	for i, count := range histogram {
		sum += float64(i * count)
	}
	
	sumB := 0.0
	wB := 0
	wF := 0
	varMax := 0.0
	threshold := 0
	
	for t := 0; t < 256; t++ {
		wB += histogram[t] // Weight Background
		if wB == 0 {
			continue
		}
		
		wF = total - wB // Weight Foreground
		if wF == 0 {
			break
		}
		
		sumB += float64(t * histogram[t])
		
		mB := sumB / float64(wB)             // Mean Background
		mF := (sum - sumB) / float64(wF)     // Mean Foreground
		
		// Calculate Between Class Variance
		varBetween := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
		
		// Check if new maximum found
		if varBetween > varMax {
			varMax = varBetween
			threshold = t
		}
	}
	
	return threshold
}

// ApplyThreshold creates a binary image using the specified threshold
func (hs *histogramService) ApplyThreshold(img image.Image, threshold int) image.Image {
	bounds := img.Bounds()
	binaryImg := image.NewGray(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			// Convert to grayscale
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			
			// Apply threshold
			if int(gray) >= threshold {
				binaryImg.SetGray(x, y, color.Gray{255}) // White
			} else {
				binaryImg.SetGray(x, y, color.Gray{0}) // Black
			}
		}
	}
	
	return binaryImg
}

// CalculateImageStats computes mean and standard deviation of image intensities
func (hs *histogramService) CalculateImageStats(img image.Image) (float64, float64) {
	bounds := img.Bounds()
	totalPixels := 0
	sum := 0.0
	
	// Calculate mean
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			sum += gray
			totalPixels++
		}
	}
	
	if totalPixels == 0 {
		return 0, 0
	}
	
	mean := sum / float64(totalPixels)
	
	// Calculate standard deviation
	sumSquared := 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			diff := gray - mean
			sumSquared += diff * diff
		}
	}
	
	variance := sumSquared / float64(totalPixels)
	stdDev := math.Sqrt(variance)
	
	return mean, stdDev
}

// HistogramEqualization performs uniform histogram equalization
func (hs *histogramService) HistogramEqualization(img image.Image) image.Image {
	bounds := img.Bounds()
	
	// Calculate histogram first
	_, _, _, grayHist := hs.CalculateHistograms(img)
	
	// Calculate total number of pixels
	totalPixels := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	
	// Calculate cumulative distribution function (CDF)
	cdf := make([]float64, 256)
	cdf[0] = float64(grayHist[0])
	for i := 1; i < 256; i++ {
		cdf[i] = cdf[i-1] + float64(grayHist[i])
	}
	
	// Normalize CDF to range [0, 255]
	lookupTable := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		lookupTable[i] = uint8(math.Round((cdf[i] / float64(totalPixels)) * 255.0))
	}
	
	// Create equalized image
	equalizedImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
			// Convert to grayscale and apply equalization
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			equalizedGray := lookupTable[gray]
			
			// Set the equalized grayscale value to all RGB channels
			equalizedImg.Set(x, y, color.RGBA{
				R: equalizedGray,
				G: equalizedGray,
				B: equalizedGray,
				A: uint8(a >> 8),
			})
		}
	}
	
	return equalizedImg
}

// GenerateHistogramData generates complete histogram data for analysis
func (hs *histogramService) GenerateHistogramData(img image.Image) *models.HistogramData {
	// Calculate histograms
	redHist, greenHist, blueHist, grayHist := hs.CalculateHistograms(img)
	
	// Calculate threshold using Otsu's method
	threshold := hs.CalculateOtsuThreshold(grayHist)
	
	// Generate binary image
	binaryImage := hs.ApplyThreshold(img, threshold)
	binaryImageData, _ := utils.ImageToBase64(binaryImage)
	
	// Calculate statistics before equalization
	meanBefore, stdDevBefore := hs.CalculateImageStats(img)
	
	// Perform histogram equalization
	equalizedImage := hs.HistogramEqualization(img)
	equalizedImageData, _ := utils.ImageToBase64(equalizedImage)
	
	// Calculate statistics after equalization
	meanAfter, stdDevAfter := hs.CalculateImageStats(equalizedImage)
	
	return &models.HistogramData{
		RedHistogram:       redHist,
		GreenHistogram:     greenHist,
		BlueHistogram:      blueHist,
		GrayscaleHistogram: grayHist,
		Threshold:          threshold,
		BinaryImageData:    binaryImageData,
		MeanBefore:         meanBefore,
		StdDevBefore:       stdDevBefore,
		MeanAfter:          meanAfter,
		StdDevAfter:        stdDevAfter,
		EqualizedImageData: equalizedImageData,
	}
}