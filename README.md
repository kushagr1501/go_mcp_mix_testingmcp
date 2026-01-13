# Filesystem Analyzer

A command-line tool that scans directories and identifies unused files, zero-byte files, and other storage issues - with special support for OneDrive and cloud storage.

## Features

### Smart File Analysis
- **Unused File Detection** - Finds files that haven't been accessed in 60+ days
- **Zero-Byte File Detection** - Identifies empty or incomplete files
- **Detailed Explanations** - Every flagged file comes with clear reasons and evidence

### Cloud Storage Support
- **OneDrive Support** - Correctly handles OneDrive placeholder files and gets real file sizes
- **Windows API Integration** - Uses native Windows APIs for accurate file information
- **Smart Fallback** - Automatically tries common extensions (.pdf, .doc, .docx, .txt) if needed

### Technical Features
- **MCP-Based Architecture** - Built on Model Context Protocol for filesystem operations
- **Comprehensive Metadata** - Extracts file size, creation/modification/accessed dates, and MIME types
- **Robust Parsing** - Handles filenames with spaces and special characters

## Installation

### Prerequisites
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Node.js 18+** - Required for mcp-filesystem-server
- **Windows OS** - Currently Windows-only due to OneDrive API integration

### Step 1: Install mcp-filesystem-server
This project depends on [mcp-filesystem-server](https://github.com/mark3labs/mcp-filesystem-server) for filesystem operations.

```bash
# Install globally via npm
npm install -g mcp-filesystem-server

# Or install locally in your project
npm install mcp-filesystem-server
```

Verify installation:
```bash
mcp-filesystem-server --help
```

### Step 2: Clone and Build This Project
```bash
# Clone repository
git clone https://github.com/yourusername/go-filesystem.git
cd go-filesystem

# Install dependencies
go mod download

# Build project
go build -o filesystem-analyzer.exe
```

## Usage

### Basic Usage
```bash
# Run analyzer
go run .

# Or use built executable
.\filesystem-analyzer.exe
```

When prompted, enter the directory path you want to analyze:
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
Reason: File is empty (0 bytes)
 - File size is 0 bytes
 - Likely placeholder or incomplete file
```

## Implementation Details

### Architecture
The project is built with a modular architecture:

```
main.go              # Entry point & orchestration
mcp_client.go        # JSON-RPC client for MCP communication
mcp_types.go        # MCP response/record type definitions
filesystem.go       # FileInfo struct for file metadata
listdirectoy.go     # Directory listing via MCP
fileinfo.go         # File metadata extraction
rules.go            # Analysis rules (unused, zero-byte)
explanation.go      # Explanation struct for findings
cloud.go           # Windows API for OneDrive support
```

### How It Works

1. **User Input** - Prompts for directory path via stdin
2. **MCP Connection** - Connects to `mcp-filesystem-server` via JSON-RPC
3. **Directory Scan** - Lists all files using `list_directory` tool
4. **File Analysis** - For each file:
   - Calls `get_file_info` to get metadata
   - Uses Windows API to get real file size (handles OneDrive)
   - Applies analysis rules (unused, zero-byte)
   - Generates explanations with evidence
5. **Output** - Prints findings to console

### OneDrive Handling
OneDrive creates placeholder files that appear as 0 bytes until downloaded. This is common even for files that exist locally but are synced with OneDrive. This project:
- Uses Windows `FindFirstFile` API to get actual file sizes
- Detects cloud placeholder attributes
- Skips placeholder files from zero-byte warnings
- Correctly identifies real files from cloud storage
- Works with both locally synced OneDrive files and cloud-only placeholders

## Analysis Rules

### Unused Files
Files are flagged as unused if:
- Not accessed in 60+ days (configurable)
- Is a regular file (not a directory)
- Has valid access time information

### Zero-Byte Files
Files are flagged as zero-byte if:
- Actual file size is 0 bytes (via Windows API)
- Not a OneDrive cloud placeholder
- Is a regular file (not a directory)

## Credits and Dependencies

### mcp-filesystem-server
This project uses [mcp-filesystem-server](https://github.com/mark3labs/mcp-filesystem-server) by mark3labs - an open-source MCP server that provides filesystem operations.

**What is MCP?**
The Model Context Protocol is an open standard that enables AI assistants to interact with external systems (like filesystems) through a consistent JSON-RPC interface.

**License:** MIT
**GitHub:** https://github.com/mark3labs/mcp-filesystem-server

## Future Enhancements

- TUI interface with interactive navigation
- Configurable time thresholds
- File type filtering
- Export results to JSON/CSV
- Recursive directory scanning
- File deletion/archiving actions
- macOS/Linux support

## License

This project is open source and available under the MIT License.

## Support

Found a bug or have a feature request? Please [open an issue](https://github.com/yourusername/go-filesystem/issues).
