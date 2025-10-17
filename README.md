# 🧩 modpack-updater

A simple Bash script that removes a predefined set of folders from a target directory, then restores them by extracting those same folders from a ZIP file.

---

## 🚀 Features

- Takes **two arguments**:
    1. Target directory path
    2. ZIP file path
- Deletes folders listed in the `folders` array.
- Extracts those same folders from the ZIP file into the target directory.

---

## 📦 Install via .deb Package (Recommended)

You can build and install a .deb package for easy installation on Debian/Ubuntu systems.

### Install the package

```bash
sudo dpkg -i modpack-updater_1.0.0_all.deb
```

If you see dependency errors, fix them with:

```bash
sudo apt-get install -f
```

After installation, you can use `modpack-updater` from anywhere as described above.

---

## 🧠 Usage

Run it from any user and directory:

```bash
modpack-updater <target_directory> <zip_file>
```

### 📘 Example

```bash
modpack-updater /home/habier/minecraft /home/habier/new_modpack.zip
```

---

## 🧰 Requirements

- `bash` – included by default on most Linux distributions
- `unzip` – required to extract files

If `unzip` is not installed:

```bash
sudo apt install unzip
```
(Use your distro’s package manager if not on Debian/Ubuntu.)

---

## 🧾 Internal Folder List

Inside the script, these are the folders that will be deleted and restored:

```bash
folders=("config" "defaultconfigs" "modernfix" "mods" "schematics")
```
You can edit this list in the script to include or remove folders.
