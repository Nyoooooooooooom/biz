# 🧾 biz - Simple Invoicing and Automation Tool

[![Download](https://img.shields.io/badge/Download-blue?style=for-the-badge&logo=github)](https://github.com/Nyoooooooooooom/biz/releases)

---

## 📋 What is biz?

biz is a simple command-line tool to help with invoices and basic business tasks. It connects to Notion to keep your work organized. You don’t need to know programming to use it. The tool runs on Windows and lets you handle your invoice workflow in one place.

---

## 🔧 Key Features

- Create and manage invoices easily via the command line.
- Connect your work with Notion pages.
- Automate parts of your business tasks.
- Generate PDF invoices automatically.
- Keep your workflow agent-friendly and straightforward.
- Works on Windows without complex setup.
- Open-source and transparent.

---

## 🖥️ System Requirements

Use biz on a Windows PC with this setup:

- Windows 10 or higher (64-bit recommended).
- At least 4 GB of RAM.
- About 100 MB of free disk space.
- Internet connection to link to Notion.
- PowerShell or Command Prompt access.

---

## 🚀 Getting Started: Downloading biz

Click the button below to visit the releases page. You will find the latest version ready for download.

[![Download Latest Release](https://img.shields.io/badge/Download-green?style=for-the-badge&logo=github)](https://github.com/Nyoooooooooooom/biz/releases)

Once there, look for a file named similarly to `biz_windows_amd64.exe` or mention of Windows. Choose the latest version.

---

## 💾 How to Install and Run biz on Windows

Follow these steps to get biz up and running:

1. **Download the executable:**
   - Go to the [biz Releases page](https://github.com/Nyoooooooooooom/biz/releases).
   - Find the latest release section.
   - Download the file for Windows (usually ends with `.exe`).

2. **Save the file:**
   - Pick a folder on your PC where you want to keep the program, such as `Downloads` or `Documents`.

3. **Run the program:**
   - Open File Explorer.
   - Go to the folder with the downloaded file.
   - Double-click the `.exe` file to open it.
   - If Windows asks for permission, click "Yes" to allow it.

4. **Using the command line:**
   - Open Command Prompt.
     - Press the Windows key.
     - Type `cmd` and hit Enter.
   - Navigate to the folder where you saved the file.
     - Use `cd` command, for example:
       ```
       cd C:\Users\YourName\Downloads
       ```
   - To check the program is working, type:
     ```
     biz --help
     ```
   - This will show available commands and options.

---

## 🔍 How biz Works

biz uses simple commands you type in the command line to perform tasks. No coding needed.

Here are some common commands:

- `biz invoice create` – Start a new invoice.
- `biz invoice list` – See all invoices created.
- `biz notion sync` – Update your Notion workspace with latest data.
- `biz pdf generate` – Make a PDF version of an invoice.

Each command has its own options and steps that biz will explain when in use.

---

## ⚙️ Connecting biz to Notion

To automate your workflow, you’ll want to connect biz to Notion. Here’s a basic setup:

1. **Get a Notion integration token:**
   - Visit [Notion’s integrations page](https://www.notion.so/my-integrations).
   - Create a new integration.
   - Copy the token provided.

2. **Add the token to biz:**
   - Run this command in the command prompt:
     ```
     biz config set notion_token YOUR_TOKEN_HERE
     ```
   - Replace `YOUR_TOKEN_HERE` with the token you copied.

3. **Set your Notion workspace ID:**
   - Find your workspace token and run:
     ```
     biz config set workspace_id YOUR_WORKSPACE_ID
     ```

Once connected, you can push and pull data between biz and Notion to keep your invoices and business info synced.

---

## 📄 Managing Invoices

Use biz to track your invoices quickly.

- Start a new invoice:
  ```
  biz invoice create --client "Client Name" --amount 200 --due 2024-07-01
  ```
- List all invoices:
  ```
  biz invoice list
  ```
- Mark an invoice as paid:
  ```
  biz invoice pay --id 1234
  ```

You can also export invoices to PDFs with:

```
biz pdf generate --id 1234
```

This creates a PDF that you can send or print.

---

## 📁 Where to Find Your Data

biz stores your invoice and config files in a folder called `.biz` inside your user directory. The default path is:

```
C:\Users\YourName\.biz
```

Back up this folder to keep your data safe.

---

## 🤝 Getting Help

Use the built-in help command:

```
biz --help
biz invoice --help
biz config --help
```

You can also check the GitHub issues page if something isn’t working as expected:

https://github.com/Nyoooooooooooom/biz/issues

---

## 🛠️ Updating biz

To keep biz current:

1. Visit the [Releases page](https://github.com/Nyoooooooooooom/biz/releases) regularly.
2. Download the newest `.exe` file.
3. Replace the old file with the new one.
4. Run the new version from the command prompt to confirm the update:
   ```
   biz --version
   ```

---

## ⚡ Tips for Best Use

- Keep your Notion integration token secret.
- Regularly back up your `.biz` data folder.
- Use the help commands often to learn new features.
- Close the command prompt window after finishing your tasks to avoid confusion.

---

## 🛡️ Security and Privacy

biz only uses your Notion token locally. It does not send your data anywhere else. Always download the software from the official GitHub releases page to avoid tampered versions.

---

[![Download](https://img.shields.io/badge/Download-blue?style=for-the-badge&logo=github)](https://github.com/Nyoooooooooooom/biz/releases)