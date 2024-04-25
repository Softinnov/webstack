#!/bin/bash
# Install mytodolist to /usr/local/bin

BINARY_PATH="./mytodolist"

INSTALL_DIR="/usr/local/bin"

cp "$BINARY_PATH" "$INSTALL_DIR/mytodolist"

# Make it executable
chmod +x "$INSTALL_DIR/mytodolist"

echo "myapp has been installed to $INSTALL_DIR"
