# Font Preview Generator (`font-gen`)

A high-performance, pure-Go CLI utility designed to generate rapid typography image previews from varied font sources. It handles inputs dynamically—including local font files, remote server URLs, and Fontsource registry configurations—and streams raw WebP binary data straight down the standard output pipeline (`stdout`).

This project is fully optimized for cross-compilation with zero CGO dependencies, enabling instant integration as an underlying engine inside runtime environments like Node.js and TypeScript.

---

## 🚀 Key Features

- **Strict CLI Stream Isolation**: Separates clean, raw binary asset payloads (`stdout`) from operational/runtime diagnostic outputs (`stderr`).
- **Multi-Source Font Processing**: Supports direct paths to local `.ttf`/`.otf` files, public web paths (`http://` or `https://`), and Fontsource registry package syntax.
- **On-the-Fly Decompression**: Dynamically sniffs binary headers for modern compressed web components (`wOF2`) and decompresses them natively.
- **Pure-Go Architecture**: Zero dependency on external system libraries or CGO, producing ultra-lightweight binaries that cross-compile perfectly.
- **Enterprise-Grade Automation**: Integrates Changeset-driven local semantic versioning alongside GitHub Actions automation.

---

## 📦 Installation & Pre-compiled Binaries

Pre-compiled, zero-dependency binaries are cross-compiled automatically for production deployment. You can download the latest executable matching your platform directly from the GitHub Releases page.

### Quick Automated Installation (CLI)

Run one of the following commands in your terminal to automatically detect your system architecture and download the latest release of `font-gen` into your current directory:

#### For Linux / macOS:
```bash
# Downloads the latest stable executable dynamically via the GitHub API redirect layer
curl -s https://api.github.com/repos/Electricitys/font-preview-gen/releases/latest \
| grep "browser_download_url" \
| grep "$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')" \
| cut -d '"' -f 4 \
| xargs curl -L -o font-gen && chmod +x font-gen
```

---

## 📂 Project Architecture

The repository adheres to the Standard Go Project Layout to isolate internal libraries from operational entry points:

```text
├── .changesets/             # Tracked intent-to-change versioning fragments
├── .github/workflows/       # GitHub Actions automated release pipeline configuration
├── cmd/
│   └── font-gen/
│       └── main.go          # The thin CLI wrapper, flag configuration, and stream routing
├── internal/
│   ├── downloader/          # Handles remote font fetches & remote registries
│   ├── processor/           # Validates binary magic-bytes and decompresses WOFF2
│   └── renderer/            # Pure Go font typography rasterization to WebP
├── Makefile                 # Development automation toolset
├── version.txt              # Unified dynamic single source of truth for versions
└── CHANGELOG.md             # Automated record of product modifications

```

---

## 🛠️ Prerequisites

* **Go**: Version `1.21` or newer.
* **Changeset CLI**: Required for local automated deployment routines.

---

## 💻 How to Use

The application takes parameters via custom sub-commands and processes flags, defaulting to streaming clean bytes to the console if an explicit output path isn't declared.

### 1. Basic CLI Redirection Syntax

```bash
# Process a font source and redirect the clean stdout buffer directly into an image asset
./dist/font-gen "<font_source>" "Preview Text Here" > preview.webp

```

### 2. Supported Configuration Formats

```bash
# Mode A: Local Filesystem Font
./dist/font-gen "./fonts/Inter-Regular.ttf" "Hello World" > output.webp

# Mode B: Fontsource Registry Token
./dist/font-gen "roboto" "Typography Preview" > output.webp

# Mode C: Direct External Web URL (Auto-sniffs and handles WOFF2 decompression)
./dist/font-gen "https://fonts.gstatic.com/s/inter/v13/UcCO3FwrK3iLTeHuS_fvQtMwCp5SR3p0.woff2" "Web Font" > output.webp

```

### 3. Native TypeScript Integration Wrapper

Because `font-gen` segregates data layers from errors cleanly, you can pipe it seamlessly into Node.js runtime environments using `child_process.execFile`:

```typescript
import { execFile } from 'child_process';

function generateFontPreview(fontSource: string, text: string): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    // Invoke the Go binary directly, omitting an output file flag to stream into stdout
    execFile(
      './dist/font-gen', 
      [fontSource, text], 
      { encoding: 'buffer' }, 
      (error, stdout, stderr) => {
        if (error) {
          // Operational problems or errors surfaced cleanly over stderr
          return reject(new Error(stderr.toString('utf8')));
        }
        // Returns the raw binary image buffer isolated perfectly from stdout
        resolve(stdout);
      }
    );
  });
}

```

---

## 🔧 Development Workflows

The project contains a comprehensive automation task runner managed through the `Makefile`.

### Local Testing and Building

Run quality sweeps and create an optimized binary stripped of debug information for your current machine architecture:

```bash
# Run tidy checks, format verification, race testing, and local compilation at once
make all

# Run specific individual milestones
make tidy      # Cleanup module dependencies
make fmt       # Synchronize code structures
make test      # Fire off localized diagnostic testing
make build     # Compile optimized local binary into /dist

```

### Dynamic Execution

To stream test applications dynamically without triggering a standalone build step manually, use the `run` directive appended with custom environment arguments:

```bash
make run ARGS="'roboto' 'Testing System Automation'"

```

---

## 📦 Lifecycle Releases & Versioning

The codebase leverages a structured two-step validation paradigm to handle code shifts, updates, changelog tracking, and artifact asset compilation automatically.

### Step 1: Create an Intention Fragment (Local Development)

Whenever changes are made, draft an automated intent fragment before pushing your branch:

```bash
make changeset

```

### Step 2: Consumer Automation & Publishing

When you are ready to prepare a production deployment cycle, clean up fragments to update version control variables and tag the Git tree locally:

```bash
# 1. Bump the text versions and compile the changelogs
make version-bump

# 2. Commit files, generate local tags, and push upstream safely
make git-release

```

Once pushed, your configured **GitHub Actions Workflow** pipeline fires automatically to cross-compile distribution binaries for the listed architectures:

* `linux/amd64`
* `linux/arm64`
* `windows/amd64`
* `darwin/arm64`

The pipeline will immediately package the binaries under a compilation matrix with `CGO_ENABLED=0` and attach them flawlessly to your repository's **GitHub Releases** portal under both the specific semantic version tag (e.g. `v0.1.2`) and the rolling `latest` tag pointer.

---

## 🧼 Housekeeping

To scrub distribution artifacts, delete previous compiled elements, or clear out the local `/dist` directory safely:

```bash
make clean

```
