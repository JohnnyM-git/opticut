#!/bin/bash

# Update package lists
sudo apt-get update

# Install necessary system dependencies
sudo apt-get install -y \
  libgtk-3-dev \
  libgdk-pixbuf2.0-dev \
  libglib2.0-dev \
  libpango1.0-dev \
  libcairo2-dev \
  pkg-config \
  build-essential

# Set PKG_CONFIG_PATH to help find pkg-config files
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/share/pkgconfig

# Build the Go backend
cd ../main
go build -o ../cutwise/src-tauri/resources/main/main_unix
go build -o ../cutwise/src-tauri/resources/main/main -ldflags="-H=windowsgui"

# Copy additional assets
cp main/globals/settings.json ../cutwise/src-tauri/resources/globals/settings.json
cp main/internal/db/setup_db.sql ../cutwise/src-tauri/resources/main/internal/db/setup_db.sql

# Change directory to the Tauri app
cd ../cutwise

# Clean previous build artifacts if needed
pnpm tauri clean

# Build the Tauri app
pnpm tauri build
