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
	@changeset version --project font-preview-gen
	@echo "=> New resolved project version is: $$(cat version.txt)"

.PHONY: git-release
git-release:
	@NEW_VERSION=$$(cat version.txt); \
	echo "=> Committing and publishing version v$${NEW_VERSION} to GitHub..."; \
	git add version.txt CHANGELOG.md .changesets/; \
	git commit -m "chore: release v$${NEW_VERSION}" || echo "=> No structural file changes to commit"; \
	git tag -f v$${NEW_VERSION}; \
	git tag -f latest; \
	echo "=> Pushing commits and release tags to remote origin..."; \
	git push origin HEAD; \
	git push origin v$${NEW_VERSION} --force; \
	git push origin latest --force; \
	echo "=> Successfully pushed! GitHub Actions will now build and index v$${NEW_VERSION} and sync the 'latest' release point."

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
