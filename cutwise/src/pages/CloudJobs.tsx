import React, {FunctionComponent, useEffect} from "react";

export const CloudJobs: FunctionComponent = () => {
    useEffect(() => {
        console.log("Launched Cloud Jobs")
    }, []);

    return (
        <div>
            <h1>Cloud Jobs</h1>
        </div>
    )
}