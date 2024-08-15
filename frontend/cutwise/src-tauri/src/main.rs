// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}
use tauri::{CustomMenuItem, Menu, MenuItem, Submenu};
// here `"quit".to_string()` defines the menu item id, and the second parameter is the menu item label.

fn create_app_menu() -> Menu {
    let open_job = CustomMenuItem::new("open_job".to_string(), "Open Job");
    // let quits = CustomMenuItem::new("quits".to_string(), "Quits");
    // let close = CustomMenuItem::new("close".to_string(), "Close");
    let jobsubmenu = Submenu::new("Job", Menu::new().add_item(open_job));
    // let filesubmenu = Submenu::new("File", Menu::new().add_item(quit).add_item(close).add_submenu
    // (jobsubmenu));
    // let menu = Menu::new()
    //     // .add_item(CustomMenuItem::new("hide", "Hide"))
    //     .add_submenu(filesubmenu);

     Menu::new()
        .add_submenu(Submenu::new("App", Menu::new().add_native_item(MenuItem::Quit).add_submenu
        (jobsubmenu)))

}

fn main() {
    // here `"quit".to_string()` defines the menu item id, and the second parameter is the menu item label.



    tauri::Builder::default()
        .menu(create_app_menu())
        .invoke_handler(tauri::generate_handler![greet])
        .on_menu_event(|event| match event.menu_item_id() {
            "open_job" => {
                // Functionality for the open_job menu item click
                let window = event.window();
                window.emit("open_job_event", "Open Job menu item clicked").unwrap();
                println!("Open Job menu item clicked");
            }
            _ => {}
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
