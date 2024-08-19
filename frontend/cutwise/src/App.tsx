// import { useState } from "react";
// import { invoke } from "@tauri-apps/api/tauri";
import "./App.css";
// import {Navbar} from "./components/navbar.tsx";
import { listen } from '@tauri-apps/api/event';
import {Navbar} from "./components/navbar.tsx";
import {useState} from "react";



function App() {
    const [jobData, setJobData] = useState()
    const [jobId, setJobId] = useState("")

  // async function getJobData(job: string) {
  //     const response = await fetch(`http://localhost:8080/api/v1/job?job_id=${job}`);
  //     const data = await response.json();
  //     setJobData(data)
  //     console.log(jobData);
  // }


    listen('open_job_event', (event):void  => {
        console.log(event.payload); // This will log "Open Job menu item clicked"
        // Add your JavaScript function call or additional functionality here
    });


    return (

      <div className="container">
    <Navbar setJobId={setJobId} jobId={jobId} setJobData={setJobData} />

      </div>
  );
}

export default App;
