#!/bin/bash

## Print the current directory
#echo "Current directory: $(pwd)"

# Navigate to the backend directory
cd ./resource/main || exit 1



# Start the Go backend in the background and redirect output to a log file
nohup ./main_unix > backend.log 2>&1 &
