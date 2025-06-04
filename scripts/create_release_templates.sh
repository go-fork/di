#!/bin/bash

# create_release_templates.sh
# Creates fresh templates in releases/next/ for new development cycle

set -e

NEXT_DIR="releases/next"

echo "🚀 Creating fresh release templates for next development cycle..."

# Create MIGRATION.md template
cat > "$NEXT_DIR/MIGRATION.md" << 'EOF'
# Migration Guide

This document provides migration instructions for the upcoming version.

## From Previous Versions

### Breaking Changes
- [List any breaking changes here]

### Deprecations  
- [List any deprecated features here]

### New Features
- [List new features that may affect existing code]

## Migration Steps

### 1. Update Dependencies
```bash
go get github.com/go-fork/di@latest
```

### 2. Code Changes
[Describe specific code changes needed]

### 3. Testing
[Describe testing considerations]

## Historical Migrations

For migration guides from older versions, see the `releases/` directory.
EOF

# Create RELEASE_NOTES.md template
cat > "$NEXT_DIR/RELEASE_NOTES.md" << 'EOF'
# Release Notes

## Upcoming Version

### New Features
- [List new features here]

### Improvements
- [List improvements here]

### Bug Fixes
- [List bug fixes here]

### Breaking Changes
- [List breaking changes here]

## Installation

```bash
go get github.com/go-fork/di@latest
```

## Documentation

- [Getting Started](docs/index.md)
- [Container Usage](docs/container.md)
- [Service Providers](docs/provider.md)
- [Migration Guide](MIGRATION.md)

## Historical Release Notes

For release notes from previous versions, see the `releases/` directory.
EOF

# Create RELEASE_SUMMARY.md template
cat > "$NEXT_DIR/RELEASE_SUMMARY.md" << 'EOF'
# Release Summary

This document tracks the completion status of the upcoming release.

## 📋 Release Status: 🔄 IN PROGRESS

**Package Name**: `go.fork.vn/di`  
**Target Version**: [To be determined]  
**Type**: [To be determined]  
**Target Release Date**: [To be determined]

## 🎯 Release Checklist

### 1. **Code Changes** 
- [ ] All planned features implemented
- [ ] Breaking changes documented
- [ ] Tests updated and passing
- [ ] Code review completed

### 2. **Documentation**
- [ ] API documentation updated
- [ ] Migration guide updated
- [ ] Release notes prepared
- [ ] Examples updated

### 3. **Quality Assurance**
- [ ] All tests passing (`go test ./...`)
- [ ] Build successful (`go build ./...`)
- [ ] Integration tests completed
- [ ] Performance benchmarks

### 4. **Version Control**
- [ ] All changes committed
- [ ] Version tag created
- [ ] Branch pushed to remote
- [ ] Tag pushed to remote

### 5. **Release Process**
- [ ] GitHub release created
- [ ] Release notes published  
- [ ] Package available for download
- [ ] Documentation deployed

### 6. **Post-Release**
- [ ] CHANGELOG.md updated
- [ ] Community notifications sent
- [ ] Support documentation updated
- [ ] Archive release docs

## 📊 Progress Tracking

| Task Category | Progress | Status |
|---------------|----------|--------|
| Code Changes | 0/4 | ⏳ Pending |
| Documentation | 0/4 | ⏳ Pending |
| Quality Assurance | 0/4 | ⏳ Pending |
| Version Control | 0/4 | ⏳ Pending |
| Release Process | 0/4 | ⏳ Pending |
| Post-Release | 0/4 | ⏳ Pending |

**Overall Progress**: 0% (0/24 tasks completed)

## 🔗 Related Documents

- [Migration Guide](MIGRATION.md)
- [Release Notes](RELEASE_NOTES.md)
- [Changelog](CHANGELOG.md)

---

**Note**: This document should be updated throughout the release process.
EOF

echo "✅ Templates created successfully in $NEXT_DIR/"
echo "📁 Files created:"
echo "   - $NEXT_DIR/MIGRATION.md"
echo "   - $NEXT_DIR/RELEASE_NOTES.md" 
echo "   - $NEXT_DIR/RELEASE_SUMMARY.md"
