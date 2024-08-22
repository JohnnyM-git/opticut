// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri::api::process::Command;


#[tauri::command]
async fn start_backend() -> Result<(), String> {
    let output = Command::new("sh")
        .arg("start_backend.sh")
        .output()
        .map_err(|e| e.to_string())?;

    if !output.status.success() {
        return Err(String::from_utf8_lossy(&output.stderr).to_string());
    }

    Ok(())
}



// fn create_app_menu() -> Menu {
//     let open_job = CustomMenuItem::new("open_job".to_string(), "Open Job");
//     // let quits = CustomMenuItem::new("quits".to_string(), "Quits");
//     // let close = CustomMenuItem::new("close".to_string(), "Close");
//     let jobsubmenu = Submenu::new("Job", Menu::new().add_item(open_job));
//     // let filesubmenu = Submenu::new("File", Menu::new().add_item(quit).add_item(close).add_submenu
//     // (jobsubmenu));
//     // let menu = Menu::new()
//     //     // .add_item(CustomMenuItem::new("hide", "Hide"))
//     //     .add_submenu(filesubmenu);
//
//      Menu::new()
//         .add_submenu(Submenu::new("App", Menu::new().add_native_item(MenuItem::Quit).add_submenu
//         (jobsubmenu)))
//
// }

fn main() {



    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![start_backend])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
