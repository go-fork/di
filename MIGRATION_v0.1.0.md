# Migration Guide - go.fork.vn/di v0.1.0

## Tổng quan

Kể từ phiên bản v0.1.0, package đã được chuyển từ `github.com/Fork/di` sang `go.fork.vn/di`. Đây là hướng dẫn chi tiết để migration code của bạn.

## 🚨 Breaking Changes

### Package Name Change

**Before (v0.0.x):**
```go
import "github.com/Fork/di"
```

**After (v0.1.0+):**
```go
import "go.fork.vn/di"
```

### Mocks Import

**Before:**
```go
import "github.com/Fork/di/mocks"
```

**After:**
```go
import "go.fork.vn/di/mocks"
```

## 📋 Migration Steps

### Step 1: Update go.mod

```bash
# Remove old dependency
go mod edit -droprequire github.com/Fork/di

# Add new dependency  
go get go.fork.vn/di@v0.1.0

# Clean up
go mod tidy
```

### Step 2: Update Import Statements

**Find and replace all imports in your project:**

```bash
# For Unix/macOS/Linux
find . -name "*.go" -type f -exec sed -i '' 's|github.com/Fork/di|go.fork.vn/di|g' {} \;

# For Linux (without '' after -i)
find . -name "*.go" -type f -exec sed -i 's|github.com/Fork/di|go.fork.vn/di|g' {} \;
```

**Or manually update each file:**

```go
// ❌ Old import
import (
    "github.com/Fork/di"
    "github.com/Fork/di/mocks"
)

// ✅ New import  
import (
    "go.fork.vn/di"
    "go.fork.vn/di/mocks"
)
```

### Step 3: Update Documentation

Update any documentation, README files, or comments that reference the old package name.

### Step 4: Verify Changes

```bash
# Build to check for errors
go build ./...

# Run tests
go test ./...

# Verify dependencies
go mod verify
```

## 🔧 No API Changes

**Good news!** Không có thay đổi nào về API. Tất cả methods, interfaces, và functionalities đều giữ nguyên:

- `Container` methods: `Bind()`, `Singleton()`, `Make()`, `MustMake()`, etc.
- `ServiceProvider` interface: `Register()`, `Boot()`, `Requires()`, `Providers()`
- `Application` interface: Tất cả methods giữ nguyên
- `ModuleLoaderContract`: Không thay đổi
- Mock objects: Cùng API, chỉ đổi import path

## 📝 Example Migration

### Before (v0.0.x)

```go
package main

import (
    "fmt"
    "github.com/Fork/di"
    "github.com/Fork/di/mocks"
)

type MyService struct {
    Name string
}

func main() {
    container := di.New()
    
    container.Bind("service", func(c *di.Container) interface{} {
        return &MyService{Name: "example"}
    })
    
    service := container.MustMake("service").(*MyService)
    fmt.Println(service.Name)
}
```

### After (v0.1.0+)

```go
package main

import (
    "fmt"
    "go.fork.vn/di"           // ✅ Updated import
    "go.fork.vn/di/mocks"     // ✅ Updated import
)

type MyService struct {
    Name string
}

func main() {
    container := di.New()     // ✅ Same API
    
    container.Bind("service", func(c *di.Container) interface{} {
        return &MyService{Name: "example"}
    })
    
    service := container.MustMake("service").(*MyService)
    fmt.Println(service.Name)
}
```

## 🧪 Testing Migration

### Test Import Updates

```go
// ❌ Old test imports
import (
    "testing"
    "github.com/Fork/di"
    "github.com/Fork/di/mocks"
    "github.com/stretchr/testify/assert"
)

// ✅ New test imports
import (
    "testing"
    "go.fork.vn/di"           // ✅ Updated
    "go.fork.vn/di/mocks"     // ✅ Updated
    "github.com/stretchr/testify/assert"
)
```

### Mock Usage (No Changes)

```go
func TestMyService(t *testing.T) {
    // Mock usage remains exactly the same
    mockApp := new(mocks.Application)
    mockContainer := new(mocks.Container)
    
    mockApp.On("Container").Return(mockContainer)
    
    // Same mock API, no changes needed
    assert.NotNil(t, mockApp)
}
```

## 🚀 New Features in v0.1.0

Ngoài việc đổi package name, v0.1.0 còn bao gồm:

### 📚 Complete Vietnamese Documentation

- `docs/container.md` - DI Container documentation
- `docs/application.md` - Application interface
- `docs/provider.md` - ServiceProvider patterns
- `docs/deferred.md` - Deferred operations
- `docs/loader.md` - Module loading
- `docs/README.md` - System overview

### 🎯 Enhanced BindingFunc Documentation

Detailed documentation cho factory functions với practical examples.

### 🏗️ Production Patterns

Documentation bao gồm enterprise-level patterns và best practices.

## ❓ Troubleshooting

### Common Issues

#### 1. Import Path Not Found

**Error:**
```
go: go.fork.vn/di@v0.1.0: reading https://go.fork.vn/di/@v/v0.1.0.mod: 404 Not Found
```

**Solution:**
Đảm bảo bạn đã update và push tag v0.1.0 lên repository chính.

#### 2. Mixed Import Paths

**Error:**
```
import "github.com/Fork/di": cannot find package
```

**Solution:**
Tìm và thay thế tất cả import paths cũ:

```bash
grep -r "github.com/Fork/di" . --include="*.go"
```

#### 3. Go Module Cache

**Error:**
Cache issues với old package name.

**Solution:**
```bash
go clean -modcache
go mod download
```

### Verification Commands

```bash
# Check for old imports
grep -r "github.com/Fork/di" . --include="*.go"

# Verify new imports work
go list -m go.fork.vn/di

# Test build
go build ./...

# Test imports resolve
go mod graph | grep "go.fork.vn/di"
```

## 📞 Support

Nếu gặp vấn đề trong quá trình migration:

1. Check documentation trong thư mục `docs/`
2. Tạo issue trên repository với tag `migration`
3. Đảm bảo bạn đã follow đúng các bước trong guide này

## 🎯 Benefits of Migration

- **Stable Package Name**: `go.fork.vn/di` là package name chính thức
- **Complete Documentation**: Tài liệu tiếng Việt comprehensive
- **Production Ready**: v0.1.0 là first stable release
- **Better Versioning**: Tuân theo semantic versioning đúng cách
- **Long-term Support**: Package name sẽ không thay đổi trong tương lai

---

**Migration completed successfully!** 🎉

Chào mừng bạn đến với `go.fork.vn/di` v0.1.0 - DI Container chính thức của Fork framework!
