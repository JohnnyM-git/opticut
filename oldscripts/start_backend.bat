@echo off

:: Navigate to the backend directory
cd .\Backend\main

:: Build the Go backend
go build -o backend main.go

:: Start the Go backend
start backend