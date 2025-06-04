# Migration Guide - v0.1.3

## From Previous Versions (v0.1.0, v0.1.1, v0.1.2)

### Breaking Changes
- **Documentation Structure**: Migration, Release Notes và Release Summary đã được di chuyển từ root directory vào `releases/next/`
- **Symlinks Removed**: Không còn symlinks ở root directory

### New Features  
- **Release Management Automation**: Scripts tự động cho việc quản lý release
- **Structured Documentation**: Tổ chức tài liệu theo phiên bản rõ ràng
- **Clean Repository Structure**: Root directory chỉ chứa source code

## Migration Steps

### 1. Update Documentation Links
Nếu bạn đang link đến documentation:

**Before:**
```markdown
[Migration Guide](MIGRATION.md)
[Release Notes](RELEASE_NOTES.md)
```

**After:**
```markdown
[Migration Guide](releases/next/MIGRATION.md)  
[Release Notes](releases/next/RELEASE_NOTES.md)
```

### 2. Update Dependencies
```bash
go get go.fork.vn/di@v0.1.3
```

### 3. No Code Changes Required
- API không thay đổi
- Existing code sẽ hoạt động bình thường
- Chỉ có cấu trúc documentation thay đổi

## Historical Migrations

For migration guides from older versions, see:
- [v0.1.2](../v0.1.2/MIGRATION_v0.1.2.md)
- [v0.1.1](../v0.1.1/MIGRATION_v0.1.1.md) 
- [v0.1.0](../v0.1.0/MIGRATION_v0.1.0.md)
