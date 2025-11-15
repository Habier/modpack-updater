#!/bin/bash

if [ -z "$1" ]; then
    echo "Error: Version number required (e.g., 1.0.0)"
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="$HOME/modpack-updater"
DESTINATION_DIR="$SCRIPT_DIR"

# Check required commands
for cmd in cp sudo mv; do
    command -v "$cmd" &> /dev/null || { echo "Error: '$cmd' command not found"; exit 1; }
done

[ -d "$TARGET_DIR" ] && { sudo rm -rf "$TARGET_DIR" || { echo "Failed to remove $TARGET_DIR"; exit 1; }; }
cp -r "$SCRIPT_DIR" "$TARGET_DIR" || { echo "Failed to copy files"; exit 1; }

# Build .deb package
cd "$TARGET_DIR" || { echo "Failed to cd to $TARGET_DIR"; exit 1; }
chmod +x build-deb.sh
./build-deb.sh "$VERSION" || { echo "Failed to build .deb package"; exit 1; }

# Move generated .deb file
DEB_FILE="modpack-updater_${VERSION}_all.deb"
if [ -f "$DEB_FILE" ]; then
    mv "$DEB_FILE" "$DESTINATION_DIR/" || { echo "Failed to move .deb file"; exit 1; }
    echo "Package created: $DESTINATION_DIR/$DEB_FILE"
else
    echo "Error: Generated .deb file not found"
    exit 1
fi
