# Navigate to the backend directory
Set-Location -Path "backend\main"

# Build the Go backend
go build -o backend main.go

# Start the Go backend and redirect output to a log file
Start-Process -FilePath ".\backend.exe" -ArgumentList "> backend.log 2>&1" -NoNewWindow -RedirectStandardOutput "backend.log" -RedirectStandardError "backend.log"
