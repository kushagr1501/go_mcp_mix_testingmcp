# Filesystem Analyzer

A command-line tool that scans directories and identifies unused files, zero-byte files, and other storage issues — with special support for OneDrive and cloud storage.

🔗 **Related Project (Non-MCP Version)**
If you're interested in a simpler, pure Go implementation **without MCP**, check out this repository:

👉 [https://github.com/yourusername/go-filesystem-native](https://github.com/kushagr1501/go-filesystem-)

🎥 **Demo Video**
A short demo of the tool in action is included in this repository under the `assets/` directory and is embedded below for quick reference.

> *(Video file is stored locally in this repo — see `assets/demo.mp4`)*

---

## Features

### Smart File Analysis

* **Unused File Detection** — Finds files that haven't been accessed in 60+ days
* **Zero-Byte File Detection** — Identifies empty or incomplete files
* **Detailed Explanations** — Every flagged file comes with clear reasons and evidence

### Cloud Storage Support

* **OneDrive Support** — Correctly handles OneDrive placeholder files and gets real file sizes
* **Windows API Integration** — Uses native Windows APIs for accurate file information
* **Smart Fallback** — Automatically tries common extensions (.pdf, .doc, .docx, .txt) if needed

### Technical Features

* **MCP-Based Architecture** — Built on Model Context Protocol for filesystem operations
* **Comprehensive Metadata** — Extracts file size, creation/modification/accessed dates, and MIME types
* **Robust Parsing** — Handles filenames with spaces and special characters

---

## Installation

### Prerequisites

* **Go 1.21+** — [https://golang.org/dl/](https://golang.org/dl/)
* **Node.js 18+** — Required for `mcp-filesystem-server`
* **Windows OS** — Currently Windows-only due to OneDrive API integration

### Step 1: Install mcp-filesystem-server

This project depends on `mcp-filesystem-server` for filesystem operations.

```bash
npm install -g mcp-filesystem-server
# or
npm install mcp-filesystem-server
```

Verify installation:

```bash
mcp-filesystem-server --help
```

### Step 2: Clone and Build This Project

```bash
git clone https://github.com/yourusername/go-filesystem.git
cd go-filesystem
go mod download
go build -o filesystem-analyzer.exe
```

---

## Usage

### Basic Usage

```bash
go run .
# or
.\\filesystem-analyzer.exe
```

When prompted, enter the directory path:

```
Enter directory path to analyze: C:\Users\YourName\Documents\Work
```

### Output Example

```
FILES COUNT: 5

[UNUSED] C:\Users\YourName\Documents\Work\old_report.pdf
  - Not accessed in last 60 days
  - Last accessed: 2024-11-15
  - Size: 1048576 bytes

[ZERO-BYTE] C:\Users\YourName\Documents\Work\temp_file.txt
  - File size is 0 bytes
  - Likely placeholder or incomplete file
```

---

## Implementation Details

### Architecture

```
main.go              # Entry point & orchestration
mcp_client.go        # JSON-RPC client for MCP communication
mcp_types.go         # MCP response/record definitions
filesystem.go        # FileInfo struct
listdirectory.go     # Directory listing via MCP
fileinfo.go          # File metadata extraction
rules.go             # Analysis rules
explanation.go       # Human-readable explanations
cloud.go             # Windows API (OneDrive handling)
```

### OneDrive Handling

* Uses Windows `FindFirstFile` API to get real file sizes
* Detects cloud placeholder attributes
* Skips cloud-only placeholders in zero-byte warnings
* Works with synced and cloud-only OneDrive files

---

## Analysis Rules

### Unused Files

* Not accessed in 60+ days
* Regular files only
* Valid access time required

### Zero-Byte Files

* Actual size is 0 bytes
* Not a OneDrive placeholder
* Regular files only

---

## Credits & Dependencies

### mcp-filesystem-server

GitHub: [https://github.com/mark3labs/mcp-filesystem-server](https://github.com/mark3labs/mcp-filesystem-server)
License: MIT

Provides filesystem access via Model Context Protocol (JSON-RPC).

---

## Future Enhancements

* Interactive TUI
* Configurable thresholds
* File type filters
* JSON / CSV export
* Recursive scans
* File cleanup actions
* macOS & Linux support

---

## License

MIT License

## Support

Found a bug or want a feature? Open an issue:
[https://github.com/yourusername/go-filesystem/issues](https://github.com/yourusername/go-filesystem/issues)
