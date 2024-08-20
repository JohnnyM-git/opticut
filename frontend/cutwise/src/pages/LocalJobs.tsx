import React, { FunctionComponent, useEffect, useState } from "react";
import { apiUrl } from "../globals.ts";
import "../styles/LocalJobs.css";
import { jobList } from "../globals.ts";
import { useNavigate } from "react-router-dom";

export const LocalJobs: FunctionComponent = () => {
  const [jobs, setJobs] = useState<jobList>({ JobList: [] });
  const navigate = useNavigate();
  useEffect(() => {
    console.log("Launched Local jobs");
    getLocalJobs();
  }, []);

  async function getLocalJobs() {
    const res = await fetch(`${apiUrl}localjobs`);
    const data = await res.json();
    console.log(data);
    setJobs(data);
  }

  return (
    <div>
      <h1 className={"heading"}>Local Jobs</h1>
      {jobs?.JobsList?.map((job, i) => (
        <div className={"job"} key={i}>
          <h2
            className={"job__number"}
            onClick={() => navigate(`/results/${job.JobNumber}`)}
          >
            Job: {job.JobNumber}
          </h2>
          <p className={"customer"}>Customer: {job.Customer}</p>
        </div>
      ))}
    </div>
  );
};
