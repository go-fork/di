#!/bin/bash

# archive_release.sh vX.X.X
# A# Step 2: Move and rename files from next/ to versioned directory
echo "📋 Archiving release documents..."
mv "$NEXT_DIR/MIGRATION.md" "$VERSION_DIR/MIGRATION_$VERSION.md"
mv "$NEXT_DIR/RELEASE_NOTES.md" "$VERSION_DIR/RELEASE_NOTES_$VERSION.md"
mv "$NEXT_DIR/RELEASE_SUMMARY.md" "$VERSION_DIR/RELEASE_SUMMARY.md"

# Step 3: Create fresh templates for next version
echo "🆕 Creating fresh templates for next development cycle..."
./scripts/create_release_templates.sh

# Step 4: Update releases/README.md to include new versionelease and creates new development cycle

set -e

VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "❌ Error: Version number required"
    echo "Usage: $0 vX.X.X"
    echo "Example: $0 v1.2.3"
    exit 1
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "❌ Error: Invalid version format"
    echo "Version must be in format vX.X.X (e.g., v1.2.3)"
    exit 1
fi

RELEASES_DIR="releases"
VERSION_DIR="$RELEASES_DIR/$VERSION"
NEXT_DIR="$RELEASES_DIR/next"

echo "🗂️  Archiving release $VERSION..."

# Check if version directory already exists
if [ -d "$VERSION_DIR" ]; then
    echo "❌ Error: Version directory $VERSION_DIR already exists"
    exit 1
fi

# Check if next directory exists
if [ ! -d "$NEXT_DIR" ]; then
    echo "❌ Error: $NEXT_DIR directory not found"
    exit 1
fi

# Step 1: Create version directory
echo "📁 Creating version directory: $VERSION_DIR"
mkdir -p "$VERSION_DIR"

# Step 2: Move and rename files from next/ to versioned directory
echo "📋 Archiving release documents..."
mv "$NEXT_DIR/MIGRATION.md" "$VERSION_DIR/MIGRATION_$VERSION.md"
mv "$NEXT_DIR/RELEASE_NOTES.md" "$VERSION_DIR/RELEASE_NOTES_$VERSION.md"
mv "$NEXT_DIR/RELEASE_SUMMARY.md" "$VERSION_DIR/RELEASE_SUMMARY.md"

# Step 4: Create fresh templates for next version
echo "🆕 Creating fresh templates for next development cycle..."
./scripts/create_release_templates.sh

# Step 5: Recreate symlinks
echo "🔗 Creating new symlinks..."
ln -s "$NEXT_DIR/MIGRATION.md" MIGRATION.md
ln -s "$NEXT_DIR/RELEASE_NOTES.md" RELEASE_NOTES.md  
ln -s "$NEXT_DIR/RELEASE_SUMMARY.md" RELEASE_SUMMARY.md

# Step 6: Update releases/README.md to include new version
echo "📝 Updating releases/README.md..."
# Add the new version to the structure in README
sed -i '' "/├── v0\.1\.2\//a\\
├── $VERSION/\\
│   ├── MIGRATION_$VERSION.md\\
│   ├── RELEASE_NOTES_$VERSION.md\\
│   └── RELEASE_SUMMARY.md\\
" "$RELEASES_DIR/README.md"

echo ""
echo "✅ Release $VERSION archived successfully!"
echo ""
echo "📁 Created:"
echo "   - $VERSION_DIR/MIGRATION_$VERSION.md"
echo "   - $VERSION_DIR/RELEASE_NOTES_$VERSION.md"
echo "   - $VERSION_DIR/RELEASE_SUMMARY.md"
echo ""
echo "📋 Next steps:"
echo "   1. Update CHANGELOG.md with $VERSION details"
echo "   2. Create git tag: git tag -a $VERSION -m 'Release $VERSION'"
echo "   3. Push tag: git push origin $VERSION"
echo "   4. Create GitHub release"
echo ""
echo "🚀 Ready for next development cycle!"
