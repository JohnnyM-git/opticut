#!/bin/bash

# Detect the operating system
OS=$(uname -s)

if [ "$OS" = "Darwin" ]; then
    echo "Building for macOS..."
    cd ../main
    go build -o ../cutwise/src-tauri/resources/main/main_unix
elif [ "$OS" = "Linux" ]; then
    echo "Building for Linux..."
    cd ../main
    go build -o ../cutwise/src-tauri/resources/main/main_unix
elif [[ "$OS" == *"NT"* || "$OS" == "MINGW"* || "$OS" == "CYGWIN"* ]]; then
    echo "Building for Windows..."
    cd ../main
    go build -o ../cutwise/src-tauri/resources/main/main.exe -ldflags=-H=windowsgui
else
    echo "Unsupported OS: $OS"
    exit 1
fi
