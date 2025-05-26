# Changelog

Tất cả các thay đổi đáng chú ý của dự án go-fork/di sẽ được ghi lại trong file này.

Định dạng dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân theo [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Thêm method `Requires()` vào ServiceProvider interface để quản lý dependencies giữa các provider
- Thêm method `Providers()` vào ServiceProvider interface để liệt kê các service được đăng ký
- Thêm method `RegisterWithDependencies()` vào Application interface để tự động sắp xếp thứ tự đăng ký provider theo dependencies
- Cập nhật mocks tương ứng cho tất cả các interface changes

### Changed
- Cải thiện dependency management với khả năng tự động phát hiện và sắp xếp thứ tự provider dependencies
- Nâng cao khả năng debug và kiểm tra service registration thông qua method `Providers()`

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
