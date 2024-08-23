import "./App.css";
import { listen } from "@tauri-apps/api/event";
import { Navbar } from "./components/navbar.tsx";
import { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Settings } from "./pages/Settings.tsx";
import { Home } from "./pages/Home.tsx";
import { LocalJobs } from "./pages/LocalJobs.tsx";
import { CloudJobs } from "./pages/CloudJobs.tsx";
import { Results } from "./pages/Results.tsx";
import { invoke } from "@tauri-apps/api";
// import Home from './pages/Home';
// import AllJobs from './pages/AllJobs';
// import JobDetail from './pages/JobDetail';

function App() {
  // const [jobData, setJobData] = useState();
  const [jobId, setJobId] = useState("");

  useEffect(() => {
    async function startBackend() {
      try {
        await invoke("start_backend");
      } catch (error) {
        console.error("Failed to start the backend:", error);
      }
    }

    startBackend();
  }, []);

  listen("open_job_event", (event): void => {
    console.log(event.payload); // This will log "Open Job menu item clicked"
    // Add your JavaScript function call or additional functionality here
  });

  return (
    <Router>
      <div className="container">
        <Navbar setJobId={setJobId} jobId={jobId} />
        <Routes>
          {/* Uncomment and use the appropriate routes */}
          <Route path="/" element={<Home />} />
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