import "./App.css";
import { Navbar } from "./components/navbar.tsx";
import { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Settings } from "./pages/Settings.tsx";
import { Home } from "./pages/Home.tsx";
import { LocalJobs } from "./pages/LocalJobs.tsx";
import { CloudJobs } from "./pages/CloudJobs.tsx";
import { Results } from "./pages/Results.tsx";
import { invoke } from "@tauri-apps/api";

function App() {

  const [jobId, setJobId] = useState("");
    // const [backendStatus, setBackendStatus] = useState('starting');
    //
    // useEffect(() => {
    //     async function startBackend() {
    //         try {
    //             await invoke("start_backend");
    //             // Check if backend is running
    //             const res = await fetch('http://localhost:2828/api/v1/hello'); // Replace with actual health check endpoint
    //             if (res.ok) {
    //                 setBackendStatus('running');
    //             } else {
    //                 setBackendStatus('failed');
    //             }
    //         } catch (error) {
    //             console.error("Failed to start the backend:", error);
    //             setBackendStatus('failed');
    //         }
    //     }
    //
    //     startBackend();
    // }, []);

  return (
    <Router>
        <div className="container">
            <Navbar setJobId={setJobId} jobId={jobId}/>
            <Routes>
                <Route path="/" element={<Home/>} />
          <Route path="/localjobs" element={<LocalJobs />} />
          <Route path="/cloudjobs" element={<CloudJobs />} />
          <Route path="/results/:jobId" element={<Results />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
