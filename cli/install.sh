#!/bin/bash
# Install mytodolist to /usr/local/bin

BINARY_PATH="./mytodolist"
INSTALL_DIR="/usr/local/bin"
NEWCONFIG_DIR="/.cfg/mytodolist"
CONFIG_DIR="./.cfg/config.json"

mkdir -p $NEWCONFIG_DIR

cp "$BINARY_PATH" "$INSTALL_DIR/"
cp -r "$CONFIG_DIR" "$NEWCONFIG_DIR/"

# Make it executable
chmod +x "$INSTALL_DIR/mytodolist"

# Permission to write and read the config
chmod 666 /.cfg/mytodolist/config.json

echo "mytodolist has been installed to $INSTALL_DIR"
