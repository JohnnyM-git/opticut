import React, { useEffect, useState } from "react";
import axios from "axios";

export const Settings = () => {
    const [settings, setSettings] = useState({ kerf: 0.0625 });

    useEffect(() => {
        // Fetch settings from backend
        axios.get("/api/v1/settings")
            .then(response => setSettings(response.data))
            .catch(error => console.error("Error fetching settings:", error));

        console.log(settings);
    }, []);

    // const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    //     setSettings({ ...settings, [event.target.name]: parseFloat(event.target.value) });
    // };
    //
    // const handleSave = () => {
    //     // Save updated settings to backend
    //     axios.post("http://localhost:8080/api/v1/setsettings", settings)
    //         .then(() => console.log("Settings saved"))
    //         .catch(error => console.error("Error saving settings:", error));
    // };

    return (
        <div>
            <h1>Settings</h1>
            <label>
                Kerf:
                <input
                    type="number"
                    name="kerf"
                    value={settings.kerf}
                    // onChange={handleInputChange}
                />
            </label>
            {/*<button onClick={handleSave}>Save</button>*/}
        </div>
    );
};
