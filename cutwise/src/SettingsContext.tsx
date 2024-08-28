// SettingsContext.tsx
import React, { createContext, useState, useEffect, ReactNode } from "react";
import { Settings as SettingsType } from "./globals";

const defaultSettings: SettingsType = {
  kerf: 0,
  units: "imperial",
};

const SettingsContext = createContext<{
  settings: SettingsType;
  setSettings: React.Dispatch<React.SetStateAction<SettingsType>>;
}>({
  settings: defaultSettings,
  setSettings: () => {},
});

export const SettingsProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [settings, setSettings] = useState<SettingsType>(defaultSettings);

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const res = await fetch("http://localhost:2828/api/v1/settings");
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        const data: SettingsType = await res.json();
        setSettings(data);
      } catch (error) {
        console.error("Error fetching settings:", error);
      }
    };

    fetchSettings();
  }, []);

  return (
    <SettingsContext.Provider value={{ settings, setSettings }}>
      {children}
    </SettingsContext.Provider>
  );
};

export const useSettings = () => React.useContext(SettingsContext);
