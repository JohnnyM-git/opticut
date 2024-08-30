#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::process::{Command, Child, Stdio};
use std::sync::{Arc, Mutex};
use tauri::Manager;
#[cfg(windows)]
use std::os::windows::process::CommandExt;
use std::path::PathBuf;
use std::fs::{create_dir_all, OpenOptions};
use std::io::prelude::*;
use tauri::api::path::{data_dir, };
use chrono::Local;

// Shared state to keep track of the backend process
struct BackendProcess(Arc<Mutex<Option<Child>>>);

#[tauri::command]
fn start_backend(state: tauri::State<BackendProcess>) -> Result<(), String> {
    #[cfg(windows)]
    {
        // Use the original path for start_backend.bat
        let mut bat_file_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
        bat_file_path.push("src/start_backend.bat");

        let mut cmd = Command::new("cmd");
        cmd.args(&["/C", bat_file_path.to_str().unwrap()])
            .creation_flags(0x08000000); // CREATE_NO_WINDOW to avoid opening a terminal window

        let child = cmd.spawn();

        match child {
            Ok(child) => {
                // Store the child process in the state
                let mut process = state.0.lock().unwrap();
                *process = Some(child);
            }
            Err(e) => {
                // Log the error
                let error_message = format!("Failed to start backend process: {}", e);
                log_to_file(error_message.clone()).unwrap_or_else(|err| {
                    eprintln!("Failed to log error: {}", err);
                });
                return Err(error_message);
            }
        }
    }

    #[cfg(unix)]
    {
        // Use relative path based on CARGO_MANIFEST_DIR
        let mut sh_file_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
        sh_file_path.push("src/start_backend.sh");

        let child = Command::new("sh")
            .args(&[sh_file_path.to_str().unwrap()])
            .stdin(Stdio::null())
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .spawn();

        match child {
            Ok(child) => {
                // Store the child process in the state
                let mut process = state.0.lock().unwrap();
                *process = Some(child);
            }
            Err(e) => {
                // Log the error
                let error_message = format!("Failed to start backend process: {}", e);
                log_to_file(error_message.clone()).unwrap_or_else(|err| {
                    eprintln!("Failed to log error: {}", err);
                });
                return Err(error_message);
            }
        }
    }

    Ok(())
}

#[tauri::command]
fn stop_backend(state: tauri::State<BackendProcess>) -> Result<(), String> {
    #[cfg(windows)]
    {
        // Use the original path for stop_backend.bat
        let mut bat_file_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
        bat_file_path.push("src/stop_backend.bat");

        let mut cmd = Command::new("cmd");
        cmd.args(&["/C", bat_file_path.to_str().unwrap()])
            .creation_flags(0x08000000); // CREATE_NO_WINDOW to avoid opening a terminal window

        let _ = cmd.spawn().expect("failed to stop backend process");

        // Remove the child process from the state
        let mut process = state.0.lock().unwrap();
        if let Some(mut child) = process.take() {
            let _ = child.kill();
        }
    }

    #[cfg(unix)]
    {
        // Use relative path based on CARGO_MANIFEST_DIR
        let mut sh_file_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
        sh_file_path.push("src-tauri/resources/stop_backend.sh");

        let child = Command::new("sh")
            .args(&[sh_file_path.to_str().unwrap()])
            .stdin(Stdio::null())
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .spawn()
            .expect("failed to stop backend process");

        // Store the child process in the state
        let mut process = state.0.lock().unwrap();
        *process = Some(child);
    }

    let _ = log_to_file("Backend server has stopped.".to_string());

    Ok(())
}

#[tauri::command]
fn log_to_file(message: String) -> Result<(), String> {
    // Get the application's data directory
    let data_dir = data_dir().ok_or("Failed to get data directory")?;

    // Construct the log directory path
    let log_dir = data_dir.join("com.cutwise.app").join("logs");

    // Ensure the directory exists
    if let Err(e) = create_dir_all(&log_dir) {
        return Err(format!("Failed to create directory: {}", e));
    }

    // Construct the log file path
    let log_path = log_dir.join("server.log");

    // Open the log file in append mode
    let mut file = match OpenOptions::new().append(true).create(true).open(&log_path) {
        Ok(f) => f,
        Err(e) => return Err(format!("Failed to open log file: {}", e)),
    };

    // Get the current date and time
    let now = Local::now();
    let timestamp = now.format("[%Y-%m-%d %H:%M:%S]").to_string();

    // Create the log entry with the date and time prepended
    let log_entry = format!("{} {}\n", timestamp, message);

    // Write the log entry to the file
    if let Err(e) = file.write_all(log_entry.as_bytes()) {
        return Err(format!("Failed to write log message: {}", e));
    }

    Ok(())
}

fn main() {
    let backend_process = BackendProcess(Arc::new(Mutex::new(None)));

    tauri::Builder::default()
        .manage(backend_process)
        .invoke_handler(tauri::generate_handler![start_backend, stop_backend, log_to_file])
        .on_window_event(|event| {
            if let tauri::WindowEvent::CloseRequested { .. } = event.event() {
                let app = event.window().app_handle();
                let backend_process = app.state::<BackendProcess>();
                let _ = stop_backend(backend_process); // Stop the backend when the window is closed
            }
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}