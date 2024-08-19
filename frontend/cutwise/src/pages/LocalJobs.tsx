import React, {FunctionComponent, useEffect, useState} from "react";
import {apiUrl} from "../globals.ts";
import "../styles/LocalJobs.css"

export const LocalJobs: FunctionComponent = () => {
    const [jobs, setJobs] = useState()
    useEffect(() => {
        console.log("Launched Local jobs")
        getLocalJobs()

    }, []);

async function getLocalJobs() {
        const res = await fetch(`${apiUrl}localjobs`)
        const data = await res.json()
        console.log(data)
        setJobs(data)


    }

    return (
        <div>
            <h1>Local Jobs</h1>
            {jobs?.JobsList?.map((job, i) => (
                <div className={"job"} key={i}>
                    <h2 className={"job__Data"}>Job Number: {job.JobNumber}</h2>
                    <h2 className={"job__Data"}>Customer: {job.Customer}</h2>
                </div>
            ))}
        </div>
    )

}