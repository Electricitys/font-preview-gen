# ====================================================================================
# Configuration Variables
# ====================================================================================
BINARY_NAME=font-gen
CMD_DIR=./cmd/font-gen
DIST_DIR=dist

# Dynamic Version Lookup from version.txt
VERSION?=$(shell cat version.txt 2>/dev/null || echo "0.1.0")

# Core Build Flags
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION}"

# Supported Platforms for Cross-Compilation
PLATFORMS=linux/amd64 linux/arm64 windows/amd64 darwin/arm64

# ====================================================================================
# Development & Quality Targets
# ====================================================================================

.PHONY: all
all: clean tidy fmt test build

.PHONY: tidy
tidy:
	@echo "=> Tidying up Go modules..."
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "=> Formatting code structures..."
	@go fmt ./...

.PHONY: test
test:
	@echo "=> Running test suites..."
	@go test -v -race ./...

# ====================================================================================
# Local Compilations
# ====================================================================================

.PHONY: build
build: tidy
	@echo "=> Compiling local binary [${BINARY_NAME}] version [${VERSION}]..."
	@mkdir -p ${DIST_DIR}
	@go build ${LDFLAGS} -o ${DIST_DIR}/${BINARY_NAME} ${CMD_DIR}
	@echo "=> Built successfully to ${DIST_DIR}/${BINARY_NAME}"

.PHONY: run
run: build
	@echo "=> Launching local instance with arguments..."
	@./${DIST_DIR}/${BINARY_NAME} $(ARGS)

# ====================================================================================
# Changeset-Driven Version Control Targets
# ====================================================================================

.PHONY: changeset
changeset:
	@echo "=> Creating a new intent-to-change fragment..."
	@changeset add

.PHONY: version-bump
version-bump:
	@echo "=> Consuming fragment files and updating version strings..."
	@changeset version
	@echo "=> New resolved project version is: $$(cat version.txt)"

# ====================================================================================
# Cross-Compilation (Production Release)
# ====================================================================================

.PHONY: release
release: clean tidy
	@echo "=> Starting multi-platform release generation for version [${VERSION}]..."
	@mkdir -p ${DIST_DIR}
	@set -e; \
	for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		OUT_NAME="$(BINARY_NAME)-$(VERSION)-$${GOOS}-$${GOARCH}"; \
		if [ "$$GOOS" = "windows" ]; then OUT_NAME="$${OUT_NAME}.exe"; fi; \
		echo "   -> Building for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 go build $(LDFLAGS) -o $(DIST_DIR)/$$OUT_NAME $(CMD_DIR); \
	done
	@echo "=> Release binaries generated in /${DIST_DIR} folder."

# ====================================================================================
# Cleanup Tasks
# ====================================================================================

.PHONY: clean
clean:
	@echo "=> Cleaning up build artifacts..."
	@rm -rf ${DIST_DIR}
