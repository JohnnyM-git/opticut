import { useEffect, useState } from "react";

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
    const res = await fetch("http://localhost:2828/api/v1/settings");
    const data = await res.json();
    setSettings(data);
    console.log(data);
  };

  const setKerf = (value: string) => {
    setSettings({ ...settings, kerf: parseFloat(value) });
  };

  const saveSettings = async () => {
    try {
      const res = await fetch("http://localhost:2828/api/v1/settings", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(settings),
      });

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      // console.log(await res.json());
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
