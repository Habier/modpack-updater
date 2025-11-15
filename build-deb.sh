#!/bin/bash
set -euo pipefail

[ -d "modpack-updater_deb" ] || { echo "Error: Run this script from the project root."; exit 1; }

# Check required commands
for cmd in dpkg-deb go; do
  command -v "$cmd" >/dev/null 2>&1 || { echo "Error: $cmd is required"; exit 1; }
done

VERSION="${1:-${VERSION:-1.0.0}}"
PKG_DIR="modpack-updater_deb"
PKG_NAME="modpack-updater_${VERSION}_all.deb"
BIN_DIR="$PKG_DIR/usr/local/bin"
PLAIN_BIN="$BIN_DIR/modpack-updater"

# Update control file
CONTROL_FILE="$PKG_DIR/DEBIAN/control"
if [ -f "$CONTROL_FILE" ]; then
  grep -qE '^Version:' "$CONTROL_FILE" && sed -i "s/^Version:.*/Version: $VERSION/" "$CONTROL_FILE"
fi

go mod tidy 2>/dev/null
GOOS=linux GOARCH=amd64 go build -o "$PLAIN_BIN" .

# Set executable permissions
chmod 755 "$PLAIN_BIN"
echo "✅ Built binary: $PLAIN_BIN"

# Handle autocomplete file if present
COMPLETION_FILE="$PKG_DIR/etc/bash_completion.d/modpack-updater"
if [ ! -f "$COMPLETION_FILE" ]; then
  echo "⚠ Warning: Autocomplete file $COMPLETION_FILE does not exist. It will not be included in the package."
else
  # Ensure readable permissions
  chmod 644 "$COMPLETION_FILE"
  echo "✅ Included autocomplete: $COMPLETION_FILE"
fi

# Handle manpage if present
MANPAGE_PLAIN="$PKG_DIR/usr/share/man/man1/modpack-updater.1"
MANPAGE_GZ="$PKG_DIR/usr/share/man/man1/modpack-updater.1.gz"
if [ -f "$MANPAGE_PLAIN" ]; then
  # Ensure correct permissions and gzip the manpage (replace plain with .gz)
  chmod 644 "$MANPAGE_PLAIN"
  gzip -n -f -9 "$MANPAGE_PLAIN"
  echo "Gzipped manpage: $MANPAGE_GZ"
else
  if [ -f "$MANPAGE_GZ" ]; then
    chmod 644 "$MANPAGE_GZ"
    echo "Manpage (gz) exists: $MANPAGE_GZ"
  else
    echo "Warning: Manpage not found at $MANPAGE_PLAIN (or .gz). No manpage will be included."
  fi
fi

# Make maintainer scripts executable if present
MAINT_SCRIPTS=("postinst" "prerm" "postrm" "preinst")
for s in "${MAINT_SCRIPTS[@]}"; do
  SCRIPT_PATH="$PKG_DIR/DEBIAN/$s"
  if [ -f "$SCRIPT_PATH" ]; then
    chmod 755 "$SCRIPT_PATH"
    echo "Set executable permission on maintainer script: $SCRIPT_PATH"
  fi
done

# Ensure DEBIAN directory and control file have acceptable permissions for dpkg-deb
if [ -d "$PKG_DIR/DEBIAN" ]; then
  chmod 0755 "$PKG_DIR/DEBIAN"
  if [ -f "$CONTROL_FILE" ]; then
    chmod 0644 "$CONTROL_FILE"
    echo "Set permissions on $PKG_DIR/DEBIAN (0755) and $CONTROL_FILE (0644)"
  fi
fi

# Remove previous package if it exists
if [ -f "$PKG_NAME" ]; then
  echo "Removing previous package: $PKG_NAME"
  rm -f "$PKG_NAME"
fi

# Build the package
echo "📦 Building Debian package $PKG_NAME..."
if dpkg-deb --build "$PKG_DIR" "$PKG_NAME" >/dev/null; then
  # Get the package size
  PKG_SIZE=$(du -h "$PKG_NAME" | cut -f1)
  
  echo ""
  echo "✨ Successfully built $PKG_NAME (${PKG_SIZE}B)"
  
  # Show package info
  echo -e "\n📋 Package information:"
  dpkg-deb --info "$PKG_NAME" | grep -E 'Package:|Version:|Architecture:|Installed-Size:'
  
  # Optional: Check the package with lintian if available
  if command -v lintian >/dev/null 2>&1; then
    echo -e "\n🔍 Running package linting..."
    lintian "$PKG_NAME" || true
  fi
  
  echo -e "\n✅ Build complete! Package: $(pwd)/$PKG_NAME"
  exit 0
else
  echo "❌ Failed to build package" >&2
  exit 1
fi

exit 0
