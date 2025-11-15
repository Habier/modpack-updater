# Modpack Updater (Go Version)

A fast and reliable tool for updating Minecraft modpack folders from ZIP archives, written in Go.

## Features

- Removes old modpack folders (config, mods, etc.)
- Extracts updated folders from a ZIP archive
- Cross-platform support (Windows, Linux, macOS)
- Clean and user-friendly output
- Bash completion support

## Requirements

- Go 1.23 or higher (for building from source)

## Usage

```bash
modpack-updater update <destination_directory> <zip_file>
```

## Installation

### From Debian Package (recommended)

1. Download the latest `.deb` package from the releases page
2. Install using `dpkg`:
   ```bash
   sudo dpkg -i modpack-updater_*.deb
   sudo apt-get install -f  # Install any missing dependencies
   ```

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/modpack-updater.git
   cd modpack-updater
   ```

2. Build the application:
   ```bash
   go build -o modpack-updater .
   ```

3. Install to your PATH (optional):
   ```bash
   sudo cp modpack-updater /usr/local/bin/
   ```

## Building the Debian Package

To build a Debian package:

```bash
# From the project root
./build-deb.sh [version]
```

If no version is specified, it will use 1.0.0 by default.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
