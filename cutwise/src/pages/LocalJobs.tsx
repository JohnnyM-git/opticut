import { FunctionComponent, useEffect, useState } from "react";
import { apiUrl } from "../globals.ts";
import styles from "../styles/LocalJobs.module.css";
import { jobListJob } from "../globals.ts"; // Update import if needed
import { useNavigate } from "react-router-dom";
import { Button } from "@mui/material";
import { Launch, Star, StarBorder } from "@mui/icons-material";

// Update interface to match API response
interface JobsResponse {
  Message: string;
  JobsList: jobListJob[];
}

export const LocalJobs: FunctionComponent = () => {
  const [jobs, setJobs] = useState<jobListJob[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    console.log("Launched Local Jobs");
    getLocalJobs();
  }, []);

  async function getLocalJobs() {
    try {
      const res = await fetch(`${apiUrl}localjobs`);
      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      const data: JobsResponse = await res.json();
      console.log(data);
      setJobs(data.JobsList); // Update state with JobsList
    } catch (error) {
      console.error("Error fetching local jobs:", error);
    }
  }

  async function toggleStar(jobNumber: string, value: number, i: number) {
    const options = {
      jobNumber: jobNumber.toString(),
      value: value,
    };

    try {
      const res = await fetch(`${apiUrl}togglestar`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(options),
      });

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      const data = await res.json();
      console.log(data);
      updateStar(i, value);
    } catch (error) {
      console.error("Error toggling star:", error);
    }
  }

  const updateStar = (index: number, newStarValue: number) => {
    setJobs(prevJobs =>
        prevJobs.map((job, i) =>
            i === index ? { ...job, Star: newStarValue } : job
        )
    );
  };

  return (
      <div>
        <h1 className={styles.heading}>Local Jobs</h1>
        <div className={styles.jobs}>
          {jobs?.map((job, i) => (
              <div className={styles.job} key={i}>
                <h2
                    className={styles.job__number}
                    onClick={() => navigate(`/results/${job.JobNumber}`)}
                >
                  Job: {job.JobNumber}
                </h2>
                <p className={styles.customer}>Customer: {job.Customer}</p>
                <div className={styles.buttons}>
                  <Button>
                    <Launch onClick={() => navigate(`/results/${job.JobNumber}`)} />
                  </Button>
                  <Button onClick={() => toggleStar(job.JobNumber, job.Star === 0 ? 1 : 0, i)}>
                    {job.Star === 0 ? <StarBorder /> : <Star />}
                  </Button>
                </div>
              </div>
          ))}
        </div>
      </div>
  );
};
