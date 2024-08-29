import { FunctionComponent, useEffect, useState } from "react";
import { apiUrl } from "../globals.ts";
import styles from "./styles/footer.module.css";
import { invoke } from "@tauri-apps/api/tauri";

export const Footer: FunctionComponent = () => {
  const [status, setStatus] = useState<string>("Offline");

  useEffect(() => {
    // Function to fetch online status
    const getOnlineStatus = async () => {
      try {
        const res = await fetch(`${apiUrl}/health`);
        if (!res.ok) {
          setStatus("Offline");
          throw new Error("Network response was not ok");
        }
        const data = await res.json();
        setStatus(data.status);
        logToFile(`Server status is ${data.status}`);
      } catch (err) {
        setStatus("Offline");
        logToFile("Failed to fetch status");
      }
    };

    // Initial fetch
    getOnlineStatus();

    // Set up interval to fetch status every 5 minutes
    const intervalId = setInterval(getOnlineStatus, 300000);

    // Clean up interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  async function logToFile(message: string) {
    try {
      await invoke("log_to_file", { message });
      console.log("Log message written successfully.");
    } catch (error) {
      console.error("Failed to write log:", error);
    }
  }

  return (
    <div className={styles.status}>
      <div
        className={
          status === "healthy"
            ? styles.status__dot__online
            : styles.status__dot__offline
        }
      ></div>
      <span className={styles.status__text}>
        {status === "healthy" ? "Online" : "Offline"}
      </span>
    </div>
  );
};
