package services

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type imageProcessor struct{}

// NewImageProcessorService creates a new image processor service
func NewImageProcessorService() ImageProcessorService {
	return &imageProcessor{}
}

// ConvertToGrayscale converts an image to grayscale
func (ip *imageProcessor) ConvertToGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	
	// Optimized grayscale conversion
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			
			// Fast integer-based grayscale conversion
			// Using luminance formula with integer arithmetic for speed
			gray := uint8((r*77 + g*151 + b*28) >> 16)
			grayImg.Set(x, y, color.Gray{gray})
		}
	}
	
	return grayImg
}

// ConvertToBinary converts an image to binary
func (ip *imageProcessor) ConvertToBinary(img image.Image) image.Image {
	bounds := img.Bounds()
	binaryImg := image.NewGray(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			
			// Fast luminance calculation with integer arithmetic
			gray := (r*77 + g*151 + b*28) >> 16
			
			// Apply threshold (128) - optimized comparison
			var binaryVal uint8
			if gray > 128 {
				binaryVal = 255 // White
			} else {
				binaryVal = 0 // Black
			}
			binaryImg.Set(x, y, color.Gray{binaryVal})
		}
	}
	
	return binaryImg
}

// AdjustBrightness adjusts the brightness of an image
func (ip *imageProcessor) AdjustBrightness(img image.Image, brightness int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	brightImg := image.NewRGBA(bounds)
	
	// Parallel processing with goroutines for better performance
	numWorkers := 4
	rowsPerWorker := height / numWorkers
	
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	
	for worker := 0; worker < numWorkers; worker++ {
		startY := worker * rowsPerWorker
		endY := startY + rowsPerWorker
		if worker == numWorkers-1 {
			endY = height // Last worker handles remaining rows
		}
		
		go func(startY, endY int) {
			defer wg.Done()
			
			// Pre-calculate brightness lookup table for faster processing
			var brightTable [256]uint8
			for i := 0; i < 256; i++ {
				val := i + brightness
				if val < 0 {
					val = 0
				}
				if val > 255 {
					val = 255
				}
				brightTable[i] = uint8(val)
			}
			
			for y := startY; y < endY; y++ {
				for x := 0; x < width; x++ {
					c := img.At(x+bounds.Min.X, y+bounds.Min.Y)
					r, g, b, a := c.RGBA()
					
					// Fast lookup table brightness adjustment
					r8 := brightTable[r>>8]
					g8 := brightTable[g>>8]
					b8 := brightTable[b>>8]
					
					brightImg.Set(x+bounds.Min.X, y+bounds.Min.Y, color.RGBA{r8, g8, b8, uint8(a >> 8)})
				}
			}
		}(startY, endY)
	}
	
	wg.Wait()
	return brightImg
}

// ResizeImage resizes an image to the specified dimensions
func (ip *imageProcessor) ResizeImage(src image.Image, newWidth, newHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()
	
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	
	// Pre-calculate ratios once
	xRatio := float64(srcWidth) / float64(newWidth)
	yRatio := float64(srcHeight) / float64(newHeight)
	
	// Parallel processing for better performance
	numWorkers := 4
	rowsPerWorker := newHeight / numWorkers
	
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	
	for worker := 0; worker < numWorkers; worker++ {
		startY := worker * rowsPerWorker
		endY := startY + rowsPerWorker
		if worker == numWorkers-1 {
			endY = newHeight
		}
		
		go func(startY, endY int) {
			defer wg.Done()
			
			for y := startY; y < endY; y++ {
				srcY := int(float64(y) * yRatio)
				if srcY >= srcHeight {
					srcY = srcHeight - 1
				}
				
				for x := 0; x < newWidth; x++ {
					srcX := int(float64(x) * xRatio)
					if srcX >= srcWidth {
						srcX = srcWidth - 1
					}
					
					dst.Set(x, y, src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
				}
			}
		}(startY, endY)
	}
	
	wg.Wait()
	return dst
}

// ZoomByFactor zooms an image by the specified factor
func (ip *imageProcessor) ZoomByFactor(src image.Image, zoomFactor float64) image.Image {
	bounds := src.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	
	// Calculate new dimensions based on zoom factor
	newWidth := int(float64(originalWidth) * zoomFactor)
	newHeight := int(float64(originalHeight) * zoomFactor)
	
	// Ensure minimum size to prevent image from becoming too small
	if newWidth < 10 {
		newWidth = 10
	}
	if newHeight < 10 {
		newHeight = 10
	}
	
	// Ensure maximum size to prevent memory issues (max 4000px)
	if newWidth > 4000 {
		ratio := 4000.0 / float64(newWidth)
		newWidth = 4000
		newHeight = int(float64(newHeight) * ratio)
	}
	if newHeight > 4000 {
		ratio := 4000.0 / float64(newHeight)
		newHeight = 4000
		newWidth = int(float64(newWidth) * ratio)
	}
	
	return ip.ResizeImage(src, newWidth, newHeight)
}

// RotateByAngle rotates an image by the specified angle
func (ip *imageProcessor) RotateByAngle(src image.Image, angle float64) image.Image {
	// Normalize angle to 0-360 range
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	
	// For performance, handle common angles with optimized functions
	if angle == 0 {
		return src
	} else if angle == 90 {
		return ip.rotate90(src)
	} else if angle == 180 {
		return ip.rotate180(src)
	} else if angle == 270 {
		return ip.rotate270(src)
	}
	
	// For other angles, use simple step rotation (90-degree increments for now)
	steps := int(angle / 90)
	result := src
	for i := 0; i < steps; i++ {
		result = ip.rotate90(result)
	}
	return result
}

// rotate90 rotates an image 90 degrees clockwise
func (ip *imageProcessor) rotate90(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Rotated image dimensions are swapped
	dst := image.NewRGBA(image.Rect(0, 0, height, width))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Rotate coordinates: (x,y) -> (height-1-y, x)
			newX := height - 1 - (y - bounds.Min.Y)
			newY := x - bounds.Min.X
			dst.Set(newX, newY, src.At(x, y))
		}
	}
	
	return dst
}

// rotate180 rotates an image 180 degrees
func (ip *imageProcessor) rotate180(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Rotate 180°: (x,y) -> (width-1-x, height-1-y)
			newX := width - 1 - (x - bounds.Min.X)
			newY := height - 1 - (y - bounds.Min.Y)
			dst.Set(newX, newY, src.At(x, y))
		}
	}
	
	return dst
}

// rotate270 rotates an image 270 degrees clockwise
func (ip *imageProcessor) rotate270(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Rotated image dimensions are swapped
	dst := image.NewRGBA(image.Rect(0, 0, height, width))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Rotate 270°: (x,y) -> (y, width-1-x)
			newX := y - bounds.Min.Y
			newY := width - 1 - (x - bounds.Min.X)
			dst.Set(newX, newY, src.At(x, y))
		}
	}
	
	return dst
}

// FlipHorizontal flips an image horizontally
func (ip *imageProcessor) FlipHorizontal(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Flip horizontally: (x,y) -> (width-1-x, y)
			newX := width - 1 - (x - bounds.Min.X)
			newY := y - bounds.Min.Y
			dst.Set(newX, newY, src.At(x, y))
		}
	}
	
	return dst
}

// FlipVertical flips an image vertically
func (ip *imageProcessor) FlipVertical(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Flip vertically: (x,y) -> (x, height-1-y)
			newX := x - bounds.Min.X
			newY := height - 1 - (y - bounds.Min.Y)
			dst.Set(newX, newY, src.At(x, y))
		}
	}
	
	return dst
}

// TranslateImage translates an image by the specified offset
func (ip *imageProcessor) TranslateImage(src image.Image, offsetX, offsetY int) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	
	// Fill with white background
	draw.Draw(dst, bounds, &image.Uniform{color.White}, image.Point{}, draw.Src)
	
	// Copy pixels with offset
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newX := x + offsetX
			newY := y + offsetY
			
			if newX >= bounds.Min.X && newX < bounds.Max.X &&
				newY >= bounds.Min.Y && newY < bounds.Max.Y {
				c := src.At(x, y)
				dst.Set(newX, newY, c)
			}
		}
	}
	
	return dst
}

// ArithmeticOperation performs arithmetic operations on two images
func (ip *imageProcessor) ArithmeticOperation(img1, img2 image.Image, operation string) (image.Image, error) {
	bounds := img1.Bounds()
	resizedImg2 := ip.ResizeImage(img2, bounds.Dx(), bounds.Dy())
	
	resultImg := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := resizedImg2.At(x-bounds.Min.X, y-bounds.Min.Y)
			
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			
			// Convert to 8-bit
			r1_8 := int(r1 >> 8)
			g1_8 := int(g1 >> 8)
			b1_8 := int(b1 >> 8)
			r2_8 := int(r2 >> 8)
			g2_8 := int(g2 >> 8)
			b2_8 := int(b2 >> 8)
			
			var rResult, gResult, bResult int
			
			switch operation {
			case "add":
				rResult = r1_8 + r2_8
				gResult = g1_8 + g2_8
				bResult = b1_8 + b2_8
			case "subtract":
				rResult = r1_8 - r2_8
				gResult = g1_8 - g2_8
				bResult = b1_8 - b2_8
			case "multiply":
				rResult = (r1_8 * r2_8) / 255
				gResult = (g1_8 * g2_8) / 255
				bResult = (b1_8 * b2_8) / 255
			}
			
			// Clamp values
			if rResult < 0 {
				rResult = 0
			}
			if rResult > 255 {
				rResult = 255
			}
			if gResult < 0 {
				gResult = 0
			}
			if gResult > 255 {
				gResult = 255
			}
			if bResult < 0 {
				bResult = 0
			}
			if bResult > 255 {
				bResult = 255
			}
			
			resultImg.Set(x, y, color.RGBA{
				uint8(rResult),
				uint8(gResult),
				uint8(bResult),
				uint8((a1 + a2) / 2 >> 8),
			})
		}
	}
	
	return resultImg, nil
}

// BooleanOperation performs boolean operations on two images
func (ip *imageProcessor) BooleanOperation(img1, img2 image.Image, operation string) (image.Image, error) {
	bounds := img1.Bounds()
	resizedImg2 := ip.ResizeImage(img2, bounds.Dx(), bounds.Dy())
	
	resultImg := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := resizedImg2.At(x-bounds.Min.X, y-bounds.Min.Y)
			
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, _ := c2.RGBA()
			
			// Convert to 8-bit
			r1_8 := uint8(r1 >> 8)
			g1_8 := uint8(g1 >> 8)
			b1_8 := uint8(b1 >> 8)
			r2_8 := uint8(r2 >> 8)
			g2_8 := uint8(g2 >> 8)
			b2_8 := uint8(b2 >> 8)
			
			var rResult, gResult, bResult uint8
			
			switch operation {
			case "and":
				rResult = r1_8 & r2_8
				gResult = g1_8 & g2_8
				bResult = b1_8 & b2_8
			case "or":
				rResult = r1_8 | r2_8
				gResult = g1_8 | g2_8
				bResult = b1_8 | b2_8
			case "xor":
				rResult = r1_8 ^ r2_8
				gResult = g1_8 ^ g2_8
				bResult = b1_8 ^ b2_8
			}
			
			resultImg.Set(x, y, color.RGBA{
				rResult,
				gResult,
				bResult,
				uint8(a1 >> 8),
			})
		}
	}
	
	return resultImg, nil
}

// ApplyAllTransformations applies multiple transformations in sequence
func (ip *imageProcessor) ApplyAllTransformations(baseImg image.Image, brightnessVal, zoomVal, rotateVal int) image.Image {
	if baseImg == nil {
		return nil
	}
	
	result := baseImg
	
	// Apply transformations in order: zoom -> rotate -> brightness
	// This order preserves the visual intent best
	
	// 1. Apply zoom if not default (0)
	if zoomVal != 0 {
		zoomValue := float64(zoomVal)
		var zoomFactor float64
		if zoomValue >= 0 {
			zoomFactor = 1.0 + (zoomValue/100.0)*4.0
		} else {
			zoomFactor = 0.1 + ((zoomValue+100.0)/100.0)*0.9
		}
		result = ip.ZoomByFactor(result, zoomFactor)
	}
	
	// 2. Apply rotation if not default (0)
	if rotateVal != 0 {
		angle := float64(rotateVal * 2)
		result = ip.RotateByAngle(result, angle)
	}
	
	// 3. Apply brightness if not default (0)
	if brightnessVal != 0 {
		result = ip.AdjustBrightness(result, brightnessVal)
	}
	
	return result
}