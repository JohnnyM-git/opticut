use tauri::api::file::write_file;
use tauri::api::path::BaseDirectory;
use std::fs::File;
use std::io::Write;

#[tauri::command]
pub async fn select_and_upload_file(file_data: Vec<u8>, file_name: String) -> Result<String, String> {
    // Save the file data to a temporary file or desired location
    let file_path = BaseDirectory::Temp.resolve_path(file_name).map_err(|e| e.to_string())?;
    let mut file = File::create(&file_path).map_err(|e| e.to_string())?;
    file.write_all(&file_data).map_err(|e| e.to_string())?;

    // Call your Go backend to handle the file
    let client = reqwest::Client::new();
    let form = reqwest::multipart::Form::new()
        .file("file", file_path.to_string_lossy())
        .map_err(|e| e.to_string())?;

    let response = client
        .post("http://localhost:2828/file-upload")
        .multipart(form)
        .send()
        .await
        .map_err(|e| e.to_string())?;

    let response_text = response.text().await.map_err(|e| e.to_string())?;

    Ok(response_text)
}
