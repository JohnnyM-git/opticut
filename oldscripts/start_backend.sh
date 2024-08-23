#!/bin/bash

## Print the current directory
#echo "Current directory: $(pwd)"

# Navigate to the backend directory
cd backend/main || exit 1

# Build the Go backend
go build -o backend main.go || exit 1

# Start the Go backend in the background and redirect output to a log file
nohup ./backend > backend.log 2>&1 &

# Navigate to the frontend directory (uncomment if needed)
# cd ../../cutwise

# Start the Tauri app (uncomment if needed)
# pnpm tauri dev
