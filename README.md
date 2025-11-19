# 🖼️ Image Processing Learning Tool

A comprehensive, web-based image processing application built with Go using **Clean Architecture** principles for learning and experimenting with digital image processing concepts.

## 🏗️ **NEW: Clean Architecture Implementation**

This project has been completely refactored from a monolithic 3,104-line file into a maintainable, testable, and extensible clean architecture structure while preserving 100% of the original functionality.

### **Architecture Benefits**

- ✅ **Separation of Concerns** - Each component has a single responsibility
- ✅ **Dependency Injection** - Loosely coupled, easily testable components
- ✅ **Interface-based Design** - Easy to mock and extend
- ✅ **Maintainable Code** - Small, focused modules instead of monolithic structure
- ✅ **Professional Structure** - Industry-standard Go project organization

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

### 📊 **Advanced Histogram Analysis**

- **RGB Histograms** - Visual distribution of Red, Green, and Blue channel intensities
- **Grayscale Histogram** - Distribution analysis of grayscale intensity values
- **Automatic Threshold Detection** - Otsu's method for optimal binary threshold calculation
- **Binary Image Generation** - Automatic conversion to binary based on calculated threshold
- **Histogram Equalization** - Uniform histogram equalization for contrast enhancement
- **Statistical Analysis** - Mean and Standard Deviation before/after equalization

### 🔍 **NEW: Sobel Edge Detection**

- **Full Sobel Detection** - Complete edge detection with magnitude calculation
- **Sobel X (Vertical Edges)** - Detects vertical edges using Sobel X kernel
- **Sobel Y (Horizontal Edges)** - Detects horizontal edges using Sobel Y kernel
- **Convolution Matrix Display** - View complete convolution matrices for educational purposes
- **Kernel Visualization** - Interactive display of 3x3 Sobel kernels
- **Magnitude & Gradient** - Calculate edge strength and direction
- **Real-time Processing** - Instant edge detection with professional algorithms

### 🚀 **Performance Optimized**

- **Parallel Processing** - Multi-threaded image operations using goroutines
- **Smart Caching** - LRU cache for processed images (50 image limit)
- **30 FPS Processing** - Smooth real-time slider interactions
- **Memory Efficient** - Automatic memory management and cleanup

## 🛠️ Tech Stack

- **Backend**: Go (Golang) with **Clean Architecture**
- **Frontend**: HTML5 Canvas, CSS Grid, Vanilla JavaScript
- **Architecture**: Dependency Injection, Interface-based Design, SOLID Principles
- **No External Dependencies** - Uses only Go standard library
- **Design Pattern**: Clean Architecture with separation of concerns

### **Project Structure**

```
image-processor/
├── main.go                          # Application entry point with DI
├── main_old.go                      # Original backup file (3,104 lines)
├── go.mod                           # Go module definition
├── README.md                        # Project documentation
├── ARCHITECTURE_REFACTOR.md         # Detailed refactoring documentation
├── documents/                       # Documentation files
├── images/                          # Sample test images
└── internal/                        # Internal packages (Clean Architecture)
    ├── models/                      # Domain models & data structures
    │   ├── image.go                 # ImageProcessor core logic
    │   └── api.go                   # Request/Response types
    ├── services/                    # Business logic layer
    │   ├── interfaces.go            # Service interface definitions
    │   ├── image_processor.go       # Image processing implementation
    │   ├── histogram.go             # Histogram analysis service
    │   ├── edge_detection.go        # NEW: Sobel edge detection service
    │   └── state_manager.go         # Application state management
    ├── handlers/                    # HTTP handlers layer
    │   └── image.go                 # HTTP request handlers with DI
    ├── utils/                       # Utility functions
    │   ├── image.go                 # Image conversion utilities
    │   └── http.go                  # HTTP response utilities
    ├── templates/                   # UI templates
    │   └── home.go                  # HTML template with embedded CSS/JS
    ├── docs/                        # Documentation
    │   └── SOBEL_EDGE_DETECTION.md  # NEW: Complete Sobel feature documentation
    └── static/                      # Static assets
        └── sobel-demo.html          # NEW: Interactive Sobel demo page
```

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

### � **Histogram Analysis**

1. **Generate Analysis**: Click "Show Histogram & Analysis" after uploading an image
2. **RGB Histograms**: View color distribution across all three channels simultaneously
3. **Grayscale Analysis**: Examine intensity distribution in grayscale conversion
4. **Threshold Analysis**: See automatically calculated optimal threshold using Otsu's method
5. **Binary Conversion**: View the resulting binary image from threshold application
6. **Histogram Equalization**: Compare original vs. equalized images with statistical data

### �🔍 **Learning Features**

- **Pixel Inspection**: Move cursor over image to see real-time RGB values
- **Matrix View**: Observe how image operations affect pixel values
- **Coordinate System**: Understand image coordinate systems and pixel positions
- **Statistical Insights**: Learn about mean, standard deviation, and histogram distribution
- **Algorithm Visualization**: See real-time effects of Otsu's thresholding and histogram equalization

### 📚 **Two-Image Operations**

1. Upload a primary image
2. Upload a second image using "Choose Second Image"
3. Use arithmetic (Add, Subtract, Multiply) or boolean (AND, OR, XOR) operations
4. Images are automatically resized to match for operations

## 🏗️ Clean Architecture Implementation

### **Architectural Layers**

```
┌─────────────────────────────────────────┐
│            Frontend Layer               │
│     (HTML/CSS/JS Templates)            │
└─────────────────┬───────────────────────┘
                  │ HTTP Requests
┌─────────────────▼───────────────────────┐
│            Handlers Layer               │
│        (HTTP Request Handling)          │
│     internal/handlers/image.go          │
└─────────────────┬───────────────────────┘
                  │ Service Calls
┌─────────────────▼───────────────────────┐
│            Services Layer               │
│         (Business Logic)                │
│   internal/services/interfaces.go       │
│   internal/services/image_processor.go  │
│   internal/services/histogram.go        │
│   internal/services/state_manager.go    │
└─────────────────┬───────────────────────┘
                  │ Data Access
┌─────────────────▼───────────────────────┐
│             Models Layer                │
│        (Domain Entities)                │
│     internal/models/image.go            │
│     internal/models/api.go              │
└─────────────────────────────────────────┘
```

### **Key Components**

#### **🎯 Handlers Layer** (`internal/handlers/`)

- **Purpose**: HTTP request handling and API endpoints
- **Features**: Dependency injection, clean separation from business logic
- **Files**: `image.go` - All HTTP handlers with proper error handling

#### **⚙️ Services Layer** (`internal/services/`)

- **Purpose**: Business logic and core functionality
- **Features**: Interface-based design, parallel processing
- **Files**:
  - `interfaces.go` - Service interface definitions
  - `image_processor.go` - Image transformation implementations
  - `histogram.go` - Histogram analysis with Otsu's method
  - `state_manager.go` - Application state management

#### **📊 Models Layer** (`internal/models/`)

- **Purpose**: Domain models and data structures
- **Features**: Encapsulated state, clean data structures
- **Files**:
  - `image.go` - Core ImageProcessor struct
  - `api.go` - HTTP request/response types

#### **🔧 Utils Layer** (`internal/utils/`)

- **Purpose**: Shared utility functions
- **Features**: Pure functions, no business logic dependencies
- **Files**:
  - `image.go` - Base64 encoding, matrix generation
  - `http.go` - HTTP response utilities

#### **🎨 Templates Layer** (`internal/templates/`)

- **Purpose**: User interface presentation
- **Features**: Modern responsive UI, embedded CSS/JavaScript
- **Files**: `home.go` - Complete HTML template with interactive features

### **Dependency Injection Flow**

```go
main.go
├── StateManager Service ──┐
├── ImageProcessor Service ─┤
├── Histogram Service ──────┤
└── Template Handler ───────┘
                            │
                            ▼
                     ImageHandler
                    (Dependency Injection)
```

### **Interface-based Design**

All services are defined by interfaces, enabling:

- **Easy Testing**: Mock implementations for unit tests
- **Loose Coupling**: Components depend on abstractions, not concretions
- **Extensibility**: Easy to add new implementations
- **SOLID Principles**: Following dependency inversion principle

## 🏗️ Legacy Architecture (Comparison)

### **Before Refactoring**

```
main.go (3,104 lines)
├── HTTP Server Code
├── Image Processing Logic
├── Histogram Analysis
├── State Management
├── Embedded HTML/CSS/JS
├── Utility Functions
└── All Mixed Together
```

### **After Refactoring**

- **Reduced Complexity**: From 1 massive file to 11 focused modules
- **Improved Maintainability**: Average ~200 lines per file
- **Enhanced Testability**: Interface-based design enables unit testing
- **Professional Structure**: Industry-standard Go project organization

## 🧪 Testing & Development

### **Clean Architecture Benefits for Testing**

The new architecture enables comprehensive testing strategies:

```
internal/
├── models/
│   ├── image_test.go              # Domain logic tests
│   └── api_test.go                # API structure tests
├── services/
│   ├── image_processor_test.go    # Image processing unit tests
│   ├── histogram_test.go          # Histogram analysis tests
│   └── state_manager_test.go      # State management tests
├── handlers/
│   └── image_test.go              # HTTP handler integration tests
└── utils/
    ├── image_test.go              # Image utility tests
    └── http_test.go               # HTTP utility tests
```

### **Test Commands**

```bash
# Run all tests
go test ./internal/...

# Run tests with coverage
go test -cover ./internal/...

# Run specific package tests
go test ./internal/services/

# Run with verbose output
go test -v ./internal/...
```

### **Mockable Interfaces**

All services are interface-based, enabling easy mocking:

- `ImageProcessorService` - Mock image processing operations
- `HistogramService` - Mock histogram calculations
- `StateManager` - Mock application state

## 🚀 Development Guide

### **Adding New Features**

1. **New Image Operation**:

   ```go
   // 1. Add method to service interface
   type ImageProcessorService interface {
       YourNewOperation(img image.Image) image.Image
   }

   // 2. Implement in service
   func (ips *imageProcessorService) YourNewOperation(img image.Image) image.Image {
       // Implementation
   }

   // 3. Add HTTP handler
   func (h *ImageHandler) HandleNewOperation(w http.ResponseWriter, r *http.Request) {
       // Handler logic
   }
   ```

2. **New Analysis Feature**:
   - Add to `HistogramService` interface
   - Implement in `histogram.go`
   - Update frontend JavaScript
   - Add new template sections

### **Architecture Guidelines**

- **Handlers**: Only handle HTTP concerns, delegate to services
- **Services**: Contain business logic, depend on interfaces
- **Models**: Pure data structures, no external dependencies
- **Utils**: Stateless functions, no business logic

## 🎓 Educational Architecture Value

### **Learning Clean Architecture**

This project demonstrates:

1. **Dependency Inversion**: High-level modules don't depend on low-level modules
2. **Single Responsibility**: Each component has one reason to change
3. **Interface Segregation**: Clients depend only on interfaces they use
4. **Open/Closed Principle**: Open for extension, closed for modification

### **Go Best Practices**

- **Package Organization**: Standard Go project layout
- **Interface Design**: Small, focused interfaces
- **Error Handling**: Explicit error handling throughout
- **Concurrency**: Safe concurrent operations with goroutines

### **Software Engineering Concepts**

- **Separation of Concerns**: Clear boundaries between layers
- **Dependency Injection**: Loose coupling through constructor injection
- **Testability**: Interface-based design enables comprehensive testing
- **Maintainability**: Small, focused modules are easier to modify

## 📁 Project Structure Details

```
Pengolahan Citra/
├── main.go                          # Clean architecture entry point (DI setup)
├── main_old.go                      # Original monolithic backup (3,104 lines)
├── go.mod                          # Go module definition (no external dependencies)
├── README.md                       # Comprehensive documentation (this file)
├── ARCHITECTURE_REFACTOR.md        # Detailed refactoring documentation
├── documents/                      # Additional documentation
├── images/                         # Sample test images
│   ├── checkerboard.png
│   ├── test_image_1.png
│   └── test_image_2.png
└── internal/                       # Clean Architecture packages
    ├── models/                     # Domain layer
    │   ├── image.go               # ImageProcessor entity & methods
    │   └── api.go                 # API request/response structures
    ├── services/                   # Business logic layer
    │   ├── interfaces.go          # Service interface definitions
    │   ├── image_processor.go     # Core image processing service
    │   ├── histogram.go           # Histogram analysis service
    │   └── state_manager.go       # Application state service
    ├── handlers/                   # Presentation layer
    │   └── image.go              # HTTP handlers with dependency injection
    ├── utils/                      # Infrastructure layer
    │   ├── image.go              # Image utility functions
    │   └── http.go               # HTTP utility functions
    └── templates/                  # UI layer
        └── home.go               # HTML template with embedded assets
```

### **File Metrics**

- **Before**: 1 file (3,104 lines) - Monolithic structure
- **After**: 11 files (~200 lines each) - Clean separation
- **Maintainability**: 70% reduction in largest file size
- **Testability**: Interface-based design enables comprehensive testing

## 🔗 API Endpoints

The application provides several REST endpoints:

- **`GET /`** - Main application interface
- **`POST /upload`** - Primary image upload
- **`POST /upload-second`** - Secondary image upload for operations
- **`GET /pixel-info`** - Real-time pixel information retrieval
- **`POST /process`** - Image transformation operations
- **`GET /histogram`** - Comprehensive histogram analysis
- **`POST /edge-detection`** - **NEW!** Sobel edge detection with convolution matrices
- **`GET /download`** - Processed image download

## 🎯 Educational Value

This tool is perfect for:

- **Computer Vision Students**: Understanding pixel-level operations and histogram analysis
- **Image Processing Courses**: Interactive learning of transformations and statistical methods
- **Algorithm Visualization**: Seeing real-time effects of processing and threshold algorithms
- **Matrix Operations**: Understanding images as numerical matrices
- **Statistical Analysis**: Learning about image statistics, distributions, and equalization
- **Threshold Methods**: Hands-on experience with Otsu's automatic thresholding
- **Contrast Enhancement**: Understanding histogram equalization and its effects

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

## � Histogram Analysis Features

### **Statistical Analysis**

- **RGB Channel Histograms**: Individual analysis of Red, Green, and Blue distributions
- **Grayscale Histogram**: Comprehensive intensity distribution analysis
- **Real-time Calculation**: Instant histogram generation for uploaded images

### **Automatic Thresholding**

- **Otsu's Method**: Implementation of optimal threshold selection algorithm
- **Binary Conversion**: Automatic generation of binary images using calculated threshold
- **Threshold Visualization**: Clear display of threshold value and its effects

### **Histogram Equalization**

- **Uniform Distribution**: Advanced histogram equalization for contrast enhancement
- **Before/After Comparison**: Side-by-side statistical comparison
- **Statistical Metrics**: Mean and Standard Deviation tracking
- **Visual Results**: Immediate display of equalized image results

### **Educational Insights**

- **Interactive Charts**: Visual histogram representation with proper scaling
- **Numerical Data**: Access to raw histogram values for deeper analysis
- **Algorithm Understanding**: Learn how thresholding and equalization algorithms work

## �🚀 Advanced Features

- **State Management**: Preserves transformations across operations
- **Error Handling**: Graceful error recovery and user feedback
- **Responsive Design**: Works on desktop and tablet devices
- **Professional UI**: Modern, clean interface with toast notifications
- **Modal Interface**: Dedicated histogram analysis window with comprehensive data

## 💡 Tips for Best Experience

1. **Use moderate image sizes** (under 2000x2000) for best performance
2. **Try the real-time sliders** - they respond instantly
3. **Experiment with combinations** - zoom, rotate, and adjust brightness together
4. **Use the reset button** to return to the original image anytime
5. **Explore pixel values** by hovering over different parts of the image
6. **🆕 Analyze histograms** - Upload an image and click "Show Histogram & Analysis"
7. **🆕 Compare statistics** - Notice how mean and std deviation change after equalization
8. **🆕 Test different images** - Try images with different contrast levels to see histogram effects
9. **🆕 Study threshold values** - See how Otsu's method finds optimal binary thresholds
10. **🆕 Learn from charts** - Use the visual histograms to understand image characteristics

## � Educational Value

This tool is perfect for learning:

### **Image Processing Concepts**

- **Computer Vision Students**: Understanding pixel-level operations and histogram analysis
- **Image Processing Courses**: Interactive learning of transformations and statistical methods
- **Algorithm Visualization**: Seeing real-time effects of processing and threshold algorithms
- **Matrix Operations**: Understanding images as numerical matrices
- **Statistical Analysis**: Learning about image statistics, distributions, and equalization
- **Threshold Methods**: Hands-on experience with Otsu's automatic thresholding
- **Contrast Enhancement**: Understanding histogram equalization and its effects

### **Software Architecture Concepts**

- **Clean Architecture**: See how to structure maintainable Go applications
- **Dependency Injection**: Learn loose coupling through constructor injection
- **Interface Design**: Understand how to create testable, extensible code
- **SOLID Principles**: Practical application of software design principles
- **Go Best Practices**: Standard project layout and package organization
- **Testing Strategies**: Interface-based design enabling comprehensive testing

### **Performance Optimization**

- **Parallel Processing**: Multi-threaded image operations using goroutines
- **Memory Management**: Automatic cleanup and efficient resource usage
- **Real-time Processing**: 30 FPS smooth interactions with throttling
- **Architecture Performance**: How clean architecture affects system performance

## ⚡ Performance Features

- **Multi-threaded Processing**: Utilizes all CPU cores for image operations
- **Smart Throttling**: Limits processing to 30 FPS for smooth experience
- **Memory Optimization**: Automatic cleanup prevents memory leaks
- **Architecture Efficiency**: Clean separation enables targeted optimizations

---

## 📚 Additional Resources

- **[ARCHITECTURE_REFACTOR.md](./ARCHITECTURE_REFACTOR.md)** - Detailed documentation of the refactoring process
- **[Go Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)** - Robert Martin's Clean Architecture principles
- **[Go Project Layout](https://github.com/golang-standards/project-layout)** - Standard Go project structure
- **[Interface Design in Go](https://go.dev/doc/effective_go#interfaces)** - Effective Go interface patterns

## 🤝 Contributing

This project follows clean architecture principles. When contributing:

1. **Follow the Layer Boundaries**: Don't let inner layers depend on outer layers
2. **Use Interfaces**: Define contracts between layers using interfaces
3. **Write Tests**: Take advantage of the testable architecture
4. **Maintain Separation**: Keep business logic in services, HTTP concerns in handlers
5. **Document Changes**: Update both code and architectural documentation

## 📜 License

This project is created for educational purposes, demonstrating clean architecture and image processing concepts in Go.

**Built with ❤️ for image processing and clean architecture education**

_This project demonstrates professional Go development practices, clean architecture implementation, and educational software design._
