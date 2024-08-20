import React, { useEffect, useState } from "react";
import axios from "axios";

export const Settings = () => {
  const [settings, setSettings] = useState({
    kerf: 0.0625,
  });

  useEffect(() => {
    fetchSettings();
    // Fetch settings from backend
    // axios
    //   .get("/api/v1/settings")
    //   .then((response) => setSettings(response.data))
    //   .catch((error) => console.error("Error fetching settings:", error));

    // console.log(settings);
  }, []);

  const fetchSettings = async () => {
    const res = await fetch("http://localhost:8080/api/v1/settings");
    const data = await res.json();
    setSettings(data);
    console.log(data);
  };

  const setKerf = (value: number) => {
    setSettings({ ...settings, kerf: value });
  };

  const saveSettings = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/v1/settings", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(settings), // assuming 'settings' is the object you want to save
      });

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }

      const data = await res.json();
      console.log("Settings saved:", data);
    } catch (error) {
      console.error("Error saving settings:", error);
    }
  };

  return (
    <div>
      <h1>Settings</h1>
      <label>
        Kerf:
        <input
          type="number"
          name="kerf"
          value={settings.kerf}
          onChange={(e) => setKerf(e.target.value)}
        />
      </label>
      <button onClick={saveSettings}>Save</button>
    </div>
  );
};
