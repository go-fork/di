# Release Notes - v0.1.3

**Release Date**: June 4, 2025  
**Type**: Documentation & Tooling Enhancement

## 🎯 Key Improvements

### Repository Structure Reorganization
- **Clean Root Directory**: Moved all release documentation to `releases/` directory
- **Structured Documentation**: Clear separation between development and released docs
- **Professional Layout**: Following industry best practices for Go projects

### Release Management Automation  
- **Automated Scripts**: Added `scripts/archive_release.sh` and `scripts/create_release_templates.sh`
- **Streamlined Workflow**: Simple process for creating new releases
- **Consistent Naming**: Standardized documentation file naming across versions

### Documentation Enhancements
- **Clear Navigation**: Easy access to both current and historical documentation
- **Comprehensive Guides**: Detailed workflow documentation in `releases/README.md`
- **Developer Experience**: Simplified development workflow

## 🛠️ New Features

### Release Automation Scripts
```bash
# Archive current release and setup new cycle
./scripts/archive_release.sh v1.2.3

# Create fresh templates for development
./scripts/create_release_templates.sh
```

### Structured Documentation
```
releases/
├── next/           # Current development docs
├── v0.1.0/        # Historical releases  
├── v0.1.1/
└── v0.1.2/
```

## 💥 Breaking Changes

- **Documentation Paths**: 
  - `MIGRATION.md` → `releases/next/MIGRATION.md`
  - `RELEASE_NOTES.md` → `releases/next/RELEASE_NOTES.md`
  - `RELEASE_SUMMARY.md` → `releases/next/RELEASE_SUMMARY.md`

## 🔧 Improvements

- **Cleaner Repository**: Root directory only contains source code and core docs
- **Better Organization**: Historical documentation properly archived
- **Automation Ready**: Scripts handle tedious release management tasks
- **Professional Structure**: Following Go community standards

## 📦 Installation

```bash
go get go.fork.vn/di@v0.1.3
```

## 📚 Documentation

- [Migration Guide](MIGRATION.md)
- [Release Process](RELEASE_SUMMARY.md)  
- [Container Usage](../../docs/container.md)
- [Service Providers](../../docs/provider.md)

## 🔗 Historical Release Notes

- [v0.1.2](../v0.1.2/RELEASE_NOTES_v0.1.2.md)
- [v0.1.0](../v0.1.0/RELEASE_NOTES_v0.1.0.md)
