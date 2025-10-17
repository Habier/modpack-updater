#!/bin/bash
set -euo pipefail

# Check that the script is run from the project root
if [ ! -d "modpack-updater_deb" ]; then
  echo "Error: Run this script from the project root."
  exit 1
fi

# Ensure dpkg-deb is available
if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "Error: dpkg-deb is required to build the package."
  exit 1
fi

# Allow version to be set as an argument or environment variable
VERSION="${1:-${VERSION:-1.0.0}}"
PKG_DIR="modpack-updater_deb"
PKG_NAME="modpack-updater_${VERSION}_all.deb"

# Update control file Version field if present
CONTROL_FILE="$PKG_DIR/DEBIAN/control"
if [ -f "$CONTROL_FILE" ]; then
  if grep -qE '^Version:' "$CONTROL_FILE"; then
    sed -i "s/^Version:.*/Version: $VERSION/" "$CONTROL_FILE"
    echo "Updated $CONTROL_FILE with Version: $VERSION"
  fi
fi

# Ensure the package contains the command without .sh
BIN_DIR="$PKG_DIR/usr/local/bin"
PLAIN_BIN="$BIN_DIR/modpack-updater"

# Ensure execution permissions for the final binary
chmod 755 "$PLAIN_BIN"

echo "Set executable permission on $PLAIN_BIN"

# Handle autocomplete file if present
COMPLETION_FILE="$PKG_DIR/etc/bash_completion.d/modpack-updater"
if [ ! -f "$COMPLETION_FILE" ]; then
  echo "Warning: Autocomplete file $COMPLETION_FILE does not exist. It will not be included in the package."
else
  # Ensure readable permissions
  chmod 644 "$COMPLETION_FILE"
  echo "Included autocomplete: $COMPLETION_FILE"
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

# Build the DEB package
echo "Building DEB package: $PKG_NAME"
if dpkg-deb --build "$PKG_DIR" "$PKG_NAME"; then
  echo "DEB package created: $PKG_NAME"
else
  echo "Error creating DEB package"
  exit 1
fi

exit 0
