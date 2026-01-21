#!/bin/bash

# urlX Installation Script
# Author: Alham Rizvi
# Repository: https://github.com/alhamrizvi-cloud/urlx
# Usage: curl -sSfL https://raw.githubusercontent.com/alhamrizvi-cloud/urlx/main/install.sh | sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="alhamrizvi-cloud/urlx"
BINARY_NAME="urlx"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}"
cat << "EOF"
 __  __  _____   _       __  __
/\ \/\ \/\  __'\/\ \     /\_\/\_\
\ \ \_\ \ \ \L\ \ \ \    \/_/\/_/
 \ \  _  \ \ ,  /\ \ \  __  /\_\
  \ \ \ \ \ \ \\ \\ \ \L\ \ \ \ \
   \ \_\ \_\ \_\ \_\ \____/  \ \_\
    \/_/\/_/\/_/\/ /\/___/    \/_/

    Installation Script
    Created by Alham Rizvi
EOF
echo -e "${NC}"

# Get latest release version
echo -e "${YELLOW}[*] Fetching latest release...${NC}"
LATEST_VERSION=$(curl -sSf "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}[!] Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "${GREEN}[+] Latest version: $LATEST_VERSION${NC}"

# Construct download URL
if [ "$OS" = "darwin" ]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME}-${LATEST_VERSION}-darwin-${ARCH}.tar.gz"
elif [ "$OS" = "linux" ]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME}-${LATEST_VERSION}-linux-${ARCH}.tar.gz"
else
    echo -e "${RED}[!] Unsupported OS: $OS${NC}"
    echo -e "${YELLOW}Please download manually from: https://github.com/$REPO/releases${NC}"
    exit 1
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download binary
echo -e "${YELLOW}[*] Downloading $BINARY_NAME...${NC}"
if ! curl -sSfL "$DOWNLOAD_URL" -o "$TMP_DIR/${BINARY_NAME}.tar.gz"; then
    echo -e "${RED}[!] Download failed${NC}"
    echo -e "${YELLOW}URL: $DOWNLOAD_URL${NC}"
    exit 1
fi

# Extract binary
echo -e "${YELLOW}[*] Extracting...${NC}"
tar -xzf "$TMP_DIR/${BINARY_NAME}.tar.gz" -C "$TMP_DIR"

# Find the binary (it might have different naming)
BINARY_PATH=$(find "$TMP_DIR" -name "${BINARY_NAME}*" -type f | head -n 1)

if [ -z "$BINARY_PATH" ]; then
    echo -e "${RED}[!] Binary not found in archive${NC}"
    exit 1
fi

# Install binary
echo -e "${YELLOW}[*] Installing to $INSTALL_DIR...${NC}"

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${YELLOW}[*] Root permissions required for installation${NC}"
    sudo mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo -e "${GREEN}[+] Successfully installed $BINARY_NAME${NC}"
    echo -e "${GREEN}[+] Version: $($BINARY_NAME -h | head -n 1 2>/dev/null || echo $LATEST_VERSION)${NC}"
    echo ""
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $BINARY_NAME -domain example.com"
    echo "  $BINARY_NAME -d domains.txt -crawl -probe -v"
    echo "  $BINARY_NAME -domain example.com -crawl -only '.pdf' -o pdfs.txt"
    echo ""
    echo -e "${YELLOW}For more information:${NC}"
    echo "  $BINARY_NAME -h"
    echo "  https://github.com/$REPO"
    echo ""
    echo -e "${GREEN}Created by Alham Rizvi${NC}"
else
    echo -e "${RED}[!] Installation completed but binary not found in PATH${NC}"
    echo -e "${YELLOW}You may need to add $INSTALL_DIR to your PATH${NC}"
    exit 1
fi
