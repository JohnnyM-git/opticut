#!/bin/bash

# Install system dependencies
sudo apt-get update
sudo apt-get install -y libgtk-3-dev

# Build the Go backend
cd ../main && go build -o ../cutwise/src-tauri/resources/main/main_unix
go build -o ../cutwise/src-tauri/resources/main/main -ldflags="-H=windowsgui"

# Copy additional assets
cp main/globals/settings.json ../cutwise/src-tauri/resources/globals/settings.json
cp main/internal/db/setup_db.sql ../cutwise/src-tauri/resources/main/internal/db/setup_db.sql

# Build the Tauri app
cd ../cutwise && pnpm tauri build