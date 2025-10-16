#!/bin/bash
set -e

# Check that the script is run from the project root
if [ ! -d "modpack-updater_deb" ]; then
  echo "Error: Run this script from the project root."
  exit 1
fi

# Allow version to be set as an argument or environment variable
VERSION="${1:-${VERSION:-1.0.0}}"
PKG_DIR="modpack-updater_deb"
PKG_NAME="modpack-updater_${VERSION}_all.deb"

# Check that the main file exists
if [ ! -f "$PKG_DIR/usr/local/bin/modpack-updater.sh" ]; then
  echo "Error: File $PKG_DIR/usr/local/bin/modpack-updater.sh does not exist."
  exit 1
fi

# Ensure execution permissions
chmod 755 "$PKG_DIR/usr/local/bin/modpack-updater.sh"

# Check that the autocomplete file exists
if [ ! -f "$PKG_DIR/etc/bash_completion.d/modpack-updater" ]; then
  echo "Warning: Autocomplete file $PKG_DIR/etc/bash_completion.d/modpack-updater does not exist."
else
  # Ensure execution permissions (readable for all)
  chmod 644 "$PKG_DIR/etc/bash_completion.d/modpack-updater"
fi

# Remove previous package if it exists
if [ -f "$PKG_NAME" ]; then
  echo "Removing previous package: $PKG_NAME"
  rm -f "$PKG_NAME"
fi

# Build the DEB package
if dpkg-deb --build "$PKG_DIR" "$PKG_NAME"; then
  echo "DEB package created: $PKG_NAME"
else
  echo "Error creating DEB package"
  exit 1
fi
