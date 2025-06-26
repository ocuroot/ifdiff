#!/bin/bash
set -e

VERSION="v0.0.2"
BINARY_NAME="ifdiff"
BUILD_DIR="./dist"
MAIN_PATH="./cmd/ifdiff/main.go"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build for each target OS/architecture
build_and_package() {
    local os=$1
    local arch=$2
    local ext=$3

    echo "Building for $os/$arch..."
    
    # Set extension based on OS (.exe for Windows)
    local binary_ext=""
    if [ "$os" == "windows" ]; then
        binary_ext=".exe"
    fi
    
    # Build binary
    GOOS=$os GOARCH=$arch go build -o ${BUILD_DIR}/${BINARY_NAME}${binary_ext} ${MAIN_PATH}
    
    # Create archive directory
    local pkg_name="${BINARY_NAME}-${VERSION}-${os}-${arch}"
    local pkg_dir="${BUILD_DIR}/${pkg_name}"
    mkdir -p ${pkg_dir}
    
    # Copy binary and license to archive directory
    cp ${BUILD_DIR}/${BINARY_NAME}${binary_ext} ${pkg_dir}/

    # Create archive
    cd ${BUILD_DIR}
    if [ "$ext" == "zip" ]; then
        zip -r "${pkg_name}.zip" "${pkg_name}"
    else
        tar -czf "${pkg_name}.tar.gz" "${pkg_name}"
    fi
    cd ..
    
    # Generate SHA256 checksum
    if [ "$os" == "darwin" ] || [ "$os" == "linux" ]; then
        shasum -a 256 "${BUILD_DIR}/${pkg_name}.${ext}" | awk '{print $1}' > "${BUILD_DIR}/${pkg_name}.${ext}.sha256"
    else
        # For Windows, use different sha256 command
        openssl dgst -sha256 "${BUILD_DIR}/${pkg_name}.${ext}" | awk '{print $2}' > "${BUILD_DIR}/${pkg_name}.${ext}.sha256"
    fi
    
    # Clean up temporary files
    rm -rf ${pkg_dir}
    rm -f ${BUILD_DIR}/${BINARY_NAME}${binary_ext}
    
    # Display info
    echo "${pkg_name}.${ext} created."
    echo "sha256:$(cat ${BUILD_DIR}/${pkg_name}.${ext}.sha256)"
    echo "$(du -h ${BUILD_DIR}/${pkg_name}.${ext} | cut -f1)"
}

echo "Building release binaries for ${BINARY_NAME} ${VERSION}"

# Build for Darwin (macOS)
build_and_package "darwin" "amd64" "tar.gz"
build_and_package "darwin" "arm64" "tar.gz"

# Build for Linux
build_and_package "linux" "amd64" "tar.gz"
build_and_package "linux" "arm64" "tar.gz"

# Build for Windows
build_and_package "windows" "amd64" "zip"

echo "Release builds complete!"
