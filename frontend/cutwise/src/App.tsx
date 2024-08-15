// import { useState } from "react";
// import { invoke } from "@tauri-apps/api/tauri";
import "./App.css";
import {Navbar} from "./components/navbar.tsx";
import { listen } from '@tauri-apps/api/event';



function App() {
  // const [greetMsg, setGreetMsg] = useState("");
  // const [name, setName] = useState("");

  // async function greet() {
  //   // Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
  //   setGreetMsg(await invoke("greet", { name }));
  // }

  // async function getJobData() {
  //     const response = await fetch(`http://localhost:8080/api/v1/job?job_id=TEST`);
  //     const data = await response.json();
  //     console.log(data);
  // }


    listen('open_job_event', (event) => {
        console.log(event.payload); // This will log "Open Job menu item clicked"
        // Add your JavaScript function call or additional functionality here
    });


    return (

      <div className="container">

          {/*<div data-tauri-drag-region className="titlebar">*/}
          {/*    /!*<div className="titlebar-button" id="titlebar-minimize">*!/*/}
          {/*    /!*    <img*!/*/}
          {/*    /!*        src="https://api.iconify.design/mdi:window-minimize.svg"*!/*/}
          {/*    /!*        alt="minimize"*!/*/}
          {/*    /!*    />*!/*/}
          {/*    /!*</div>*!/*/}
          {/*    /!*<div className="titlebar-button" id="titlebar-maximize">*!/*/}
          {/*    /!*    <img*!/*/}
          {/*    /!*        src="https://api.iconify.design/mdi:window-maximize.svg"*!/*/}
          {/*    /!*        alt="maximize"*!/*/}
          {/*    /!*    />*!/*/}
          {/*    /!*</div>*!/*/}
          {/*    /!*<div className="titlebar-button" id="titlebar-close">*!/*/}
          {/*    /!*    <img src="https://api.iconify.design/mdi:close.svg"*!/*/}
          {/*    /!*         alt="close"/>*!/*/}
          {/*    /!*</div>*!/*/}
          {/*</div>*/}
          <Navbar/>

      </div>
  );
}

export default App;
