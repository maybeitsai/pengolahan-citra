package models

// PixelInfo represents pixel information for JSON response
type PixelInfo struct {
	X         int     `json:"x"`
	Y         int     `json:"y"`
	R         uint8   `json:"r"`
	G         uint8   `json:"g"`
	B         uint8   `json:"b"`
	A         uint8   `json:"a"`
	Grayscale uint8   `json:"grayscale"`
	Success   bool    `json:"success"`
	Message   string  `json:"message,omitempty"`
}

// ProcessRequest represents a processing request
type ProcessRequest struct {
	Operation string `json:"operation"`
	Value     int    `json:"value"`
}

// ProcessResponse represents the response after processing an image
type ProcessResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ImageData   string `json:"imageData"`
	PixelMatrix string `json:"pixelMatrix"`
}

// HistogramData represents histogram information
type HistogramData struct {
	RedHistogram       []int   `json:"redHistogram"`
	GreenHistogram     []int   `json:"greenHistogram"`
	BlueHistogram      []int   `json:"blueHistogram"`
	GrayscaleHistogram []int   `json:"grayscaleHistogram"`
	Threshold          int     `json:"threshold"`
	BinaryImageData    string  `json:"binaryImageData"`
	MeanBefore         float64 `json:"meanBefore"`
	StdDevBefore       float64 `json:"stdDevBefore"`
	MeanAfter          float64 `json:"meanAfter"`
	StdDevAfter        float64 `json:"stdDevAfter"`
	EqualizedImageData string  `json:"equalizedImageData"`
}

// HistogramResponse represents the histogram analysis response
type HistogramResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Data    HistogramData `json:"data,omitempty"`
}