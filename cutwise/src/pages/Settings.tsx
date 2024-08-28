// Settings.tsx
import { useEffect } from "react";
import styles from "../styles/Settings.module.css";
import { StyledInput } from "../components/StyledInput";
import { useSettings } from "../SettingsContext";

export const Settings = () => {
  const { settings, setSettings } = useSettings();

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const res = await fetch("http://localhost:2828/api/v1/settings");
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        const data = await res.json();
        setSettings(data);
      } catch (error) {
        console.error("Error fetching settings:", error);
      }
    };

    fetchSettings();
  }, [setSettings]);

  const setKerf = (value: string) => {
    setSettings((prevSettings) => ({
      ...prevSettings,
      kerf: parseFloat(value),
    }));
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
      const data = await res.json();
      console.log("Settings saved:", data);
    } catch (error) {
      console.error("Error saving settings:", error);
    }
  };

  const handleUnitChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSettings((prevSettings) => ({
      ...prevSettings,
      units: e.target.value,
    }));
  };

  return (
    <div className={styles.settings}>
      <h1>Settings</h1>
      <div className={styles.input}>
        <StyledInput
          type="number"
          value={
            settings.units === "imperial"
              ? settings.kerf.toFixed(4)
              : (settings.kerf * 25.4).toFixed(4)
          }
          placeholder="Kerf"
          onChange={(e) => setKerf(e.target.value)}
          step="0.01"
        />

        <button onClick={saveSettings}>Save</button>
      </div>
      <fieldset>
        <label>
          <StyledInput
            type="radio"
            value="imperial"
            name="units"
            checked={settings.units === "imperial"}
            onChange={handleUnitChange}
          />
          Imperial - Ft/In
        </label>
        <label>
          <StyledInput
            type="radio"
            value="metric"
            name="units"
            checked={settings.units === "metric"}
            onChange={handleUnitChange}
          />
          Metric - M/mm
        </label>
      </fieldset>
    </div>
  );
};
