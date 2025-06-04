# Release Management Scripts

This directory contains automation scripts for managing releases and documentation.

## Scripts

### `create_release_templates.sh`
Creates fresh templates in `releases/next/` for new development cycle.

**Usage:**
```bash
./scripts/create_release_templates.sh
```

**What it does:**
- Creates template `MIGRATION.md` with standard structure
- Creates template `RELEASE_NOTES.md` with standard sections
- Creates template `RELEASE_SUMMARY.md` with release checklist

### `archive_release.sh`
Archives current release and sets up new development cycle.

**Usage:**
```bash
./scripts/archive_release.sh vX.X.X
```

**Example:**
```bash
./scripts/archive_release.sh v1.2.3
```

**What it does:**
1. Creates new version directory `releases/vX.X.X/`
2. Moves and renames files from `releases/next/` to version directory
3. Creates fresh templates for next development cycle
4. Updates `releases/README.md` structure

## Workflow Integration

### During Development
Work directly with files in `releases/next/`:
- Edit `releases/next/MIGRATION.md`
- Edit `releases/next/RELEASE_NOTES.md`  
- Track progress in `releases/next/RELEASE_SUMMARY.md`

### When Ready to Release
```bash
# Archive current work and create new cycle
./scripts/archive_release.sh v1.2.3

# Update CHANGELOG.md manually
# Create git tag and push
# Create GitHub release
```

## Benefits

- **Automated**: No manual file moving or renaming
- **Consistent**: Standard templates and naming conventions
- **Safe**: Validates input and checks for existing versions
- **Complete**: Updates all related documentation automatically

## Requirements

- Bash shell
- `sed` command (for updating README)
- Proper directory structure (`releases/next/` must exist)

These scripts maintain the clean separation between development and released documentation while automating the tedious parts of the release process.
