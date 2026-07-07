#!/bin/bash
PACKAGE_NAME="font-gen"
PLATFORMS=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/arm64")

# Ensure the dist folder exists
mkdir -p dist

for platform in "${PLATFORMS[@]}"; do
    IFS="/" read -r -a split <<< "$platform"
    GOOS=${split[0]}
    GOARCH=${split[1]}

    OUTPUT_NAME="${PACKAGE_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi

    echo "Compiling for ${GOOS}/${GOARCH}..."

    # 1. Map Go targets to Zig triple targets
    case "$platform" in
        "linux/amd64")  ZIG_TARGET="x86_64-linux-gnu" ;;
        "linux/arm64")  ZIG_TARGET="aarch64-linux-gnu" ;;
        "windows/amd64") ZIG_TARGET="x86_64-windows-gnu" ;;
        "darwin/arm64")  ZIG_TARGET="aarch64-macos" ;;
        *)               ZIG_TARGET="" ;;
    esac

    # 2. Build using explicit CGO options
    if [ -n "$ZIG_TARGET" ]; then
        env CGO_ENABLED=1 \
            GOOS=$GOOS \
            GOARCH=$GOARCH \
            CC="zig cc -target $ZIG_TARGET" \
            CXX="zig c++ -target $ZIG_TARGET" \
            go build -o dist/$OUTPUT_NAME .
    else
        # Fallback to standard cross-compilation if the target mapping isn't defined
        env GOOS=$GOOS GOARCH=$GOARCH go build -o dist/$OUTPUT_NAME .
    fi
done
