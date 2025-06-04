# Changelog

Tất cả các thay đổi đáng chú ý của dự án go.fork.vn/di sẽ được ghi lại trong file này.
Định dạng dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân theo [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Repository structure now organized with releases/ directory
- Documentation moved to releases/next/ for development
- Historical releases archived in releases/vX.X.X/ directories

## [0.1.3] - 2025-06-04

### Added
- **Release Management Automation**: Scripts tự động cho việc quản lý release
  - `scripts/archive_release.sh`: Tự động archive release và tạo development cycle mới
  - `scripts/create_release_templates.sh`: Tạo templates cho development documentation
  - Comprehensive workflow documentation trong `scripts/README.md`

### Changed  
- **Repository Structure Reorganization**: Tổ chức lại cấu trúc documentation
  - Di chuyển release documentation từ root vào `releases/` directory
  - Tạo `releases/next/` cho work-in-progress documentation
  - Archive historical documentation vào `releases/vX.X.X/` directories
  - Xóa symlinks ở root directory để có cấu trúc sạch sẽ
- **Documentation Paths**: Cập nhật tất cả links đến documentation paths mới
  - `MIGRATION.md` → `releases/next/MIGRATION.md`
  - `RELEASE_NOTES.md` → `releases/next/RELEASE_NOTES.md` 
  - `RELEASE_SUMMARY.md` → `releases/next/RELEASE_SUMMARY.md`

### Fixed
- **Clean Root Directory**: Root chỉ chứa source code và core documentation
- **Professional Structure**: Tuân theo Go community best practices

### Documentation
- **Comprehensive Release Workflow**: Complete guide trong `releases/README.md`
- **Automation Documentation**: Detailed scripts usage và workflow
- **Historical Preservation**: Tất cả documentation của các version trước được bảo tồn

### Breaking Changes
- **Documentation Paths**: Links đến migration guides và release notes đã thay đổi
- **Repository Structure**: Developers cần update bookmarks và documentation links

### Migration
- Xem [Migration Guide](releases/next/MIGRATION.md) để biết chi tiết về việc update documentation links

## [0.1.2] - 2025-06-02

### Added
- **Container Interface**: Chuyển đổi container từ struct sang interface
  - Cải thiện khả năng testing với dependency injection
  - Hỗ trợ nhiều container implementations tùy chỉnh
  - Tuân theo Dependency Inversion Principle triệt để hơn
- **Regenerated Mocks**: Cập nhật tất cả mock objects để tương thích với Container interface mới
  - Cập nhật mocks.Container để implement interface mới
  - Cải thiện khả năng mock testing
  - Đảm bảo type safety trong tests

### Documentation
- **Interface Documentation**: Cập nhật tài liệu để phản ánh Container interface mới
  - Cập nhật `docs/container.md` với Container interface
  - Thêm hướng dẫn chuyển đổi và best practices
  - Bổ sung testing patterns với Container interface
- **Migration Guide**: Tạo MIGRATION_v0.1.2.md với hướng dẫn chi tiết về việc cập nhật
- **Release Notes**: Tạo RELEASE_NOTES_v0.1.2.md với thông tin về phiên bản mới

### Technical Improvements
- **API Consistency**: Đảm bảo tất cả interface đều nhất quán và có tài liệu đầy đủ
- **Type Safety**: Cải thiện type safety trong toàn bộ hệ thống DI
- **Testing Support**: Triển khai mẫu để dễ dàng mock Container interface

## [0.1.1] - 2025-06-02

### Documentation
- **Type Safety Improvements**: Migrated all `app interface{}` parameters to `app Application` across documentation
  - Updated `docs/deferred.md` with 17+ method signatures using strongly typed Application interface
  - Improved ServiceProvider and ServiceProviderDeferred interface documentation with proper typing
  - Enhanced type safety examples and best practices
- **Documentation Restructure**: 
  - Replaced `docs/README.md` with `docs/index.md` for better documentation structure
  - Created comprehensive `docs/overview.md` with detailed DI system architecture
  - Added complete migration documentation in `docs/MIGRATION_COMPLETE.md`
- **Enhanced Documentation Quality**:
  - All method signatures now use proper Application interface typing
  - Improved code examples with type-safe implementations
  - Better documentation organization and navigation
  - Maintained backward compatibility while improving type clarity

### Technical Improvements
- **Interface Consistency**: All interfaces now consistently use `Application` type instead of `interface{}`
- **Type Safety**: Enhanced type safety across all ServiceProvider implementations
- **Documentation Coverage**: Comprehensive documentation updates covering all recent changes

## [0.1.0] - 2025-05-30

### Added
- **🎉 First Official Release**: Fork framework DI container v0.1.0
- **Complete Vietnamese Documentation**: 5000+ lines of comprehensive technical documentation
  - `docs/container.md` - DI Container API reference và patterns (500+ lines)
  - `docs/provider.md` - ServiceProvider implementation patterns (1000+ lines) 
  - `docs/deferred.md` - Deferred operations và async processing (800+ lines)
  - `docs/loader.md` - Module loading và dynamic registration (950+ lines)
  - `docs/application.md` - Application interface và integration (600+ lines)
  - `docs/README.md` - System overview và architecture (400+ lines)
- **Enhanced BindingFunc Documentation**: Detailed documentation cho factory functions với practical examples
- **Production Patterns**: Enterprise-level implementation patterns và best practices
- **Migration Guide**: Comprehensive MIGRATION.md với step-by-step instructions
- **Release Documentation**: Detailed RELEASE_NOTES_v0.1.0.md với features overview
- **Testing Support**: Complete mock objects và testing strategies documentation

### Changed
- **Package Name**: Migrated từ `github.com/go-fork/di` sang `go.fork.vn/di`
- **Module Path**: Updated go.mod với new module name
- **Import Paths**: Updated tất cả internal imports và documentation
- **Installation Instructions**: Updated README.md với new package path

### Breaking Changes
- **Package Import Path**: 
  - Old: `import "github.com/go-fork/di"`
  - New: `import "go.fork.vn/di"`
- **Mocks Import Path**:
  - Old: `import "github.com/go-fork/di/mocks"`
  - New: `import "go.fork.vn/di/mocks"`

### Migration
- See [MIGRATION.md](MIGRATION.md) for detailed migration instructions
- No API changes - chỉ cần update import paths
- Automated migration commands provided

### Documentation
- All documentation viết hoàn toàn bằng tiếng Việt
- Comprehensive examples từ basic đến enterprise-level
- Production-ready patterns và architectural guidance
- Complete troubleshooting và performance considerations

### Technical Improvements
- Enhanced BindingFunc documentation với usage patterns
- Improved concurrent safety documentation
- Better error handling strategies
- Performance optimization guidelines

### Changed
- **BREAKING**: Package name changed from `github.com/go-fork/di` to `go.fork.vn/di`
- All documentation updated to reflect new package name
- Installation instructions updated for new import path

### Documentation
- Added comprehensive Vietnamese technical documentation:
  - `docs/container.md`: Complete DI Container documentation
  - `docs/application.md`: Application interface documentation  
  - `docs/provider.md`: ServiceProvider interface documentation
  - `docs/deferred.md`: ServiceProviderDeferred interface documentation
  - `docs/loader.md`: ModuleLoaderContract documentation
  - `docs/README.md`: Complete system overview and architecture guide
- All documentation includes practical examples, best practices, and production patterns
- Added detailed BindingFunc documentation with usage patterns

### Migration
- Update import statements from `github.com/go-fork/di` to `go.fork.vn/di`
- Update mocks import from `github.com/go-fork/di/mocks` to `go.fork.vn/di/mocks`
- No API changes - seamless transition for existing code

## [0.0.5] - 2025-05-26

### Added
- Thêm method `Requires()` vào ServiceProvider interface để quản lý dependencies giữa các provider
- Thêm method `Providers()` vào ServiceProvider interface để liệt kê các service được đăng ký
- Thêm method `RegisterWithDependencies()` vào Application interface để tự động sắp xếp thứ tự đăng ký provider theo dependencies
- Cập nhật mocks tương ứng cho tất cả các interface changes
- Cải thiện dependency management với khả năng tự động phát hiện và sắp xếp thứ tự provider dependencies
- Nâng cao khả năng debug và kiểm tra service registration thông qua method `Providers()`
### Breaking Changes
- ServiceProvider interface giờ yêu cầu implement thêm methods `Requires()` và `Providers()`

## [0.0.4] - 2025-05-25
### Changed
- Thiết lập GitHub Actions cho CI/CD pipeline
- Thêm .goreleaser.yml cho tự động release

## [0.0.3] - 2025-05-24
### Added
- Thư mục mocks/ chứa các mock objects cho tất cả interface chính
- Tích hợp mockery để tạo mocks tự động
- Tài liệu về cách sử dụng mocks trong testing

## [0.0.2] - 2025-05-20
### Added
- Hỗ trợ Go 1.23.9
- API đầy đủ tài liệu cho Container, Application, ServiceProvider
- Tài liệu đầy đủ theo chuẩn godoc
- Test code coverage rates 100% cho toàn bộ package

### Changed
- Tối ưu hoá hiệu suất cho Container.Call()
- Cải thiện xử lý lỗi và báo cáo lỗi

## [0.0.1] - 2024-09-21

### Added
- DI Container với quản lý binding, singleton, instance, alias
- Interface ServiceProvider và ServiceProviderDeferred
- Interface Application để quản lý container và service provider
- Hệ thống deferred service provider
- Hỗ trợ module loading
- Tự động resolve dependency thông qua reflection

### Changed
- Cải thiện API dựa trên phản hồi từ người dùng

### Fixed
- Xử lý các lỗi edge case khi resolve dependencies
- Cải thiện đồng thời (concurrency) cho container
