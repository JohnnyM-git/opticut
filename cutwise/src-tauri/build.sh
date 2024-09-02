#!/bin/bash

# Build the Linux executable
cd ../main && go build -o ../cutwise/src-tauri/resources/main/main_unix

# Build the Windows executable (optional in a Linux CI environment)
go build -o ../cutwise/src-tauri/resources/main/main -ldflags="-H=windowsgui"

# Copy the settings.json file
cp ../main/globals/settings.json ../cutwise/src-tauri/resources/globals/settings.json

# Copy the setup_db.sql file
cp ../main/internal/db/setup_db.sql ../cutwise/src-tauri/resources/main/internal/db/setup_db.sql
