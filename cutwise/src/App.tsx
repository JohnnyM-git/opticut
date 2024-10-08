import "./App.css";
import { Navbar } from "./components/navbar.tsx";
import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { SettingsPage } from "./pages/Settings.tsx";
import { Home } from "./pages/Home.tsx";
import { LocalJobs } from "./pages/LocalJobs.tsx";
import { CloudJobs } from "./pages/CloudJobs.tsx";
import { Results } from "./pages/Results.tsx";
import { invoke } from "@tauri-apps/api/tauri";
import { apiUrl } from "./globals.ts";
import { Status } from "./pages/Status.tsx";
import { Footer } from "./components/footer.tsx";
import { SettingsProvider } from "./SettingsContext.tsx";
import {
  checkUpdate,
  installUpdate,
  onUpdaterEvent,
} from "@tauri-apps/api/updater";
import { relaunch } from "@tauri-apps/api/process";

function App() {
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Function to toggle the modalOpen state
  const toggleModal = () => {
    setIsModalOpen((prevState) => !prevState);
  };

  useEffect(() => {
    startBackend();
  }, []);

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
          await new Promise((res) => setTimeout(res, delay));
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
      await invoke("log_to_file", { message });
      console.log("Log message written successfully.");
    } catch (error) {
      console.error("Failed to write log:", error);
    }
  }

  // async function pingBackend() {
  //   const res = await fetch(`${apiUrl}/health`);
  //   const data = res.json();
  //   console.log(data);
  // }

  async function checkForUpdates() {
    try {
      const { shouldUpdate, manifest } = await checkUpdate();

      if (shouldUpdate) {
        console.log(`Update available: ${manifest?.version}`);

        if (
          window.confirm(
            `Update available to version ${manifest?.version}. Do you want to update now?`,
          )
        ) {
          await installUpdate();
          console.log("Update installed. Restarting application...");
          await relaunch(); // This restarts the app after the update
        }
      } else {
        console.log("No updates available.");
      }
    } catch (error) {
      console.error("Error checking for updates:", error);
    }
  }

  onUpdaterEvent(({ error, status }) => {
    if (error) {
      console.error("Update error:", error);
    } else {
      console.log("Update status:", status);
    }
  });

  return (
    <SettingsProvider>
      <Router>
        {/* Main content */}
        <div className={`container ${isModalOpen ? "modalOpen" : ""}`}>
          <Navbar
            // startbackend={pingBackend}
            checkForUpdates={checkForUpdates}
          />

          <div className="container__secondary">
            <Routes>
              <Route path="/" element={<Home toggleModal={toggleModal} />} />
              <Route path="/localjobs" element={<LocalJobs />} />
              <Route path="/cloudjobs" element={<CloudJobs />} />
              <Route path="/results/:jobId" element={<Results />} />
              <Route path="/settings" element={<SettingsPage />} />
              <Route path="/healthstatus" element={<Status />} />
            </Routes>

            <div className="footer">
              <Footer />
            </div>
          </div>
        </div>
      </Router>
    </SettingsProvider>
  );
}

export default App;
