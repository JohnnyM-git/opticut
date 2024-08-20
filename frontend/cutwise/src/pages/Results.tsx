import {useEffect, useState} from "react";
import {apiUrl} from "../globals.ts";
import {useParams} from "react-router-dom";
import styles from "../styles/Results.module.css";
import {Button} from "@mui/material";
import {Star, StarBorder} from "@mui/icons-material";



export const Results: Function = (() => {

    interface JobDetails {
        Job: string;
        Customer: string;
        Star: number;
    }

    interface Job {
        message: string;
        JobData: Array<any>;
        MaterialData: Array<any>;
        Job: JobDetails;
    }

    const [job, setJob] = useState<Job>({
        message: '',
        JobData: [],
        MaterialData: [],
        Job: {
            Job: '',
            Customer: '',
            Star: 0,
        },
    });
    const { jobId } = useParams();

    useEffect(() => {
        getJob()
    }, [jobId])

    async function getJob() {
        console.log(jobId);
        // const jobNum = query
        const res = await fetch(`${apiUrl}job?job_id=${jobId}`)
        if (!res.ok) {
            // console.error(res);
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        const data = await res.json();
        console.log(data);
        setJob(data);


    }

    async function toggleStar(jobNumber: string, value: number,) {
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
            updateStar(value);
        } catch (error) {
            console.error("Error toggling star:", error);
        }
    }

    const updateStar = (newStarValue: number) => {
        setJob((prevJob) => ({
            ...prevJob,
            Job: {
                ...prevJob.Job,
                Star: newStarValue,
            },
        }));
    };



    return (
        <div>
            <h1 className={styles.heading}>Results</h1>
            <div className={styles.job}>
                <div className={styles.job__info}>
                    <h2>Job Info</h2>
                    <h3>Job Number: {job.Job.Job}</h3>
                    <h4>Customer: {job.Job.Customer}</h4>
                </div>
                <div className={styles.buttons}>
                    <Button
                        onClick={() => toggleStar(job.Job.Job, job.Job.Star === 0 ? 1 : 0)}>
                        {job.Job.Star === 0 ? <StarBorder/> : <Star/>}
                    </Button>
                </div>
            </div>

            <div className={styles.material}>
                <div className={styles.material__info}>
                    <h2>Material Info</h2>
                    <h3>Job Number: {job.Job.Job}</h3>
                    {/*<h4>Customer: {job.Job.Customer}</h4>*/}
                </div>
                <div className={styles.material__totals}>
                    {job.MaterialData
                        .sort((a, b) => a.MaterialCode.localeCompare(b.MaterialCode))
                        .map((item, i) => (
                            <div key={i}>
                                <h3 className={styles.material__code}>Material Code: {item.MaterialCode} <span className={styles.tip}>(click to filter results)</span> </h3>
                                <p className={styles.material__info__property}>Total Length: {(item.TotalLength / 12).toFixed(2)}'</p>
                                <p className={styles.material__info__property}>QTY: {item.TotalQuantity} | Stock Length: {(item.StockLength / 12).toFixed(2)}' | {item.StockLength}"</p>

                            </div>
                        ))}

                </div>
            </div>


        </div>
    )
})