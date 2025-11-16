# Modpack Updater

A fast and reliable tool for updating Minecraft modpack folders from ZIP
archives, written in Go.

It will replace the following folders:

-   `config`
-   `defaultconfigs`
-   `modernfix`
-   `mods`
-   `schematics`

## 🚀 Features

-   Removes old modpack folders (config, mods, etc.)
-   Extracts updated folders from a ZIP archive
-   Cross-platform support (Windows, Linux, macOS)
-   Clean and user-friendly output
-   Bash completion support

## 📦 Requirements

-   Go 1.23 or higher (for building from source)

## 🧰 Usage

``` bash
modpack-updater update <destination_directory> <zip_file>
```

## 📥 Installation

### From Debian Package (recommended)

1.  Download the latest `.deb` package from the releases page.

2.  Install using `dpkg`:

    ``` bash
    sudo dpkg -i modpack-updater_*.deb
    sudo apt-get install -f
    ```

### From Source

1.  Clone the repository:

    ``` bash
    git clone https://github.com/yourusername/modpack-updater.git
    cd modpack-updater
    ```

2.  Build the application:

    ``` bash
    go build -o modpack-updater .
    ```

3.  (Optional) Install to your PATH:

    ``` bash
    sudo cp modpack-updater /usr/local/bin/
    ```

## 🏗️ Building the Debian Package

To build a Debian package:

``` bash
./build-deb.sh [version]
```

If no version is specified, it will default to **1.0.0**.

## 📄 License

This project is licensed under the MIT License. See the
[LICENSE](LICENSE) file for details.
