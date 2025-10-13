# 🖼️ Image Processing Learning Tool

A comprehensive, web-based image processing application built with Go for learning and experimenting with digital image processing concepts.

## ✨ Features

### 🎛️ **Real-Time Interactive Processing**

- **Brightness Control** - Adjust image brightness with real-time slider
- **Zoom Control** - Zoom in/out with smooth scaling (0.1x to 5.0x)
- **Rotation Control** - Rotate images with precise angle control (-360° to +360°)

### 🔧 **Image Transformations**

- **Color Conversions**
  - Grayscale conversion
  - Binary threshold conversion
- **Geometric Operations**
  - Horizontal flip
  - Vertical flip
  - Image translation
- **Arithmetic Operations** (Two-image operations)
  - Addition
  - Subtraction
  - Multiplication
- **Boolean Operations** (Two-image operations)
  - AND operation
  - OR operation
  - XOR operation

### 📊 **Educational Features**

- **Real-Time Pixel Inspection** - Click anywhere on the image to see RGB values
- **Complete Pixel Matrix Display** - View the entire image as a matrix of pixel values
- **Interactive Matrix Highlighting** - Cursor position highlights corresponding matrix cell
- **Image Information Display** - Shows dimensions, coordinates, and color values

### 🚀 **Performance Optimized**

- **Parallel Processing** - Multi-threaded image operations using goroutines
- **Smart Caching** - LRU cache for processed images (50 image limit)
- **30 FPS Processing** - Smooth real-time slider interactions
- **Memory Efficient** - Automatic memory management and cleanup

## 🛠️ Tech Stack

- **Backend**: Go (Golang) with standard libraries only
- **Frontend**: HTML5 Canvas, CSS Grid, Vanilla JavaScript
- **Architecture**: Clean architecture with separation of concerns
- **No External Dependencies** - Uses only Go standard library

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Web browser (Chrome, Firefox, Safari, Edge)

### Installation & Running

1. **Clone or download the project**

   ```bash
   cd "path/to/Pengolahan Citra"
   ```

2. **Run the application**

   ```bash
   go run main.go
   ```

3. **Open your web browser**

   - Navigate to `http://localhost:8080`
   - The application will start automatically

4. **Start processing images**
   - Upload an image using the file picker
   - Experiment with real-time sliders
   - Try different image operations
   - View pixel matrices and RGB values

## 📖 How to Use

### 🖼️ **Basic Operations**

1. **Upload Image**: Click "Choose Image File" to upload PNG, JPEG, or GIF files
2. **Real-Time Adjustments**: Use sliders for brightness, zoom, and rotation
3. **Apply Filters**: Click buttons for grayscale, binary, or flip operations
4. **Reset**: Return to original image at any time
5. **Download**: Save processed images to your computer

### 🔍 **Learning Features**

- **Pixel Inspection**: Move cursor over image to see real-time RGB values
- **Matrix View**: Observe how image operations affect pixel values
- **Coordinate System**: Understand image coordinate systems and pixel positions

### 📚 **Two-Image Operations**

1. Upload a primary image
2. Upload a second image using "Choose Second Image"
3. Use arithmetic (Add, Subtract, Multiply) or boolean (AND, OR, XOR) operations
4. Images are automatically resized to match for operations

## 🏗️ Architecture

### **Clean Architecture Principles**

```
├── Frontend (HTML/CSS/JS)
│   ├── User Interface Layer
│   ├── Real-time Processing
│   └── Canvas Rendering
│
├── Backend (Go)
│   ├── HTTP Server & Routing
│   ├── Image Processing Engine
│   ├── Transformation Pipeline
│   └── Pixel Matrix Generation
│
└── Performance Layer
    ├── Parallel Processing
    ├── Memory Management
    └── Caching System
```

### **Key Components**

- **ImageProcessor**: Core struct managing image state and transformations
- **Transformation Pipeline**: Applies multiple operations in optimal order
- **Cache System**: LRU cache for performance optimization
- **Matrix Generator**: Creates educational pixel matrices

## ⚡ Performance Features

- **Multi-threaded Processing**: Utilizes all CPU cores for image operations
- **Smart Throttling**: Limits processing to 30 FPS for smooth experience
- **Memory Optimization**: Automatic cleanup prevents memory leaks
- **Cache Management**: Intelligent caching reduces redundant calculations

## 📁 Project Structure

```
Pengolahan Citra/
├── main.go           # Complete application (web server + image processing)
├── go.mod           # Go module definition (no external dependencies)
├── README.md        # This documentation
└── images/          # Sample test images
    ├── checkerboard.png
    ├── test_image_1.png
    └── test_image_2.png
```

## 🎯 Educational Value

This tool is perfect for:

- **Computer Vision Students**: Understanding pixel-level operations
- **Image Processing Courses**: Interactive learning of transformations
- **Algorithm Visualization**: Seeing real-time effects of processing
- **Matrix Operations**: Understanding images as numerical matrices

## 🔧 Supported Image Formats

- **PNG** - Full support with transparency
- **JPEG** - Full color support
- **GIF** - Static image support

## 📊 Technical Specifications

- **Maximum Image Size**: 4000x4000 pixels (auto-scaled for performance)
- **Color Depth**: 8-bit per channel (RGB + Alpha)
- **Processing Speed**: Up to 30 FPS real-time processing
- **Memory Usage**: Optimized with automatic cleanup
- **Browser Compatibility**: All modern browsers

## 🚀 Advanced Features

- **State Management**: Preserves transformations across operations
- **Error Handling**: Graceful error recovery and user feedback
- **Responsive Design**: Works on desktop and tablet devices
- **Professional UI**: Modern, clean interface with toast notifications

## 💡 Tips for Best Experience

1. **Use moderate image sizes** (under 2000x2000) for best performance
2. **Try the real-time sliders** - they respond instantly
3. **Experiment with combinations** - zoom, rotate, and adjust brightness together
4. **Use the reset button** to return to the original image anytime
5. **Explore pixel values** by hovering over different parts of the image

## 🎓 Learning Objectives

By using this tool, you will understand:

- How digital images are represented as pixel matrices
- The effects of various image processing operations
- Real-time processing and optimization techniques
- The relationship between mathematical operations and visual results

---

**Built with ❤️ for image processing education**

_This project demonstrates clean architecture, performance optimization, and educational software design in Go._
