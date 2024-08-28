import "./App.css";
import { Navbar } from "./components/navbar.tsx";
import { useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { Settings } from "./pages/Settings.tsx";
import { Home } from "./pages/Home.tsx";
import { LocalJobs } from "./pages/LocalJobs.tsx";
import { CloudJobs } from "./pages/CloudJobs.tsx";
import { Results } from "./pages/Results.tsx";
import { invoke } from "@tauri-apps/api/tauri";
import {apiUrl} from "./globals.ts";
import {Status} from "./pages/Status.tsx";


function App() {
    // const [appData, setAppData] = useState<string>()
    // const [backendStatus, setBackendStatus] = useState<string>('starting');


    // useEffect(() => {
    //     startBackend()
    // }, []);

    async function startBackend() {
        try {
            await logToFile("Starting server...");

            // Start the backend asynchronously
            await invoke("start_backend");
            await logToFile("Server started.");

            // Poll the backend health endpoint until it's up or timeout
            const checkBackendStatus = async () => {
                try {
                    const res = await fetch(`${apiUrl}/hello`); // Replace with actual health check endpoint
                    if (res.ok) {
                        await logToFile("Server is running.");
                        return true;
                    }
                } catch (error) {
                    await logToFile("Server is not up yet.");
                }
                return false;
            };

            let attempts = 0;
            const maxAttempts = 3;
            const delay = 1000; // 1 second

            const pollBackend = async () => {
                while (attempts < maxAttempts) {
                    const success = await checkBackendStatus();
                    if (success) {
                        break;
                    }
                    attempts++;
                    await new Promise(res => setTimeout(res, delay));
                }

                if (attempts >= maxAttempts) {
                    await logToFile("Backend failed to start.");
                }
            };

            pollBackend();
        } catch (error) {
            console.error("Failed to start the backend:", error);
            await logToFile("Failed to start the backend.");
        }
    }




    async function logToFile(message: string) {
        try {
            await invoke('log_to_file', { message });
            console.log('Log message written successfully.');
        } catch (error) {
            console.error('Failed to write log:', error);
        }
    }





    // async function stopBackend() {
    //     try {
    //         console.log(backendStatus)
    //         // Start the backend asynchronously
    //         await invoke("stop_backend");
    //         console.log(backendStatus)
    //         // Poll the backend health endpoint until it's up or timeout
    //         const checkBackendStatus = async () => {
    //             try {
    //                 console.log(backendStatus)
    //                 const res = await fetch('http://localhost:2828/api/v1/hello'); // Replace with actual health check endpoint
    //                 if (res.ok) {
    //                     setBackendStatus('running');
    //                     console.log(backendStatus)
    //                     return true;
    //                 }
    //             } catch (error) {
    //                 console.log(backendStatus)
    //                 console.log(error)
    //                 // Backend might not be up yet, so just catch the error
    //             }
    //             return false;
    //         };
    //
    //         let attempts = 0;
    //         const maxAttempts = 10;
    //         const delay = 1000; // 1 second
    //
    //         const pollBackend = async () => {
    //             while (attempts < maxAttempts) {
    //                 const success = await checkBackendStatus();
    //                 if (success) {
    //                     break;
    //                 }
    //                 attempts++;
    //                 await new Promise(res => setTimeout(res, delay));
    //             }
    //
    //             if (attempts >= maxAttempts) {
    //                 setBackendStatus('failed');
    //             }
    //         };
    //
    //         pollBackend();
    //     } catch (error) {
    //         console.error("Failed to start the backend:", error);
    //         setBackendStatus('failed');
    //     }
    // }


    return (
        <Router>
            <div className="container">

                <Navbar startbackend={startBackend} />

                <Routes>
                    <Route path="/" element={<Home/>}/>
                    <Route path="/localjobs" element={<LocalJobs/>}/>
                    <Route path="/cloudjobs" element={<CloudJobs/>}/>
                    <Route path="/results/:jobId" element={<Results/>}/>
                    <Route path="/settings" element={<Settings/>}/>
                    <Route path="/healthstatus" element={<Status />}/>
                </Routes>
            </div>
        </Router>
    );
}

export default App;
