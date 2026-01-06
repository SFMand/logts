#!/bin/bash
set -e

OWNER="SFMand"
REPO="logts"
BINARY="logts"
INSTALL_DIR="/usr/local/bin"

OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    Linux)  OS="Linux" ;;
    Darwin) OS="Darwin" ;;
    *)      echo "OS $OS not supported"; exit 1 ;;
esac

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    arm64)  ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;; # Linux often reports arm64 as aarch64
    *)      echo "Architecture $ARCH not supported"; exit 1 ;;
esac

# GoReleaser default format: logts_Linux_x86_64.tar.gz
FILE_NAME="${BINARY}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${FILE_NAME}"

echo "Downloading $BINARY from $DOWNLOAD_URL..."

TMP_DIR=$(mktemp -d)
curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$FILE_NAME"

echo "Extracting..."
tar -xzf "$TMP_DIR/$FILE_NAME" -C "$TMP_DIR"

echo "Installing to $INSTALL_DIR (requires password)..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
else
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
fi

chmod +x "$INSTALL_DIR/$BINARY"
rm -rf "$TMP_DIR"

echo "Success, $BINARY has been installed to $INSTALL_DIR"
echo "Try running: $BINARY --help"