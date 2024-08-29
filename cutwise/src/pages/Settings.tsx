// Settings.tsx
import { useEffect, useState } from "react";
import styles from "../styles/Settings.module.css";
import { StyledInput } from "../components/StyledInput";
import { useSettings } from "../SettingsContext";
import { Settings } from "../globals.ts";

export const SettingsPage = () => {
  const { settings, setSettings } = useSettings();
  const [settingsChanged, setSettingsChanged] = useState(false);

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

  // const setKerf = (value: string) => {
  //   setSettings((prevSettings) => ({
  //     ...prevSettings,
  //     kerf: parseFloat(value),
  //   }));
  //   setSettingsChanged(true);
  // };

  function updateSettings(key: keyof Settings, value: string | number): void {
    setSettingsChanged(true);
    console.log(settingsChanged);
    if (key === "kerf") {
      if (typeof value === "string") {
        const parsedValue = parseFloat(value);
        if (isNaN(parsedValue)) {
          console.error("Invalid value for kerf. Must be a number.");
          return; // Exit early if parsing fails
        }
        value = parsedValue;
      }
    }
    setSettings((prev) => ({
      ...prev,
      [key]: value,
    }));
  }

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
      setSettingsChanged(false);
    } catch (error) {
      console.error("Error saving settings:", error);
    }
  };

  // const handleUnitChange = (e: React.ChangeEvent<HTMLInputElement>) => {
  //   setSettings((prevSettings) => ({
  //     ...prevSettings,
  //     units: e.target.value,
  //   }));
  //   setSettingsChanged(true);
  // };

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
          onChange={(e) => updateSettings("kerf", e.target.value)}
          step="0.01"
        />
      </div>
      <div className={styles.unit__selection}>
        <p>Select Units</p>
        <fieldset className={styles.unit__selection__options}>
          <label>
            <StyledInput
              type="radio"
              value="imperial"
              name="units"
              checked={settings.units === "imperial"}
              onChange={(e) => updateSettings("units", e.target.value)}
            />
            Imperial - Ft/In
          </label>
          <label>
            <StyledInput
              type="radio"
              value="metric"
              name="units"
              checked={settings.units === "metric"}
              onChange={(e) => updateSettings("units", e.target.value)}
            />
            Metric - M/mm
          </label>
        </fieldset>
      </div>
      <button
        className={styles.button}
        onClick={saveSettings}
        disabled={!settingsChanged}
      >
        Save
      </button>
    </div>
  );
};
