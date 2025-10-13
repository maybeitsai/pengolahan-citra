package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	"image/png"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
)

// ImageProcessor represents the main application structure
type ImageProcessor struct {
	originalImage    image.Image // Store the true original image
	currentImage     image.Image // Working image that can be modified
	secondImage      image.Image
	processedImage   image.Image
	brightnessValue  int
	zoomValue        int
	rotateValue      int
}

// PixelInfo represents pixel information for JSON response
type PixelInfo struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	R uint8  `json:"r"`
	G uint8  `json:"g"`
	B uint8  `json:"b"`
	A uint8  `json:"a"`
}

// ProcessResponse represents the response after processing an image
type ProcessResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ImageData   string `json:"imageData"`
	PixelMatrix string `json:"pixelMatrix"`
}

var processor = &ImageProcessor{}

func main() {
	// Serve static files and handle API endpoints
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/upload-second", handleUploadSecond)
	http.HandleFunc("/pixel-info", handlePixelInfo)
	http.HandleFunc("/process", handleProcess)
	http.HandleFunc("/download", handleDownload)
	
	fmt.Println("🖼️  Image Processing Learning Tool")
	fmt.Println("📡 Server starting on http://localhost:8080")
	fmt.Println("🌐 Open your web browser and navigate to the URL above")
	fmt.Println("⚡ Press Ctrl+C to stop the server")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// UI represents the user interface module
type UI struct {
	template *template.Template
}

// NewUI creates a new UI instance
func NewUI() *UI {
	return &UI{}
}

// serveHome serves the main HTML page with modern UI/UX
func serveHome(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Image Processing Learning Tool</title>
    <style>
        :root {
            --primary-color: #3b82f6;
            --primary-hover: #2563eb;
            --success-color: #10b981;
            --error-color: #ef4444;
            --warning-color: #f59e0b;
            --bg-primary: #ffffff;
            --bg-secondary: #f8fafc;
            --bg-tertiary: #f1f5f9;
            --text-primary: #0f172a;
            --text-secondary: #64748b;
            --border-color: #e2e8f0;
            --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
            --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
            --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
            --border-radius: 8px;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Inter', 'Segoe UI', system-ui, sans-serif;
            background: var(--bg-secondary);
            color: var(--text-primary);
            line-height: 1.6;
        }
        
        .app-container {
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }

        .header {
            background: var(--bg-primary);
            border-bottom: 1px solid var(--border-color);
            padding: 1rem 2rem;
            box-shadow: var(--shadow-sm);
        }
        
        .header-content {
            max-width: 1400px;
            margin: 0 auto;
            display: flex;
            align-items: center;
            justify-content: between;
        }

        .header h1 {
            font-size: 1.5rem;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .header p {
            color: var(--text-secondary);
            font-size: 0.875rem;
            margin-top: 0.25rem;
        }
        
        .main-content {
            flex: 1;
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
            display: grid;
            grid-template-columns: 1fr 400px;
            gap: 2rem;
        }

        @media (max-width: 1024px) {
            .main-content {
                grid-template-columns: 1fr;
                gap: 1.5rem;
            }
        }
        
        .card {
            background: var(--bg-primary);
            border: 1px solid var(--border-color);
            border-radius: var(--border-radius);
            padding: 1.5rem;
            margin-bottom: 1.5rem;
            box-shadow: var(--shadow-sm);
            transition: all 0.2s ease;
        }

        .card:hover {
            box-shadow: var(--shadow-md);
        }
        
        .card h3 {
            color: var(--text-primary);
            margin-bottom: 1rem;
            font-size: 1rem;
            font-weight: 600;
        }
        
        .upload-area {
            border: 2px dashed var(--border-color);
            border-radius: var(--border-radius);
            padding: 2rem;
            text-align: center;
            transition: all 0.2s ease;
            cursor: pointer;
            background: var(--bg-tertiary);
        }
        
        .upload-area:hover {
            border-color: var(--primary-color);
            background: var(--bg-primary);
        }
        
        .upload-area.dragover {
            border-color: var(--success-color);
            background: var(--bg-primary);
            transform: scale(1.02);
        }
        
        .canvas-container {
            position: relative;
            display: inline-block;
            max-width: 100%;
        }

        #imageCanvas {
            max-width: 100%;
            border: 1px solid var(--border-color);
            border-radius: var(--border-radius);
            cursor: crosshair;
            display: block;
            box-shadow: var(--shadow-sm);
        }
        
        .btn {
            background: var(--primary-color);
            color: white;
            border: none;
            padding: 0.5rem 0.75rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.8rem;
            font-weight: 500;
            transition: all 0.2s ease;
            margin-bottom: 0.375rem;
            width: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.375rem;
            min-height: 36px;
        }
        
        .btn:hover {
            background: var(--primary-hover);
            transform: translateY(-1px);
            box-shadow: var(--shadow-md);
        }

        .btn:active {
            transform: translateY(0);
        }
        
        .btn-secondary {
            background: var(--text-secondary);
        }

        .btn-secondary:hover {
            background: var(--text-primary);
        }
        
        .btn-danger {
            background: var(--error-color);
        }

        .btn-danger:hover {
            background: #dc2626;
        }
        
        .btn-success {
            background: var(--success-color);
        }

        .btn-success:hover {
            background: #059669;
        }
        
        .slider-container {
            margin: 1rem 0;
        }

        .slider-label {
            display: flex;
            justify-content: between;
            align-items: center;
            margin-bottom: 0.5rem;
            font-size: 0.875rem;
            font-weight: 500;
        }
        
        .slider {
            width: 100%;
            -webkit-appearance: none;
            height: 6px;
            border-radius: 3px;
            background: var(--bg-tertiary);
            outline: none;
            transition: all 0.2s ease;
        }
        
        .slider::-webkit-slider-thumb {
            -webkit-appearance: none;
            appearance: none;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background: var(--primary-color);
            cursor: pointer;
            box-shadow: var(--shadow-sm);
            transition: all 0.2s ease;
        }

        .slider::-webkit-slider-thumb:hover {
            background: var(--primary-hover);
            transform: scale(1.1);
        }
        
        .pixel-info {
            background: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            border-radius: var(--border-radius);
            padding: 0.75rem;
            margin: 1rem 0;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.875rem;
            min-height: 3rem;
            display: flex;
            align-items: center;
        }
        
        .matrix-container {
            background: var(--bg-primary);
            border: 1px solid var(--border-color);
            border-radius: var(--border-radius);
            padding: 1rem;
            max-height: 400px;
            overflow: auto;
        }

        .matrix-table {
            width: 100%;
            border-collapse: collapse;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
        }

        .matrix-table th,
        .matrix-table td {
            border: 1px solid var(--border-color);
            padding: 0.25rem;
            text-align: center;
            min-width: 60px;
        }

        .matrix-table th {
            background: var(--bg-tertiary);
            font-weight: 600;
            position: sticky;
            top: 0;
            z-index: 10;
        }

        .matrix-table td {
            transition: all 0.2s ease;
        }

        .matrix-table td.highlighted {
            background: var(--primary-color) !important;
            color: white;
            transform: scale(1.05);
            box-shadow: var(--shadow-md);
            z-index: 5;
            position: relative;
        }

        .matrix-table td.neighbor {
            background: rgba(59, 130, 246, 0.1);
        }
        
        .section {
            margin-bottom: 1.5rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid var(--border-color);
        }
        
        .section:last-child {
            border-bottom: none;
        }
        
        .section h4 {
            color: var(--text-primary);
            margin-bottom: 0.75rem;
            font-size: 0.875rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        /* Toast Notifications */
        .toast-container {
            position: fixed;
            top: 1rem;
            right: 1rem;
            z-index: 1000;
            pointer-events: none;
        }

        .toast {
            background: var(--bg-primary);
            border: 1px solid var(--border-color);
            border-radius: var(--border-radius);
            padding: 1rem 1.5rem;
            margin-bottom: 0.5rem;
            box-shadow: var(--shadow-lg);
            display: flex;
            align-items: center;
            gap: 0.75rem;
            min-width: 300px;
            max-width: 400px;
            pointer-events: all;
            transform: translateX(100%);
            transition: all 0.3s ease;
        }

        .toast.show {
            transform: translateX(0);
        }

        .toast.success {
            border-left: 4px solid var(--success-color);
        }

        .toast.error {
            border-left: 4px solid var(--error-color);
        }

        .toast.info {
            border-left: 4px solid var(--primary-color);
        }

        .toast-icon {
            width: 1.25rem;
            height: 1.25rem;
            flex-shrink: 0;
        }

        .toast-content {
            flex: 1;
        }

        .toast-title {
            font-weight: 600;
            font-size: 0.875rem;
            margin-bottom: 0.25rem;
        }

        .toast-message {
            font-size: 0.875rem;
            color: var(--text-secondary);
        }

        .toast-close {
            background: none;
            border: none;
            color: var(--text-secondary);
            cursor: pointer;
            padding: 0.25rem;
            border-radius: 4px;
            transition: all 0.2s ease;
        }

        .toast-close:hover {
            background: var(--bg-tertiary);
            color: var(--text-primary);
        }
        
        .loading-overlay {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(15, 23, 42, 0.5);
            backdrop-filter: blur(4px);
            display: none;
            align-items: center;
            justify-content: center;
            z-index: 999;
        }

        .loading-spinner {
            background: var(--bg-primary);
            border-radius: var(--border-radius);
            padding: 2rem;
            box-shadow: var(--shadow-lg);
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;
        }

        .spinner {
            width: 2rem;
            height: 2rem;
            border: 3px solid var(--border-color);
            border-top: 3px solid var(--primary-color);
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateX(100%);
            }
            to {
                opacity: 1;
                transform: translateX(0);
            }
        }
        
        @keyframes slideOut {
            from {
                opacity: 1;
                transform: translateX(0);
            }
            to {
                opacity: 0;
                transform: translateX(100%);
            }
        }
        
        /* Additional Matrix and Control Panel Styles */
        .matrix-panel {
            grid-column: span 2;
        }
        
        .matrix-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
            padding-bottom: 0.75rem;
            border-bottom: 1px solid var(--border-color);
        }
        
        .matrix-info {
            display: flex;
            gap: 1rem;
            font-size: 0.875rem;
            color: var(--text-secondary);
        }
        
        .cursor-pos, .matrix-size {
            background: var(--bg-tertiary);
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-family: 'JetBrains Mono', monospace;
        }
        
        .matrix-table-container {
            max-height: 400px;
            overflow: auto;
            background: var(--bg-secondary);
            border-radius: var(--border-radius);
            border: 1px solid var(--border-color);
        }
        
        .matrix-table {
            width: 100%;
            border-collapse: collapse;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
        }
        
        .matrix-table th,
        .matrix-table td {
            padding: 0.25rem 0.5rem;
            text-align: center;
            border: 1px solid var(--border-color);
            min-width: 60px;
        }
        
        .matrix-table th {
            background: var(--bg-tertiary);
            font-weight: 600;
            position: sticky;
            top: 0;
            z-index: 10;
        }
        
        .matrix-table td.highlight {
            background: var(--primary-color);
            color: white;
            font-weight: 600;
        }
        
        .matrix-placeholder {
            text-align: center;
            color: var(--text-secondary);
            padding: 2rem;
            font-style: italic;
        }
        
        .operation-section {
            margin-bottom: 1.5rem;
        }
        
        .operation-section h4 {
            color: var(--text-secondary);
            font-size: 0.875rem;
            font-weight: 600;
            margin-bottom: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
        
        .button-group {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 0.375rem;
        }
        
        .button-group.single-column {
            grid-template-columns: 1fr;
        }
        
        .btn-icon {
            margin-right: 0.375rem;
            font-size: 0.9rem;
        }
        
        .upload-section {
            margin-bottom: 1rem;
        }
        
        .file-input {
            width: 100%;
            padding: 0.75rem;
            border: 2px dashed var(--border-color);
            border-radius: var(--border-radius);
            background: var(--bg-secondary);
            color: var(--text-primary);
            cursor: pointer;
            transition: all 0.2s ease;
            font-size: 0.875rem;
        }
        
        .file-input:hover {
            border-color: var(--primary-color);
            background: var(--bg-tertiary);
        }
        
        .canvas-container {
            position: relative;
            display: inline-block;
            border-radius: var(--border-radius);
            overflow: hidden;
            border: 1px solid var(--border-color);
        }
        
        #imageCanvas {
            display: block;
            max-width: 100%;
            height: auto;
        }
    </style>
</head>
<body>
    <div class="app-container">
        <div class="header">
            <div class="header-content">
                <div>
                    <h1>🖼️ Image Processing Learning Tool</h1>
                    <p>Learn image processing concepts with interactive pixel manipulation</p>
                </div>
            </div>
        </div>
        
        <div class="main-content">
            <div class="image-panel">
                <div class="card">
                    <h3>📂 Upload Image</h3>
                    <div class="upload-area" onclick="document.getElementById('imageUpload').click()">
                        <div>
                            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" style="color: var(--text-secondary); margin-bottom: 1rem;">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                            </svg>
                            <p style="font-weight: 500; margin-bottom: 0.5rem;">Drop your image here</p>
                            <p style="font-size: 0.875rem; color: var(--text-secondary);">or click to browse files</p>
                        </div>
                        <input type="file" id="imageUpload" accept="image/*" style="display: none;">
                    </div>
                </div>
                
                <div class="card">
                    <h3>🖼️ Image Canvas</h3>
                    <div class="canvas-container">
                        <canvas id="imageCanvas" width="400" height="300"></canvas>
                    </div>
                    <div class="pixel-info" id="pixelInfo">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" style="color: var(--text-secondary);">
                            <circle cx="12" cy="12" r="10"/>
                            <path d="m9 12 2 2 4-4"/>
                        </svg>
                        Move cursor over image to inspect pixels
                    </div>
                </div>
            </div>
            
            <div class="control-panel">
                <div class="card">
                    <h3>⚙️ Image Operations</h3>
                    
                    <div class="operation-section">
                        <h4>Basic Processing</h4>
                        <div class="button-group">
                            <button class="btn btn-primary" onclick="processImage('grayscale')">
                                <span class="btn-icon">🎨</span>
                                Grayscale
                            </button>
                            <button class="btn btn-primary" onclick="processImage('binary')">
                                <span class="btn-icon">⚫</span>
                                Binary
                            </button>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Brightness Control</h4>
                        <div class="slider-container">
                            <label class="slider-label">
                                <span>Brightness </span>
                                <span class="slider-value" id="brightnessValue">0</span>
                            </label>
                            <input type="range" id="brightnessSlider" min="-100" max="100" value="0" class="slider">
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Zoom Control</h4>
                        <div class="slider-container">
                            <label class="slider-label">
                                <span>Zoom Level </span>
                                <span class="slider-value" id="zoomValue">1.0x</span>
                            </label>
                            <input type="range" id="zoomSlider" min="-100" max="100" value="0" class="slider">
                            <div style="display: flex; justify-content: space-between; font-size: 0.75rem; color: var(--text-secondary); margin-top: 0.25rem;">
                                <span>0.1x</span>
                                <span>1.0x</span>
                                <span>5.0x</span>
                            </div>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Rotate Control</h4>
                        <div class="slider-container">
                            <label class="slider-label">
                                <span>Rotation </span>
                                <span class="slider-value" id="rotateValue">0°</span>
                            </label>
                            <input type="range" id="rotateSlider" min="-180" max="180" value="0" class="slider">
                            <div style="display: flex; justify-content: space-between; font-size: 0.75rem; color: var(--text-secondary); margin-top: 0.25rem;">
                                <span>-360°</span>
                                <span>0°</span>
                                <span>360°</span>
                            </div>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Image Arithmetic</h4>
                        <div class="upload-section">
                            <input type="file" id="secondImageUpload" accept="image/*" class="file-input" placeholder="Select second image">
                        </div>
                        <div class="button-group">
                            <button class="btn btn-primary" onclick="processImage('add')">
                                <span class="btn-icon">➕</span>
                                Add
                            </button>
                            <button class="btn btn-primary" onclick="processImage('subtract')">
                                <span class="btn-icon">➖</span>
                                Subtract
                            </button>
                            <button class="btn btn-primary" onclick="processImage('multiply')">
                                <span class="btn-icon">✖️</span>
                                Multiply
                            </button>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Boolean Logic</h4>
                        <div class="button-group">
                            <button class="btn btn-primary" onclick="processImage('and')">
                                <span class="btn-icon">🔗</span>
                                AND
                            </button>
                            <button class="btn btn-primary" onclick="processImage('or')">
                                <span class="btn-icon">🔀</span>
                                OR
                            </button>
                            <button class="btn btn-primary" onclick="processImage('xor')">
                                <span class="btn-icon">⚡</span>
                                XOR
                            </button>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>Flip Operations</h4>
                        <div class="button-group">
                            <button class="btn btn-primary" onclick="processImage('flipH')">
                                <span class="btn-icon">↔️</span>
                                Flip H
                            </button>
                            <button class="btn btn-primary" onclick="processImage('flipV')">
                                <span class="btn-icon">↕️</span>
                                Flip V
                            </button>
                        </div>
                    </div>
                    

                    
                    <div class="operation-section">
                        <h4>Actions</h4>
                        <div class="button-group">
                            <button class="btn btn-secondary" onclick="resetImage()">
                                <span class="btn-icon">🔄</span>
                                Reset
                            </button>
                            <button class="btn btn-success" onclick="downloadImage()">
                                <span class="btn-icon">📥</span>
                                Download
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="matrix-panel">
                <div class="card">
                    <div class="matrix-header">
                        <h3>🔢 Pixel Matrix</h3>
                        <div class="matrix-info">
                            <span class="cursor-pos" id="cursorPosition">Cursor: (0, 0)</span>
                            <span class="matrix-size" id="matrixSize">Size: 0×0</span>
                            <span class="cursor-pos" style="color: var(--success-color);">⚡ Optimized for Performance</span>
                        </div>
                    </div>
                    <div class="matrix-container">
                        <div class="matrix-table-container" id="matrixDisplay">
                            <p class="matrix-placeholder">Load an image to see the complete pixel matrix</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <!-- Toast Notifications Container -->
    <div class="toast-container" id="toastContainer"></div>
    
    <!-- Loading Overlay -->
    <div class="loading-overlay" id="loadingOverlay">
        <div class="loading-spinner">
            <div class="spinner"></div>
            <p>Processing image...</p>
        </div>
    </div>

    <script>
        // Global state
        let canvas = document.getElementById('imageCanvas');
        let ctx = canvas.getContext('2d');
        let currentImageData = null;
        let imageWidth = 0;
        let imageHeight = 0;
        let pixelMatrix = null;
        let currentCursorX = -1;
        let currentCursorY = -1;
        
        // Performance optimization: Cache and throttling
        let imageCache = new Map();
        let lastProcessTime = 0;
        let isProcessing = false;
        
        // Initialize the application
        function init() {
            setupEventListeners();
            showToast('Welcome to Image Processing Tool!', 'info');
        }
        
        function setupEventListeners() {
            // Upload handlers
            document.getElementById('imageUpload').addEventListener('change', uploadImage);
            document.getElementById('secondImageUpload').addEventListener('change', uploadSecondImage);
            
            // Canvas mouse tracking
            canvas.addEventListener('mousemove', handleMouseMove);
            canvas.addEventListener('mouseleave', () => {
                updateCursorPosition(-1, -1);
                clearMatrixHighlight();
            });
            
            // Drag and drop functionality
            setupDragAndDrop();
            
            // Brightness slider - optimized real-time processing
            document.getElementById('brightnessSlider').addEventListener('input', function() {
                document.getElementById('brightnessValue').textContent = this.value;
                
                if (currentImageData && !isProcessing) {
                    processImageOptimized('brightness');
                }
            });
            
            // Zoom slider - real-time processing
            document.getElementById('zoomSlider').addEventListener('input', function() {
                // Convert slider value to zoom factor for display
                let zoomValue = parseFloat(this.value);
                let zoomFactor;
                if (zoomValue >= 0) {
                    // Positive values: 0 to 100 maps to 1.0 to 5.0
                    zoomFactor = 1.0 + (zoomValue/100.0)*4.0;
                } else {
                    // Negative values: -100 to 0 maps to 0.1 to 1.0
                    zoomFactor = 0.1 + ((zoomValue+100.0)/100.0)*0.9;
                }
                document.getElementById('zoomValue').textContent = zoomFactor.toFixed(1) + 'x';
                
                if (currentImageData && !isProcessing) {
                    processImageWithValueOptimized('zoomSlider', parseInt(this.value));
                }
            });
            
            // Rotate slider - real-time processing
            document.getElementById('rotateSlider').addEventListener('input', function() {
                // Convert slider value to degrees (double the range for smoother control)
                let rotateValue = parseInt(this.value);
                let degrees = rotateValue * 2; // -180 to 180 becomes -360 to 360
                document.getElementById('rotateValue').textContent = degrees + '°';
                
                if (currentImageData && !isProcessing) {
                    processImageWithValueOptimized('rotateSlider', rotateValue);
                }
            });
        }
        
        function setupDragAndDrop() {
            let uploadArea = document.querySelector('.upload-area');
            
            uploadArea.addEventListener('dragover', function(e) {
                e.preventDefault();
                this.style.borderColor = 'var(--primary-color)';
                this.style.background = 'var(--bg-tertiary)';
            });
            
            uploadArea.addEventListener('dragleave', function(e) {
                this.style.borderColor = 'var(--border-color)';
                this.style.background = 'var(--bg-secondary)';
            });
            
            uploadArea.addEventListener('drop', function(e) {
                e.preventDefault();
                this.style.borderColor = 'var(--border-color)';
                this.style.background = 'var(--bg-secondary)';
                
                let files = e.dataTransfer.files;
                if (files.length > 0) {
                    uploadImageFile(files[0]);
                }
            });
        }
        
        function uploadImage() {
            let file = document.getElementById('imageUpload').files[0];
            uploadImageFile(file);
        }
        
        function uploadImageFile(file) {
            if (!file) return;
            
            console.log('Starting image upload:', file.name);
            showLoading(true);
            let formData = new FormData();
            formData.append('image', file);
            
            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                console.log('Upload response status:', response.status);
                return response.json();
            })
            .then(data => {
                console.log('Upload response data:', data);
                showLoading(false);
                if (data.success) {
                    console.log('Displaying image and matrix...');
                    displayImage(data.imageData);
                    displayMatrix(data.pixelMatrix);
                    showToast('Image uploaded successfully!', 'success');
                } else {
                    showToast('Failed to upload image: ' + data.message, 'error');
                }
            })
            .catch(error => {
                console.error('Upload error:', error);
                showLoading(false);
                showToast('Upload failed: ' + error.message, 'error');
            });
        }
        
        function uploadSecondImage() {
            let file = document.getElementById('secondImageUpload').files[0];
            if (!file) return;
            
            showLoading(true);
            let formData = new FormData();
            formData.append('image', file);
            
            fetch('/upload-second', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    showToast('Second image uploaded successfully!', 'success');
                } else {
                    showToast('Failed to upload second image: ' + data.message, 'error');
                }
            })
            .catch(error => {
                showLoading(false);
                showToast('Second image upload failed: ' + error.message, 'error');
            });
        }
        
        function handleMouseMove(e) {
            if (!currentImageData) return;
            
            let rect = canvas.getBoundingClientRect();
            let scaleX = canvas.width / rect.width;
            let scaleY = canvas.height / rect.height;
            let x = Math.floor((e.clientX - rect.left) * scaleX);
            let y = Math.floor((e.clientY - rect.top) * scaleY);
            
            if (x >= 0 && x < canvas.width && y >= 0 && y < canvas.height) {
                currentCursorX = x;
                currentCursorY = y;
                updateCursorPosition(x, y);
                highlightMatrixCell(x, y);
                getPixelInfo(x, y);
            }
        }
        
        function updateCursorPosition(x, y) {
            let posElement = document.getElementById('cursorPosition');
            if (x >= 0 && y >= 0) {
                posElement.textContent = 'Cursor: (' + x + ', ' + y + ')';
            } else {
                posElement.textContent = 'Cursor: (--, --)';
            }
        }
        
        function getPixelInfo(x, y) {
            fetch('/pixel-info?x=' + x + '&y=' + y)
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('pixelInfo').innerHTML = 
                        '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" style="color: var(--success-color);">' +
                        '<path d="m9 12 2 2 4-4"/>' +
                        '<circle cx="12" cy="12" r="10"/>' +
                        '</svg>' +
                        'Position: (' + x + ', ' + y + ') | RGB: (' + data.r + ', ' + data.g + ', ' + data.b + ') | Gray: ' + data.grayscale;
                }
            })
            .catch(error => {
                document.getElementById('pixelInfo').innerHTML = 
                    '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" style="color: var(--error-color);">' +
                    '<circle cx="12" cy="12" r="10"/>' +
                    '<line x1="15" y1="9" x2="9" y2="15"/>' +
                    '<line x1="9" y1="9" x2="15" y2="15"/>' +
                    '</svg>' +
                    'Error reading pixel data';
            });
        }
        
        function displayImage(imageData) {
            console.log('Displaying image, data length:', imageData ? imageData.length : 'null');
            currentImageData = imageData;
            let img = new Image();
            img.onload = function() {
                console.log('Image loaded successfully, size:', img.width + 'x' + img.height);
                canvas.width = img.width;
                canvas.height = img.height;
                imageWidth = img.width;
                imageHeight = img.height;
                ctx.drawImage(img, 0, 0);
                
                // Update matrix size display
                document.getElementById('matrixSize').textContent = 'Size: ' + imageWidth + '×' + imageHeight;
            };
            img.onerror = function() {
                console.error('Failed to load image');
                showToast('Failed to load image', 'error');
            };
            img.src = 'data:image/png;base64,' + imageData;
        }
        
        function displayMatrix(matrixData) {
            // Skip matrix updates for real-time operations to improve performance
            if (isProcessing) return;
            
            let matrixContainer = document.getElementById('matrixDisplay');
            
            if (!matrixData) {
                matrixContainer.innerHTML = '<p class="matrix-placeholder">No pixel data available</p>';
                return;
            }
            
            // Parse the JSON string to get the matrix array
            let matrix;
            try {
                matrix = typeof matrixData === 'string' ? JSON.parse(matrixData) : matrixData;
            } catch (e) {
                console.error('Failed to parse matrix data:', e);
                matrixContainer.innerHTML = '<p class="matrix-placeholder">Error parsing matrix data</p>';
                return;
            }
            
            pixelMatrix = matrix;
            
            if (!matrix || matrix.length === 0) {
                matrixContainer.innerHTML = '<p class="matrix-placeholder">No pixel data available</p>';
                return;
            }
            
            // Only update matrix if significantly different from current display
            let currentSize = matrixContainer.querySelector('.matrix-table');
            if (currentSize && matrix.length > 100) {
                // Skip frequent matrix updates for large matrices during real-time operations
                return;
            }
            
            // Fast matrix rendering with optimized DOM manipulation
            let html = '<table class="matrix-table">';
            
            // Add header row
            html += '<tr><th></th>';
            for (let j = 0; j < matrix[0].length; j++) {
                html += '<th>' + j + '</th>';
            }
            html += '</tr>';
            
            // Add data rows with batch processing for better performance
            for (let i = 0; i < matrix.length; i++) {
                html += '<tr><th>' + i + '</th>';
                for (let j = 0; j < matrix[i].length; j++) {
                    let value = matrix[i][j] || 0;
                    html += '<td id="cell-' + i + '-' + j + '">' + value + '</td>';
                }
                html += '</tr>';
            }
            html += '</table>';
            
            // Single DOM update for better performance
            matrixContainer.innerHTML = html;
            
            // Show matrix info
            console.log('Matrix displayed: ' + matrix.length + 'x' + matrix[0].length + ' (optimized for performance)');
        }
        
        function highlightMatrixCell(x, y) {
            clearMatrixHighlight();
            let cell = document.getElementById('cell-' + y + '-' + x);
            if (cell) {
                cell.classList.add('highlight');
            }
        }
        
        function clearMatrixHighlight() {
            let highlightedCells = document.querySelectorAll('.matrix-table td.highlight');
            highlightedCells.forEach(cell => cell.classList.remove('highlight'));
        }
        
        function processImage(operation) {
            if (!currentImageData && operation !== 'reset') {
                showToast('Please upload an image first', 'error');
                return;
            }
            
            showLoading(true);
            
            let requestData = {
                operation: operation
            };
            
            // Add brightness value for brightness operations
            if (operation === 'brightness') {
                requestData.value = parseInt(document.getElementById('brightnessSlider').value);
            }
            
            fetch('/process', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    if (data.imageData) {
                        displayImage(data.imageData);
                    }
                    if (data.pixelMatrix) {
                        displayMatrix(data.pixelMatrix);
                    }
                    showToast(operation.charAt(0).toUpperCase() + operation.slice(1) + ' operation completed!', 'success');
                } else {
                    showToast('Processing failed: ' + data.message, 'error');
                }
            })
            .catch(error => {
                showLoading(false);
                showToast('Processing failed: ' + error.message, 'error');
            });
        }
        
        function processImageWithValue(operation, value) {
            if (!currentImageData) {
                showToast('Please upload an image first', 'error');
                return;
            }
            
            showLoading(true);
            
            let requestData = {
                operation: operation,
                value: value
            };
            
            fetch('/process', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    if (data.imageData) {
                        displayImage(data.imageData);
                    }
                    if (data.pixelMatrix) {
                        displayMatrix(data.pixelMatrix);
                    }
                    // Don't show toast for slider operations to avoid spam
                } else {
                    showToast('Processing failed: ' + data.message, 'error');
                }
            })
            .catch(error => {
                showLoading(false);
                showToast('Processing failed: ' + error.message, 'error');
            });
        }
        
        // Optimized processing functions with caching and throttling
        function processImageOptimized(operation) {
            if (!currentImageData || isProcessing) return;
            
            isProcessing = true;
            let now = Date.now();
            
            // Throttle requests to max 30fps (33ms intervals)
            if (now - lastProcessTime < 33) {
                setTimeout(() => {
                    isProcessing = false;
                    processImageOptimized(operation);
                }, 33 - (now - lastProcessTime));
                return;
            }
            
            lastProcessTime = now;
            let value = parseInt(document.getElementById('brightnessSlider').value);
            let cacheKey = operation + '_' + value;
            
            // Check cache first
            if (imageCache.has(cacheKey)) {
                let cachedData = imageCache.get(cacheKey);
                displayImage(cachedData.imageData);
                if (cachedData.pixelMatrix) {
                    displayMatrix(cachedData.pixelMatrix);
                }
                isProcessing = false;
                return;
            }
            
            // Process normally if not in cache
            processImageWithCaching(operation, value, cacheKey);
        }
        
        function processImageWithValueOptimized(operation, value) {
            if (!currentImageData || isProcessing) return;
            
            isProcessing = true;
            let now = Date.now();
            
            // Throttle requests to max 30fps
            if (now - lastProcessTime < 33) {
                setTimeout(() => {
                    isProcessing = false;
                    processImageWithValueOptimized(operation, value);
                }, 33 - (now - lastProcessTime));
                return;
            }
            
            lastProcessTime = now;
            let cacheKey = operation + '_' + value;
            
            // Check cache first
            if (imageCache.has(cacheKey)) {
                let cachedData = imageCache.get(cacheKey);
                displayImage(cachedData.imageData);
                if (cachedData.pixelMatrix) {
                    displayMatrix(cachedData.pixelMatrix);
                }
                isProcessing = false;
                return;
            }
            
            // Process normally if not in cache
            processImageWithCaching(operation, value, cacheKey);
        }
        
        function processImageWithCaching(operation, value, cacheKey) {
            let requestData = { operation: operation };
            if (value !== undefined) {
                requestData.value = value;
            }
            
            fetch('/process', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(requestData)
            })
            .then(response => response.json())
            .then(data => {
                isProcessing = false;
                if (data.success && data.imageData) {
                    // Cache the result
                    imageCache.set(cacheKey, {
                        imageData: data.imageData,
                        pixelMatrix: data.pixelMatrix
                    });
                    
                    // Limit cache size to prevent memory issues
                    if (imageCache.size > 50) {
                        let firstKey = imageCache.keys().next().value;
                        imageCache.delete(firstKey);
                    }
                    
                    displayImage(data.imageData);
                    if (data.pixelMatrix) {
                        displayMatrix(data.pixelMatrix);
                    }
                } else {
                    showToast('Processing failed: ' + data.message, 'error');
                }
            })
            .catch(error => {
                isProcessing = false;
                showToast('Processing failed: ' + error.message, 'error');
            });
        }
        
        function resetImage() {
            // Clear the image cache to ensure fresh start
            imageCache.clear();
            
            // Reset sliders to default values
            document.getElementById('brightnessSlider').value = 0;
            document.getElementById('brightnessValue').textContent = '0';
            document.getElementById('zoomSlider').value = 0;
            document.getElementById('zoomValue').textContent = '1.0x';
            document.getElementById('rotateSlider').value = 0;
            document.getElementById('rotateValue').textContent = '0°';
            
            // Reset processing state
            isProcessing = false;
            
            processImage('reset');
            showToast('Image reset to original state', 'success');
        }
        
        function downloadImage() {
            if (!currentImageData) {
                showToast('No processed image to download', 'error');
                return;
            }
            
            fetch('/download')
            .then(response => response.blob())
            .then(blob => {
                let url = window.URL.createObjectURL(blob);
                let a = document.createElement('a');
                a.href = url;
                a.download = 'processed_image.png';
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                window.URL.revokeObjectURL(url);
                showToast('Image downloaded successfully!', 'success');
            })
            .catch(error => {
                showToast('Download failed: ' + error.message, 'error');
            });
        }
        
        function showLoading(show) {
            document.getElementById('loadingOverlay').style.display = show ? 'flex' : 'none';
        }
        
        function showToast(message, type = 'info') {
            const toastContainer = document.getElementById('toastContainer');
            const toastId = 'toast-' + Date.now();
            
            const toast = document.createElement('div');
            toast.className = 'toast toast-' + type;
            toast.id = toastId;
            
            const icon = getToastIcon(type);
            
            toast.innerHTML = 
                '<div class="toast-content">' +
                    '<div class="toast-icon">' + icon + '</div>' +
                    '<div>' +
                        '<div class="toast-title">' + getToastTitle(type) + '</div>' +
                        '<div class="toast-message">' + message + '</div>' +
                    '</div>' +
                '</div>' +
                '<button class="toast-close" onclick="removeToast(\'' + toastId + '\')">×</button>';
            
            toastContainer.appendChild(toast);
            
            // Auto-remove after delay
            setTimeout(() => removeToast(toastId), type === 'error' ? 6000 : 4000);
        }
        
        function getToastIcon(type) {
            const icons = {
                success: '✅',
                error: '❌', 
                warning: '⚠️',
                info: 'ℹ️'
            };
            return icons[type] || icons.info;
        }
        
        function getToastTitle(type) {
            const titles = {
                success: 'Success',
                error: 'Error',
                warning: 'Warning', 
                info: 'Information'
            };
            return titles[type] || titles.info;
        }
        
        function removeToast(toastId) {
            const toast = document.getElementById(toastId);
            if (toast) {
                toast.style.animation = 'slideOut 0.3s ease-in-out';
                setTimeout(() => {
                    if (toast.parentNode) {
                        toast.parentNode.removeChild(toast);
                    }
                }, 300);
            }
        }
        
        // Initialize on page load
        document.addEventListener('DOMContentLoaded', init);
    </script>
</body>
</html>
`
	
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, nil)
}

// handleUpload handles the primary image upload
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	file, _, err := r.FormFile("image")
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to read uploaded file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to decode image: " + err.Error(),
		})
		return
	}
	
	processor.originalImage = img   // Store the true original
	processor.currentImage = img    // Working copy
	processor.processedImage = img
	
	// Convert image to base64 for display
	imageData, err := imageToBase64(img)
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to encode image: " + err.Error(),
		})
		return
	}
	
	// Generate pixel matrix
	pixelMatrix := generatePixelMatrix(img)
	
	sendJSONResponse(w, ProcessResponse{
		Success:     true,
		Message:     "Image uploaded successfully",
		ImageData:   imageData,
		PixelMatrix: pixelMatrix,
	})
}

// handleUploadSecond handles the second image upload
func handleUploadSecond(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	file, _, err := r.FormFile("image")
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to read uploaded file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to decode image: " + err.Error(),
		})
		return
	}
	
	processor.secondImage = img
	
	sendJSONResponse(w, ProcessResponse{
		Success: true,
		Message: "Second image uploaded successfully",
	})
}

// handlePixelInfo handles pixel information requests
func handlePixelInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse query parameters
	x := 0
	y := 0
	var err error
	
	if xStr := r.URL.Query().Get("x"); xStr != "" {
		x, err = strconv.Atoi(xStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid x coordinate",
			})
			return
		}
	}
	
	if yStr := r.URL.Query().Get("y"); yStr != "" {
		y, err = strconv.Atoi(yStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid y coordinate",
			})
			return
		}
	}
	
	// Use current image if processed image is nil
	img := processor.processedImage
	if img == nil {
		img = processor.currentImage
	}
	
	if img == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No image loaded",
		})
		return
	}
	
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X ||
		y < bounds.Min.Y || y >= bounds.Max.Y {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Coordinates out of bounds",
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
	
	response := map[string]interface{}{
		"success":   true,
		"x":         x,
		"y":         y,
		"r":         rVal,
		"g":         gVal,
		"b":         bVal,
		"a":         uint8(alpha >> 8),
		"grayscale": grayscale,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleProcess handles image processing operations
func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		Operation string `json:"operation"`
		Value     int    `json:"value"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}
	
	if processor.currentImage == nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "No image loaded",
		})
		return
	}
	
	var resultImage image.Image
	var err2 error
	var skipMatrix bool = false // Skip matrix generation for performance
	
	switch request.Operation {
	case "grayscale":
		grayImg := convertToGrayscale(processor.currentImage)
		processor.currentImage = grayImg // Update base image
		processor.brightnessValue = 0
		processor.zoomValue = 0
		processor.rotateValue = 0
		resultImage = grayImg
	case "binary":
		binaryImg := convertToBinary(processor.currentImage)
		processor.currentImage = binaryImg // Update base image
		processor.brightnessValue = 0
		processor.zoomValue = 0
		processor.rotateValue = 0
		resultImage = binaryImg
	case "brightness":
		processor.brightnessValue = request.Value
		resultImage = applyAllTransformations(processor.currentImage, processor.brightnessValue, processor.zoomValue, processor.rotateValue)
		// Skip matrix generation for real-time operations
		skipMatrix = true
	case "add", "subtract", "multiply":
		if processor.secondImage == nil {
			sendJSONResponse(w, ProcessResponse{
				Success: false,
				Message: "Second image required for arithmetic operations",
			})
			return
		}
		resultImage, err2 = arithmeticOperation(processor.currentImage, processor.secondImage, request.Operation)
	case "and", "or", "xor":
		if processor.secondImage == nil {
			sendJSONResponse(w, ProcessResponse{
				Success: false,
				Message: "Second image required for boolean operations",
			})
			return
		}
		resultImage, err2 = booleanOperation(processor.currentImage, processor.secondImage, request.Operation)

	case "flipH":
		// Apply flip as a transformation, don't permanently modify currentImage
		baseImg := applyAllTransformations(processor.currentImage, processor.brightnessValue, processor.zoomValue, processor.rotateValue)
		resultImage = flipHorizontal(baseImg)
	case "flipV":
		// Apply flip as a transformation, don't permanently modify currentImage  
		baseImg := applyAllTransformations(processor.currentImage, processor.brightnessValue, processor.zoomValue, processor.rotateValue)
		resultImage = flipVertical(baseImg)
	case "zoomSlider":
		processor.zoomValue = request.Value
		resultImage = applyAllTransformations(processor.currentImage, processor.brightnessValue, processor.zoomValue, processor.rotateValue)
		// Skip matrix generation for real-time operations
		skipMatrix = true
	case "rotateSlider":
		processor.rotateValue = request.Value
		resultImage = applyAllTransformations(processor.currentImage, processor.brightnessValue, processor.zoomValue, processor.rotateValue)
		// Skip matrix generation for real-time operations  
		skipMatrix = true
	case "translate":
		sourceImg := processor.processedImage
		if sourceImg == nil {
			sourceImg = processor.currentImage
		}
		resultImage = translateImage(sourceImg, 50, 50)
	case "reset":
		// Check if we have an original image
		if processor.originalImage == nil {
			sendJSONResponse(w, ProcessResponse{
				Success: false,
				Message: "No original image available to reset to",
			})
			return
		}
		// Reset all transformation values
		processor.brightnessValue = 0
		processor.zoomValue = 0
		processor.rotateValue = 0
		// Restore the true original image
		processor.currentImage = processor.originalImage
		resultImage = processor.originalImage
	default:
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Unknown operation: " + request.Operation,
		})
		return
	}
	
	if err2 != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Processing failed: " + err2.Error(),
		})
		return
	}
	
	processor.processedImage = resultImage
	
	// Convert image to base64
	imageData, err := imageToBase64(resultImage)
	if err != nil {
		sendJSONResponse(w, ProcessResponse{
			Success: false,
			Message: "Failed to encode processed image: " + err.Error(),
		})
		return
	}
	
	// Generate pixel matrix only if needed (skip for real-time operations)
	var pixelMatrix string
	if !skipMatrix {
		pixelMatrix = generatePixelMatrix(resultImage)
	}
	
	sendJSONResponse(w, ProcessResponse{
		Success:     true,
		Message:     "Image processed successfully",
		ImageData:   imageData,
		PixelMatrix: pixelMatrix,
	})
}

// handleDownload serves the processed image for download
func handleDownload(w http.ResponseWriter, r *http.Request) {
	if processor.processedImage == nil {
		http.Error(w, "No processed image available", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=processed_image.png")
	
	err := png.Encode(w, processor.processedImage)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
	}
}

// Helper functions for image processing

func convertToGrayscale(img image.Image) image.Image {
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

func convertToBinary(img image.Image) image.Image {
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

func adjustBrightness(img image.Image, brightness int) image.Image {
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
				if val < 0 { val = 0 }
				if val > 255 { val = 255 }
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

func resizeImage(src image.Image, newWidth, newHeight int) image.Image {
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
				if srcY >= srcHeight { srcY = srcHeight - 1 }
				
				for x := 0; x < newWidth; x++ {
					srcX := int(float64(x) * xRatio)
					if srcX >= srcWidth { srcX = srcWidth - 1 }
					
					dst.Set(x, y, src.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
				}
			}
		}(startY, endY)
	}
	
	wg.Wait()
	return dst
}

func arithmeticOperation(img1, img2 image.Image, operation string) (image.Image, error) {
	bounds := img1.Bounds()
	resizedImg2 := resizeImage(img2, bounds.Dx(), bounds.Dy())
	
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
			if rResult < 0 { rResult = 0 }
			if rResult > 255 { rResult = 255 }
			if gResult < 0 { gResult = 0 }
			if gResult > 255 { gResult = 255 }
			if bResult < 0 { bResult = 0 }
			if bResult > 255 { bResult = 255 }
			
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

func booleanOperation(img1, img2 image.Image, operation string) (image.Image, error) {
	bounds := img1.Bounds()
	resizedImg2 := resizeImage(img2, bounds.Dx(), bounds.Dy())
	
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

func rotate90(src image.Image) image.Image {
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

func flipHorizontal(src image.Image) image.Image {
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

func flipVertical(src image.Image) image.Image {
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

func zoomByFactor(src image.Image, zoomFactor float64) image.Image {
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
	
	return resizeImage(src, newWidth, newHeight)
}

func rotateByAngle(src image.Image, angle float64) image.Image {
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
		return rotate90(src)
	} else if angle == 180 {
		return rotate180(src)
	} else if angle == 270 {
		return rotate270(src)
	}
	
	// For other angles, use simple step rotation (90-degree increments for now)
	steps := int(angle / 90)
	result := src
	for i := 0; i < steps; i++ {
		result = rotate90(result)
	}
	return result
}

func rotate180(src image.Image) image.Image {
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

func rotate270(src image.Image) image.Image {
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

func translateImage(src image.Image, offsetX, offsetY int) image.Image {
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

// Utility functions

func imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func generatePixelMatrix(img image.Image) string {
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
		sampledImg = createSampledMatrix(img, matrixWidth, matrixHeight)
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

// Fast sampling function that creates a representative matrix without full resize
func createSampledMatrix(img image.Image, targetWidth, targetHeight int) image.Image {
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

func applyAllTransformations(baseImg image.Image, brightnessVal, zoomVal, rotateVal int) image.Image {
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
		result = zoomByFactor(result, zoomFactor)
	}
	
	// 2. Apply rotation if not default (0)
	if rotateVal != 0 {
		angle := float64(rotateVal * 2)
		result = rotateByAngle(result, angle)
	}
	
	// 3. Apply brightness if not default (0)
	if brightnessVal != 0 {
		result = adjustBrightness(result, brightnessVal)
	}
	
	return result
}

func sendJSONResponse(w http.ResponseWriter, response ProcessResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}