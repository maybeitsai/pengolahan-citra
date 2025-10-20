package models

import "image"

// ImageProcessor represents the main application structure
type ImageProcessor struct {
	OriginalImage   image.Image // Store the true original image
	CurrentImage    image.Image // Working image that can be modified
	SecondImage     image.Image
	ProcessedImage  image.Image
	BrightnessValue int
	ZoomValue       int
	RotateValue     int
}

// NewImageProcessor creates a new ImageProcessor instance
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		BrightnessValue: 0,
		ZoomValue:       0,
		RotateValue:     0,
	}
}

// SetOriginalImage sets the original image and working copy
func (ip *ImageProcessor) SetOriginalImage(img image.Image) {
	ip.OriginalImage = img
	ip.CurrentImage = img
	ip.ProcessedImage = img
}

// SetSecondImage sets the second image for arithmetic operations
func (ip *ImageProcessor) SetSecondImage(img image.Image) {
	ip.SecondImage = img
}

// UpdateProcessedImage updates the processed image
func (ip *ImageProcessor) UpdateProcessedImage(img image.Image) {
	ip.ProcessedImage = img
}

// UpdateCurrentImage updates the current working image
func (ip *ImageProcessor) UpdateCurrentImage(img image.Image) {
	ip.CurrentImage = img
}

// Reset resets the processor to the original image state
func (ip *ImageProcessor) Reset() {
	if ip.OriginalImage != nil {
		ip.CurrentImage = ip.OriginalImage
		ip.ProcessedImage = ip.OriginalImage
		ip.BrightnessValue = 0
		ip.ZoomValue = 0
		ip.RotateValue = 0
	}
}

// SetBrightness sets the brightness value
func (ip *ImageProcessor) SetBrightness(value int) {
	ip.BrightnessValue = value
}

// SetZoom sets the zoom value
func (ip *ImageProcessor) SetZoom(value int) {
	ip.ZoomValue = value
}

// SetRotation sets the rotation value
func (ip *ImageProcessor) SetRotation(value int) {
	ip.RotateValue = value
}