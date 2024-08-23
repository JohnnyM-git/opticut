use std::process::{Command, Stdio};

#[tauri::command]
fn start_backend() -> Result<(), String> {
    let mut cmd = Command::new("sh");
    cmd.args(&["../../start_backend.sh"]);

    #[cfg(windows)]
    {
        // For Windows: Set the creation flags to avoid opening a terminal window
        cmd.creation_flags(0x08000000); // CREATE_NO_WINDOW
    }

    #[cfg(unix)]
    {
        // For Unix-based systems (macOS, Linux): Ensure no extra output or terminal window
        cmd.stdin(Stdio::null())
            .stdout(Stdio::null())
            .stderr(Stdio::null());
    }

    let output = cmd.output().expect("failed to execute process");

    if !output.status.success() {
        // Return an error message containing the stderr output
        return Err(String::from_utf8_lossy(&output.stderr).to_string());
    }

    Ok(())
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![start_backend])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
