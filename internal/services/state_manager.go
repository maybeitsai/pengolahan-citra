package services

import (
	"image"
	"image-processor/internal/models"
)

type stateManager struct {
	processor *models.ImageProcessor
}

// NewStateManager creates a new state manager
func NewStateManager() StateManager {
	return &stateManager{
		processor: models.NewImageProcessor(),
	}
}

// GetProcessor returns the image processor instance
func (sm *stateManager) GetProcessor() *models.ImageProcessor {
	return sm.processor
}

// SetOriginalImage sets the original image and working copy
func (sm *stateManager) SetOriginalImage(img image.Image) {
	sm.processor.SetOriginalImage(img)
}

// SetSecondImage sets the second image for arithmetic operations
func (sm *stateManager) SetSecondImage(img image.Image) {
	sm.processor.SetSecondImage(img)
}

// UpdateProcessedImage updates the processed image
func (sm *stateManager) UpdateProcessedImage(img image.Image) {
	sm.processor.UpdateProcessedImage(img)
}

// Reset resets the processor to original state
func (sm *stateManager) Reset() {
	sm.processor.Reset()
}

// GetCurrentImage returns the current working image
func (sm *stateManager) GetCurrentImage() image.Image {
	return sm.processor.CurrentImage
}

// GetProcessedImage returns the processed image
func (sm *stateManager) GetProcessedImage() image.Image {
	return sm.processor.ProcessedImage
}