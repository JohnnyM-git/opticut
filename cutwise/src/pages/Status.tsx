import styles from "../styles/Status.module.css";
import { useEffect, useState } from "react";
import { apiUrl } from "../globals.ts";
import { invoke } from "@tauri-apps/api/tauri";

interface healthStatus {
  database: string;
  service_name: string;
  status: string;
  uptime: string;
  version: string;
}

export const Status = () => {
  const [health, setHealth] = useState<healthStatus>();
  const offlineState = {
    database: "Offline",
    service_name: "Cutwise",
    status: "Offline",
    uptime: "Offline",
    version: "Offline",
  };

  useEffect(() => {
    async function checkStatus() {
      const res = await fetch(`${apiUrl}/health`);
      let data = await res.json();
      if (data) {
        data.uptime = formatDuration(data.uptime);
        console.log(data);

        setHealth(data);
      } else {
        setHealth(offlineState);
      }
    }

    checkStatus();
  }, []);

  function formatDuration(duration: string): string {
    // Regular expression to match minutes and seconds with optional fractional seconds
    const regex = /(?:(\d+)m)?(\d+)(\.\d+)?s/;
    const match = duration.match(regex);

    if (!match) {
      throw new Error("Invalid duration format");
    }

    // Extract the minutes, seconds, and fractional part (if any)
    const minutes = match[1] ? parseInt(match[1], 10) : 0;
    const seconds = parseFloat(match[2]);

    // Format the duration
    if (minutes > 0) {
      return `${minutes}m ${seconds}s`;
    } else {
      return `${seconds}s`;
    }
  }

  async function restartService() {
    try {
      await logToFile("Starting server...");

      // Start the backend asynchronously
      await invoke("start_backend");
      await logToFile("Server started.");

      // Poll the backend health endpoint until it's up or timeout
      const checkBackendStatus = async () => {
        try {
          const res = await fetch(`${apiUrl}/health`); // Replace with actual health check endpoint
          if (res.ok) {
            setHealth(await res.json());
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

  return (
    <>
      {!health ? (
        <div>Please wait...</div>
      ) : (
        <div className={styles.healthDashboard}>
          <h2>{health.service_name} Health Status</h2>
          <div className={styles.status}>
            <span>Status:</span>
            <span
              className={`${styles.statusIndicator} ${
                health.status === "healthy" ? styles.healthy : styles.offline
              }`}
            >
              {health.status}
            </span>
          </div>
          <div className={styles.details}>
            <p>
              <strong>Database:</strong> {health.database}
            </p>
            <p>
              <strong>Version:</strong> {health.version}
            </p>
            <p>
              <strong>Uptime:</strong> {health.uptime}
            </p>
          </div>
          {health.status === "Offline" && (
            <>
              <h3>Service is offline, please try restarting</h3>
              <h4 onClick={restartService}>Restart Service</h4>
            </>
          )}
          <h4 onClick={() => window.location.reload()}>Refresh</h4>
        </div>
      )}
    </>
  );
};
