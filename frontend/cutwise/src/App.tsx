import "./App.css";
import { listen } from '@tauri-apps/api/event';
import { Navbar } from "./components/navbar.tsx";
import { useState } from "react";
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Settings } from "./pages/Settings.tsx";
import {Home} from "./pages/Home.tsx";
import {LocalJobs} from "./pages/LocalJobs.tsx";
import {CloudJobs} from "./pages/CloudJobs.tsx";
// import Home from './pages/Home';
// import AllJobs from './pages/AllJobs';
// import JobDetail from './pages/JobDetail';

function App() {
    const [jobData, setJobData] = useState();
    const [jobId, setJobId] = useState("");

    listen('open_job_event', (event): void => {
        console.log(event.payload); // This will log "Open Job menu item clicked"
        // Add your JavaScript function call or additional functionality here
    });

    return (
        <Router>
            <div className="container">
                <Navbar setJobId={setJobId} jobId={jobId} setJobData={setJobData} />
                <Routes>
                    {/* Uncomment and use the appropriate routes */}
                     <Route path="/" element={<Home />} />
                     <Route path="/localjobs" element={<LocalJobs />} />
                     <Route path="/cloudjobs" element={<CloudJobs />} />
                    {/* <Route path="/jobs/:jobId" element={<JobDetail />} /> */}
                    <Route path="/settings" element={<Settings />} />
                </Routes>
            </div>
        </Router>
    );
}

export default App;
