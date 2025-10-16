# 🧩 modpack-updater.sh

A simple Bash script that removes a predefined set of folders from a target directory, then restores them by extracting those same folders from a ZIP file.

---

## 🚀 Features

- Takes **two arguments**:
    1. Target directory path
    2. ZIP file path
- Deletes folders listed in the `folders` array.
- Extracts those same folders from the ZIP file into the target directory.

---

## 📦 System-wide Installation (Linux)

To make the script available for **all users** on your Linux system:

1. Copy the script to `/usr/local/bin`:
   ```bash
   sudo cp modpack-updater.sh /usr/local/bin/
   ```
2. Make it executable:
   ```bash
   sudo chmod +x /usr/local/bin/modpack-updater.sh
   ```
3. Verify installation:
   ```bash
   which modpack-updater.sh
   ```
   If it returns `/usr/local/bin/modpack-updater.sh`, it’s successfully installed.

---

## 🧠 Usage

Run it from any user and directory:

```bash
modpack-updater.sh <target_directory> <zip_file>
```

### 📘 Example

```bash
modpack-updater.sh /home/javier/minecraft /home/javier/backups/mods_backup.zip
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

---

## ⚡ Optional: Enable Bash Autocompletion

To make the script easier to use (auto-completing file and folder names):

1. Create a completion script file:
   ```bash
   sudo nano /etc/bash_completion.d/modpack-updater
   ```
2. Paste the following content:
   ```bash
   _modpack-updater_completions() {
       local cur prev
       COMPREPLY=()
       cur="${COMP_WORDS[COMP_CWORD]}"
       prev="${COMP_WORDS[COMP_CWORD-1]}"
       # Autocomplete directories for the first argument
       if [[ $COMP_CWORD -eq 1 ]]; then
           COMPREPLY=( $(compgen -d -- "$cur") )
       # Autocomplete .zip files for the second argument
       elif [[ $COMP_CWORD -eq 2 ]]; then
           COMPREPLY=( $(compgen -f -- "$cur" | grep -E '\.zip$') )
       fi
       return 0
   }
   complete -F _modpack-updater_completions modpack-updater.sh
   ```
3. Save and exit (Ctrl+O, Enter, Ctrl+X).
4. Reload completions:
   ```bash
   source /etc/bash_completion.d/modpack-updater
   ```

Now, when you type:
```bash
modpack-updater.sh [TAB]
```
you’ll get folder suggestions, and after that, ZIP file suggestions.

---

## 🧹 Uninstallation

To remove the script and its autocompletion:

```bash
sudo rm /usr/local/bin/modpack-updater.sh
sudo rm /etc/bash_completion.d/modpack-updater
```

---

## ✅ Example Workflow

```bash
# Backup existing Minecraft folders
zip -r /home/javier/backups/mods_backup.zip config defaultconfigs modernfix mods schematics

# Update with new files
modpack-updater.sh /home/javier/minecraft /home/javier/backups/mods_backup.zip
```
