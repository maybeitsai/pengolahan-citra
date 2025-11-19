package templates

import (
	"html/template"
	"net/http"
)

// TemplateHandler handles template rendering
type TemplateHandler struct {
	template *template.Template
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

// ServeHome serves the main HTML page with modern UI/UX
func (th *TemplateHandler) ServeHome(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `<!DOCTYPE html>
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

        .btn-secondary {
            background: var(--text-secondary);
        }

        .btn-secondary:hover {
            background: var(--text-primary);
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
            justify-content: space-between;
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

        .matrix-table td.highlight {
            background: var(--primary-color);
            color: white;
            font-weight: 600;
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
        
        .matrix-placeholder {
            text-align: center;
            color: var(--text-secondary);
            padding: 2rem;
            font-style: italic;
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
                        <h4>Histogram Analysis</h4>
                        <div class="button-group single-column">
                            <button class="btn btn-success" onclick="showHistogramAnalysis()">
                                <span class="btn-icon">📊</span>
                                Show Histogram & Analysis
                            </button>
                        </div>
                    </div>
                    
                    <div class="operation-section">
                        <h4>🔍 Edge Detection</h4>
                        <div class="button-group">
                            <button class="btn btn-primary" onclick="performEdgeDetection('sobel')">
                                <span class="btn-icon">🔍</span>
                                Full Sobel
                            </button>
                            <button class="btn btn-primary" onclick="performEdgeDetection('sobelX')">
                                <span class="btn-icon">↕️</span>
                                Sobel X
                            </button>
                            <button class="btn btn-primary" onclick="performEdgeDetection('sobelY')">
                                <span class="btn-icon">↔️</span>
                                Sobel Y
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
        
        // Initialize the application
        function init() {
            setupEventListeners();
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
            
            // Brightness slider
            document.getElementById('brightnessSlider').addEventListener('input', function() {
                document.getElementById('brightnessValue').textContent = this.value;
                if (currentImageData) {
                    processImageWithValue('brightness', parseInt(this.value));
                }
            });
            
            // Zoom slider
            document.getElementById('zoomSlider').addEventListener('input', function() {
                let zoomValue = parseFloat(this.value);
                let zoomFactor;
                if (zoomValue >= 0) {
                    zoomFactor = 1.0 + (zoomValue/100.0)*4.0;
                } else {
                    zoomFactor = 0.1 + ((zoomValue+100.0)/100.0)*0.9;
                }
                document.getElementById('zoomValue').textContent = zoomFactor.toFixed(1) + 'x';
                
                if (currentImageData) {
                    processImageWithValue('zoomSlider', parseInt(this.value));
                }
            });
            
            // Rotate slider
            document.getElementById('rotateSlider').addEventListener('input', function() {
                let rotateValue = parseInt(this.value);
                let degrees = rotateValue * 2;
                document.getElementById('rotateValue').textContent = degrees + '°';
                
                if (currentImageData) {
                    processImageWithValue('rotateSlider', rotateValue);
                }
            });
        }
        
        function uploadImage() {
            let file = document.getElementById('imageUpload').files[0];
            uploadImageFile(file);
        }
        
        function uploadImageFile(file) {
            if (!file) return;
            
            showLoading(true);
            let formData = new FormData();
            formData.append('image', file);
            
            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    displayImage(data.imageData);
                    displayMatrix(data.pixelMatrix);
                } else {
                    alert('Failed to upload image: ' + data.message);
                }
            })
            .catch(error => {
                showLoading(false);
                alert('Upload failed: ' + error.message);
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
                    alert('Second image uploaded successfully!');
                } else {
                    alert('Failed to upload second image: ' + data.message);
                }
            })
            .catch(error => {
                showLoading(false);
                alert('Second image upload failed: ' + error.message);
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
                    document.getElementById('pixelInfo').textContent = 
                        'Position: (' + x + ', ' + y + ') | RGB: (' + data.r + ', ' + data.g + ', ' + data.b + ') | Gray: ' + data.grayscale;
                }
            })
            .catch(error => {
                document.getElementById('pixelInfo').textContent = 'Error reading pixel data';
            });
        }
        
        function displayImage(imageData) {
            currentImageData = imageData;
            let img = new Image();
            img.onload = function() {
                canvas.width = img.width;
                canvas.height = img.height;
                imageWidth = img.width;
                imageHeight = img.height;
                ctx.drawImage(img, 0, 0);
                
                document.getElementById('matrixSize').textContent = 'Size: ' + imageWidth + '×' + imageHeight;
            };
            img.src = 'data:image/png;base64,' + imageData;
        }
        
        function displayMatrix(matrixData) {
            let matrixContainer = document.getElementById('matrixDisplay');
            
            if (!matrixData) {
                matrixContainer.innerHTML = '<p class="matrix-placeholder">No pixel data available</p>';
                return;
            }
            
            let matrix;
            try {
                matrix = typeof matrixData === 'string' ? JSON.parse(matrixData) : matrixData;
            } catch (e) {
                matrixContainer.innerHTML = '<p class="matrix-placeholder">Error parsing matrix data</p>';
                return;
            }
            
            pixelMatrix = matrix;
            
            if (!matrix || matrix.length === 0) {
                matrixContainer.innerHTML = '<p class="matrix-placeholder">No pixel data available</p>';
                return;
            }
            
            let html = '<table class="matrix-table">';
            
            // Add header row
            html += '<tr><th></th>';
            for (let j = 0; j < matrix[0].length; j++) {
                html += '<th>' + j + '</th>';
            }
            html += '</tr>';
            
            // Add data rows
            for (let i = 0; i < matrix.length; i++) {
                html += '<tr><th>' + i + '</th>';
                for (let j = 0; j < matrix[i].length; j++) {
                    let value = matrix[i][j] || 0;
                    html += '<td id="cell-' + i + '-' + j + '">' + value + '</td>';
                }
                html += '</tr>';
            }
            html += '</table>';
            
            matrixContainer.innerHTML = html;
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
                alert('Please upload an image first');
                return;
            }
            
            showLoading(true);
            
            let requestData = {
                operation: operation
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
                } else {
                    alert('Processing failed: ' + data.message);
                }
            })
            .catch(error => {
                showLoading(false);
                alert('Processing failed: ' + error.message);
            });
        }
        
        function processImageWithValue(operation, value) {
            if (!currentImageData) {
                alert('Please upload an image first');
                return;
            }
            
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
                if (data.success) {
                    if (data.imageData) {
                        displayImage(data.imageData);
                    }
                    if (data.pixelMatrix) {
                        displayMatrix(data.pixelMatrix);
                    }
                } else {
                    alert('Processing failed: ' + data.message);
                }
            })
            .catch(error => {
                alert('Processing failed: ' + error.message);
            });
        }
        
        function resetImage() {
            document.getElementById('brightnessSlider').value = 0;
            document.getElementById('brightnessValue').textContent = '0';
            document.getElementById('zoomSlider').value = 0;
            document.getElementById('zoomValue').textContent = '1.0x';
            document.getElementById('rotateSlider').value = 0;
            document.getElementById('rotateValue').textContent = '0°';
            
            processImage('reset');
        }
        
        function downloadImage() {
            if (!currentImageData) {
                alert('No processed image to download');
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
            })
            .catch(error => {
                alert('Download failed: ' + error.message);
            });
        }
        
        function showHistogramAnalysis() {
            if (!currentImageData) {
                alert('Please upload an image first');
                return;
            }
            
            showLoading(true);
            
            fetch('/histogram', {
                method: 'GET'
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    displayHistogramAnalysis(data.data);
                } else {
                    alert('Failed to generate histogram: ' + data.message);
                }
            })
            .catch(error => {
                showLoading(false);
                alert('Failed to generate histogram: ' + error.message);
            });
        }
        
        function displayHistogramAnalysis(histogramData) {
            // Create modal window for histogram display
            const modal = document.createElement('div');
            modal.style.cssText = ` + "`" + `
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: rgba(0, 0, 0, 0.8);
                z-index: 1000;
                display: flex;
                align-items: center;
                justify-content: center;
                backdrop-filter: blur(4px);
            ` + "`" + `;
            
            const modalContent = document.createElement('div');
            modalContent.style.cssText = ` + "`" + `
                background: white;
                border-radius: 12px;
                padding: 2rem;
                max-width: 90vw;
                max-height: 90vh;
                overflow-y: auto;
                box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
                position: relative;
            ` + "`" + `;
            
            // Create histogram visualization
            const histogramHTML = ` + "`" + `
                <div style="font-family: 'Inter', sans-serif;">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; border-bottom: 2px solid #e2e8f0; padding-bottom: 1rem;">
                        <h2 style="margin: 0; color: #1e293b; font-size: 1.5rem;">📊 Histogram Analysis</h2>
                        <button onclick="this.closest('.modal').remove()" style="background: #ef4444; color: white; border: none; border-radius: 50%; width: 2rem; height: 2rem; cursor: pointer; font-size: 1rem;">×</button>
                    </div>
                    
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-bottom: 2rem;">
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">📈 RGB Histograms</h3>
                            <canvas id="rgbHistogramCanvas" width="400" height="200" style="border: 1px solid #e2e8f0; border-radius: 8px;"></canvas>
                        </div>
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">⚫ Grayscale Histogram</h3>
                            <canvas id="grayHistogramCanvas" width="400" height="200" style="border: 1px solid #e2e8f0; border-radius: 8px;"></canvas>
                        </div>
                    </div>
                    
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-bottom: 2rem;">
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">🎯 Binary Image (Otsu's Method)</h3>
                            <div style="text-align: center;">
                                <p style="margin: 0.5rem 0; font-weight: 500;">Threshold: ` + "${histogramData.threshold}" + `</p>
                                <img src="data:image/png;base64,` + "${histogramData.binaryImageData}" + `" style="max-width: 100%; border: 1px solid #e2e8f0; border-radius: 8px;" alt="Binary Image">
                            </div>
                        </div>
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">⚖️ Histogram Equalized</h3>
                            <div style="text-align: center;">
                                <img src="data:image/png;base64,` + "${histogramData.equalizedImageData}" + `" style="max-width: 100%; border: 1px solid #e2e8f0; border-radius: 8px;" alt="Equalized Image">
                            </div>
                        </div>
                    </div>
                    
                    <div style="background: #f8fafc; border-radius: 8px; padding: 1.5rem; margin-bottom: 1rem;">
                        <h3 style="color: #374151; margin-bottom: 1rem;">📊 Statistical Analysis</h3>
                        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem;">
                            <div>
                                <h4 style="margin: 0 0 0.5rem 0; color: #6b7280;">Original Image</h4>
                                <p style="margin: 0.25rem 0; font-family: 'JetBrains Mono', monospace;">Mean: ` + "${histogramData.meanBefore.toFixed(2)}" + `</p>
                                <p style="margin: 0.25rem 0; font-family: 'JetBrains Mono', monospace;">Std Dev: ` + "${histogramData.stdDevBefore.toFixed(2)}" + `</p>
                            </div>
                            <div>
                                <h4 style="margin: 0 0 0.5rem 0; color: #6b7280;">After Equalization</h4>
                                <p style="margin: 0.25rem 0; font-family: 'JetBrains Mono', monospace;">Mean: ` + "${histogramData.meanAfter.toFixed(2)}" + `</p>
                                <p style="margin: 0.25rem 0; font-family: 'JetBrains Mono', monospace;">Std Dev: ` + "${histogramData.stdDevAfter.toFixed(2)}" + `</p>
                            </div>
                        </div>
                    </div>
                    
                    <div style="text-align: center;">
                        <button onclick="this.closest('.modal').remove()" style="background: #3b82f6; color: white; border: none; padding: 0.75rem 2rem; border-radius: 6px; cursor: pointer; font-size: 1rem; font-weight: 500;">Close Analysis</button>
                    </div>
                </div>
            ` + "`" + `;
            
            modalContent.innerHTML = histogramHTML;
            modal.appendChild(modalContent);
            modal.className = 'modal';
            document.body.appendChild(modal);
            
            // Draw histograms after modal is added to DOM
            setTimeout(() => {
                drawRGBHistogram(histogramData);
                drawGrayscaleHistogram(histogramData);
            }, 100);
        }
        
        function drawRGBHistogram(data) {
            const canvas = document.getElementById('rgbHistogramCanvas');
            if (!canvas) return;
            
            const ctx = canvas.getContext('2d');
            const width = canvas.width;
            const height = canvas.height;
            
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Find maximum value for scaling
            const maxRed = Math.max(...data.redHistogram);
            const maxGreen = Math.max(...data.greenHistogram);
            const maxBlue = Math.max(...data.blueHistogram);
            const maxValue = Math.max(maxRed, maxGreen, maxBlue);
            
            if (maxValue === 0) return;
            
            // Draw histograms
            const barWidth = width / 256;
            
            // Red histogram
            ctx.fillStyle = 'rgba(255, 0, 0, 0.6)';
            for (let i = 0; i < 256; i++) {
                const barHeight = (data.redHistogram[i] / maxValue) * height;
                ctx.fillRect(i * barWidth, height - barHeight, barWidth, barHeight);
            }
            
            // Green histogram
            ctx.fillStyle = 'rgba(0, 255, 0, 0.6)';
            for (let i = 0; i < 256; i++) {
                const barHeight = (data.greenHistogram[i] / maxValue) * height;
                ctx.fillRect(i * barWidth, height - barHeight, barWidth, barHeight);
            }
            
            // Blue histogram
            ctx.fillStyle = 'rgba(0, 0, 255, 0.6)';
            for (let i = 0; i < 256; i++) {
                const barHeight = (data.blueHistogram[i] / maxValue) * height;
                ctx.fillRect(i * barWidth, height - barHeight, barWidth, barHeight);
            }
            
            // Draw axes
            ctx.strokeStyle = '#374151';
            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.moveTo(0, height);
            ctx.lineTo(width, height);
            ctx.moveTo(0, 0);
            ctx.lineTo(0, height);
            ctx.stroke();
        }
        
        function drawGrayscaleHistogram(data) {
            const canvas = document.getElementById('grayHistogramCanvas');
            if (!canvas) return;
            
            const ctx = canvas.getContext('2d');
            const width = canvas.width;
            const height = canvas.height;
            
            // Clear canvas
            ctx.clearRect(0, 0, width, height);
            
            // Find maximum value for scaling
            const maxValue = Math.max(...data.grayscaleHistogram);
            if (maxValue === 0) return;
            
            // Draw histogram
            const barWidth = width / 256;
            ctx.fillStyle = 'rgba(75, 85, 99, 0.8)';
            
            for (let i = 0; i < 256; i++) {
                const barHeight = (data.grayscaleHistogram[i] / maxValue) * height;
                ctx.fillRect(i * barWidth, height - barHeight, barWidth, barHeight);
            }
            
            // Draw threshold line
            ctx.strokeStyle = '#ef4444';
            ctx.lineWidth = 2;
            ctx.beginPath();
            const thresholdX = (data.threshold / 255) * width;
            ctx.moveTo(thresholdX, 0);
            ctx.lineTo(thresholdX, height);
            ctx.stroke();
            
            // Draw axes
            ctx.strokeStyle = '#374151';
            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.moveTo(0, height);
            ctx.lineTo(width, height);
            ctx.moveTo(0, 0);
            ctx.lineTo(0, height);
            ctx.stroke();
        }
        
        function performEdgeDetection(operation) {
            if (!currentImageData) {
                alert('Please upload an image first');
                return;
            }
            
            showLoading(true);
            
            fetch('/edge-detection', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    operation: operation
                })
            })
            .then(response => response.json())
            .then(data => {
                showLoading(false);
                if (data.success) {
                    displayEdgeDetectionResults(data.data, operation);
                } else {
                    alert('Edge detection failed: ' + data.message);
                }
            })
            .catch(error => {
                showLoading(false);
                alert('Edge detection error: ' + error.message);
            });
        }
        
        function displayEdgeDetectionResults(data, operation) {
            // Display the main edge detection result
            if (data.edgeImage) {
                displayImage(data.edgeImage);
            }
            
            // Create modal window for detailed results
            const modal = document.createElement('div');
            modal.style.cssText = ` + "`" + `
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: rgba(0, 0, 0, 0.8);
                z-index: 1000;
                display: flex;
                align-items: center;
                justify-content: center;
                backdrop-filter: blur(4px);
            ` + "`" + `;
            
            const modalContent = document.createElement('div');
            modalContent.style.cssText = ` + "`" + `
                background: white;
                border-radius: 12px;
                padding: 2rem;
                max-width: 95vw;
                max-height: 95vh;
                overflow-y: auto;
                box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
                position: relative;
            ` + "`" + `;
            
            let resultsHTML = ` + "`" + `
                <div style="font-family: 'Inter', sans-serif;">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; border-bottom: 2px solid #e2e8f0; padding-bottom: 1rem;">
                        <h2 style="margin: 0; color: #1e293b; font-size: 1.5rem;">🔍 Sobel Edge Detection Results</h2>
                        <button onclick="this.closest('.modal').remove()" style="background: #ef4444; color: white; border: none; border-radius: 50%; width: 2rem; height: 2rem; cursor: pointer; font-size: 1rem;">×</button>
                    </div>
                    
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-bottom: 2rem;">
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">📐 Sobel Kernels</h3>
            ` + "`" + `;
            
            // Display kernels based on operation type
            if (operation === 'sobel') {
                // Show both kernels for full Sobel
                resultsHTML += ` + "`" + `
                            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                                <div>
                                    <h4 style="text-align: center; margin-bottom: 0.5rem; color: #6b7280;">Sobel X (Vertical Edges)</h4>
                                    <div id="sobelXKernel" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; max-width: 120px; margin: 0 auto; border: 2px solid #3b82f6; padding: 8px; border-radius: 8px; background: #eff6ff;"></div>
                                </div>
                                <div>
                                    <h4 style="text-align: center; margin-bottom: 0.5rem; color: #6b7280;">Sobel Y (Horizontal Edges)</h4>
                                    <div id="sobelYKernel" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; max-width: 120px; margin: 0 auto; border: 2px solid #10b981; padding: 8px; border-radius: 8px; background: #ecfdf5;"></div>
                                </div>
                            </div>
                ` + "`" + `;
            } else if (operation === 'sobelX') {
                // Show only Sobel X kernel
                resultsHTML += ` + "`" + `
                            <div style="display: flex; justify-content: center;">
                                <div>
                                    <h4 style="text-align: center; margin-bottom: 0.5rem; color: #6b7280;">Sobel X (Vertical Edges)</h4>
                                    <div id="sobelXKernel" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; max-width: 120px; margin: 0 auto; border: 2px solid #3b82f6; padding: 8px; border-radius: 8px; background: #eff6ff;"></div>
                                </div>
                            </div>
                ` + "`" + `;
            } else if (operation === 'sobelY') {
                // Show only Sobel Y kernel
                resultsHTML += ` + "`" + `
                            <div style="display: flex; justify-content: center;">
                                <div>
                                    <h4 style="text-align: center; margin-bottom: 0.5rem; color: #6b7280;">Sobel Y (Horizontal Edges)</h4>
                                    <div id="sobelYKernel" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; max-width: 120px; margin: 0 auto; border: 2px solid #10b981; padding: 8px; border-radius: 8px; background: #ecfdf5;"></div>
                                </div>
                            </div>
                ` + "`" + `;
            }
            
            resultsHTML += ` + "`" + `
                        </div>
                        <div>
                            <h3 style="color: #374151; margin-bottom: 1rem;">📊 Edge Detection Results</h3>
                            <div style="display: grid; gap: 1rem;">
            ` + "`" + `;
            
            // Add images based on operation
            if (operation === 'sobel') {
                resultsHTML += ` + "`" + `
                    <div style="text-align: center;">
                        <h4 style="margin-bottom: 0.5rem;">Combined Edge Detection</h4>
                        <img src="data:image/png;base64,` + "${data.edgeImage}" + `" style="max-width: 200px; border: 1px solid #e2e8f0; border-radius: 6px;" alt="Combined Edges">
                    </div>
                ` + "`" + `;
                if (data.magnitudeImage) {
                    resultsHTML += ` + "`" + `
                        <div style="text-align: center;">
                            <h4 style="margin-bottom: 0.5rem;">Magnitude Image</h4>
                            <img src="data:image/png;base64,` + "${data.magnitudeImage}" + `" style="max-width: 200px; border: 1px solid #e2e8f0; border-radius: 6px;" alt="Magnitude">
                        </div>
                    ` + "`" + `;
                }
            } else {
                resultsHTML += ` + "`" + `
                    <div style="text-align: center;">
                        <h4 style="margin-bottom: 0.5rem;">` + "${operation.toUpperCase()}" + ` Result</h4>
                        <img src="data:image/png;base64,` + "${data.edgeImage}" + `" style="max-width: 250px; border: 1px solid #e2e8f0; border-radius: 6px;" alt="` + "${operation}" + ` Edges">
                    </div>
                ` + "`" + `;
            }
            
            resultsHTML += ` + "`" + `
                            </div>
                        </div>
                    </div>
            ` + "`" + `;
            
            // Add convolution matrices
            if (operation === 'sobel' && data.sobelXResult && data.sobelYResult) {
                resultsHTML += ` + "`" + `
                    <div style="margin-bottom: 2rem;">
                        <h3 style="color: #374151; margin-bottom: 1rem;">🔢 Convolution Matrices</h3>
                        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem;">
                            <div>
                                <h4 style="margin-bottom: 1rem; color: #6b7280;">Sobel X Convolution Matrix</h4>
                                <div style="max-height: 200px; overflow: auto; border: 1px solid #e2e8f0; border-radius: 6px;">
                                    <table style="width: 100%; border-collapse: collapse; font-size: 0.75rem; font-family: monospace;">
                ` + "`" + `;
                
                // Add Sobel X matrix (first 10x10 for performance)
                const xMatrix = data.sobelXResult.convolutionData;
                const maxRows = Math.min(xMatrix.length, 10);
                const maxCols = Math.min(xMatrix[0] ? xMatrix[0].length : 0, 10);
                
                for (let i = 0; i < maxRows; i++) {
                    resultsHTML += '<tr>';
                    for (let j = 0; j < maxCols; j++) {
                        const value = xMatrix[i] && xMatrix[i][j] !== undefined ? xMatrix[i][j].toFixed(1) : '0.0';
                        resultsHTML += ` + "`<td style=\"padding: 2px 4px; border: 1px solid #e2e8f0; text-align: right;\">${value}</td>`" + `;
                    }
                    resultsHTML += '</tr>';
                }
                
                resultsHTML += ` + "`" + `
                                    </table>
                                </div>
                            </div>
                            <div>
                                <h4 style="margin-bottom: 1rem; color: #6b7280;">Sobel Y Convolution Matrix</h4>
                                <div style="max-height: 200px; overflow: auto; border: 1px solid #e2e8f0; border-radius: 6px;">
                                    <table style="width: 100%; border-collapse: collapse; font-size: 0.75rem; font-family: monospace;">
                ` + "`" + `;
                
                // Add Sobel Y matrix (first 10x10 for performance)  
                const yMatrix = data.sobelYResult.convolutionData;
                
                for (let i = 0; i < maxRows; i++) {
                    resultsHTML += '<tr>';
                    for (let j = 0; j < maxCols; j++) {
                        const value = yMatrix[i] && yMatrix[i][j] !== undefined ? yMatrix[i][j].toFixed(1) : '0.0';
                        resultsHTML += ` + "`<td style=\"padding: 2px 4px; border: 1px solid #e2e8f0; text-align: right;\">${value}</td>`" + `;
                    }
                    resultsHTML += '</tr>';
                }
                
                resultsHTML += ` + "`" + `
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                ` + "`" + `;
            } else if ((operation === 'sobelX' && data.sobelXResult) || (operation === 'sobelY' && data.sobelYResult)) {
                const result = operation === 'sobelX' ? data.sobelXResult : data.sobelYResult;
                resultsHTML += ` + "`" + `
                    <div style="margin-bottom: 2rem;">
                        <h3 style="color: #374151; margin-bottom: 1rem;">🔢 Convolution Matrix (` + "${operation.toUpperCase()}" + `)</h3>
                        <div style="max-height: 300px; overflow: auto; border: 1px solid #e2e8f0; border-radius: 6px;">
                            <table style="width: 100%; border-collapse: collapse; font-size: 0.75rem; font-family: monospace;">
                ` + "`" + `;
                
                const matrix = result.convolutionData;
                const maxRows = Math.min(matrix.length, 15);
                const maxCols = Math.min(matrix[0] ? matrix[0].length : 0, 15);
                
                for (let i = 0; i < maxRows; i++) {
                    resultsHTML += '<tr>';
                    for (let j = 0; j < maxCols; j++) {
                        const value = matrix[i] && matrix[i][j] !== undefined ? matrix[i][j].toFixed(1) : '0.0';
                        resultsHTML += ` + "`<td style=\"padding: 3px 6px; border: 1px solid #e2e8f0; text-align: right;\">${value}</td>`" + `;
                    }
                    resultsHTML += '</tr>';
                }
                
                resultsHTML += ` + "`" + `
                            </table>
                        </div>
                    </div>
                ` + "`" + `;
            }
            
            resultsHTML += ` + "`" + `
                    <div style="text-align: center;">
                        <button onclick="this.closest('.modal').remove()" style="background: #3b82f6; color: white; border: none; padding: 0.75rem 2rem; border-radius: 6px; cursor: pointer; font-size: 1rem; font-weight: 500;">Close Results</button>
                    </div>
                </div>
            ` + "`" + `;
            
            modalContent.innerHTML = resultsHTML;
            modal.appendChild(modalContent);
            modal.className = 'modal';
            document.body.appendChild(modal);
            
            // Display kernels
            setTimeout(() => {
                if (data.sobelXKernel) {
                    displayKernelInModal('sobelXKernel', data.sobelXKernel);
                }
                if (data.sobelYKernel) {
                    displayKernelInModal('sobelYKernel', data.sobelYKernel);
                }
            }, 100);
        }
        
        function displayKernelInModal(containerId, kernel) {
            const container = document.getElementById(containerId);
            if (!container || !kernel) return;
            
            container.innerHTML = '';
            kernel.forEach(row => {
                row.forEach(value => {
                    const cell = document.createElement('div');
                    cell.style.cssText = 'background: white; border: 1px solid #d1d5db; padding: 6px; text-align: center; font-weight: bold; font-size: 0.75rem; border-radius: 4px;';
                    cell.textContent = value;
                    container.appendChild(cell);
                });
            });
        }
        
        function showLoading(show) {
            document.getElementById('loadingOverlay').style.display = show ? 'flex' : 'none';
        }
        
        // Initialize on page load
        document.addEventListener('DOMContentLoaded', init);
    </script>
</body>
</html>`
	
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, nil)
}