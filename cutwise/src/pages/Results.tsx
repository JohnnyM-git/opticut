import { FunctionComponent, useEffect, useState } from "react";
import { apiUrl } from "../globals.ts";
import { useParams } from "react-router-dom";
import styles from "../styles/Results.module.css";
import { Button } from "@mui/material";
import { Print, Star, StarBorder } from "@mui/icons-material";
import { useSettings } from "../SettingsContext.tsx";
import { inToFt, inToM, inToMm } from "../functions/unitConverter.ts";

export const Results: FunctionComponent = () => {
  interface JobDetails {
    Job: string;
    Customer: string;
    Star: number;
  }

  interface Job {
    message: string;
    job_data_materials: Array<any>;
    job_data_parts: Array<any>;
    material_data: Array<any>;
    job_info: JobDetails;
  }

  const [job, setJob] = useState<Job>({
    message: "",
    job_data_materials: [],
    job_data_parts: [],
    material_data: [],
    job_info: {
      Job: "",
      Customer: "",
      Star: 0,
    },
  });
  const { jobId } = useParams();
  const { settings, setSettings } = useSettings();
  // const [printable, setPrintable] = useState(true)

  useEffect(() => {
    getJob();
  }, [jobId]);

  async function getJob() {
    console.log(jobId);
    // const jobNum = query
    const res = await fetch(`${apiUrl}job?job_id=${jobId}`);
    if (!res.ok) {
      // console.error(res);
      throw new Error(`HTTP error! status: ${res.status}`);
    }
    const data = await res.json();
    console.log(data);
    setJob(data);
  }

  async function toggleStar(jobNumber: string, value: number) {
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
        ...prevJob.job_info,
        Star: newStarValue,
      },
    }));
  };

  return (
    <div>
      <h1 className={styles.heading}>Results</h1>

      {/*Job Info*/}
      <div className={styles.job}>
        <div className={styles.job__info}>
          <h2>Job Info</h2>
          <h3>Job Number: {job.job_info.Job}</h3>
          <h4>Customer: {job.job_info.Customer}</h4>
        </div>
        <div className={styles.buttons}>
          <Button
            onClick={() =>
              toggleStar(job.job_info.Job, job.job_info.Star === 0 ? 1 : 0)
            }
          >
            {job.job_info.Star === 0 ? <StarBorder /> : <Star />}
          </Button>
        </div>
      </div>
      {/*Job Info*/}

      {/*Visual display of materials totals*/}
      <div className={styles.material}>
        <div className={styles.material__info}>
          <div className={styles.material__info__left}>
            <h2>Material Info</h2>
            <h3>Job Number: {job.job_info.Job}</h3>
            {/*<h4>Customer: {job.Job.Customer}</h4>*/}
          </div>
          <div className={styles.material__info__right}>
            <Button onClick={() => window.print()}>
              <Print />
              Print
            </Button>
          </div>
        </div>
        <div className={styles.material__totals}>
          {job.material_data
            .sort((a, b) => a.material_code.localeCompare(b.material_code))
            .map((item) => (
              <div key={item.id}>
                <h3 className={styles.material__code}>
                  Material Code: {item.material_code}{" "}
                  <span className={styles.tip}>(click to filter results)</span>{" "}
                </h3>
                {settings.units === "imperial" && (
                  <p className={styles.material__info__property}>
                    Total Length: {inToFt(item.total_stock_length).toFixed(2)}'
                    | {item.total_stock_length.toFixed(2)}"
                  </p>
                )}

                {settings.units === "metric" && (
                  <p className={styles.material__info__property}>
                    Total Length: {inToM(item.total_stock_length).toFixed(2)}M |{" "}
                    {inToMm(item.total_stock_length).toFixed(2)}mm
                  </p>
                )}

                {settings.units === "imperial" && (
                  <p className={styles.material__info__property}>
                    Total Used: {inToFt(item.total_used_length).toFixed(2)}' |{" "}
                    {item.total_used_length.toFixed(2)}"
                  </p>
                )}

                {settings.units === "metric" && (
                  <p className={styles.material__info__property}>
                    Total Used: {inToM(item.total_used_length).toFixed(2)}M |{" "}
                    {inToMm(item.total_used_length).toFixed(2)}mm
                  </p>
                )}

                {/*<p className={styles.material__info__property}>*/}
                {/*  Total Used: {(item.total_used_length / 12).toFixed(2)}' |{" "}*/}
                {/*  {item.total_used_length.toFixed(2)}"*/}
                {/*</p>*/}

                {settings.units === "imperial" && (
                  <p className={styles.material__info__property}>
                    QTY: {item.total_quantity} | Stock Length:{" "}
                    {inToFt(item.stock_length).toFixed(2)}' |{" "}
                    {item.stock_length.toFixed(2)}"
                  </p>
                )}

                {settings.units === "metric" && (
                  <p className={styles.material__info__property}>
                    QTY: {item.total_quantity} | Stock Length:{" "}
                    {inToM(item.stock_length).toFixed(2)}M |{" "}
                    {inToMm(item.stock_length).toFixed(2)}mm
                  </p>
                )}

                {/*<p className={styles.material__info__property}>*/}
                {/*  QTY: {item.total_quantity} | Stock Length:{" "}*/}
                {/*  {(item.stock_length / 12).toFixed(2)}' |{" "}*/}
                {/*  {item.stock_length.toFixed(2)}"*/}
                {/*</p>*/}
              </div>
            ))}
        </div>
      </div>
      {/*Visual display of materials totals*/}

      {/*Visual display of materials*/}
      {job.job_data_materials.map((item) => (
        <div className={styles.material__display} key={item.cut_material_id}>
          <div className={styles.material__display__info}>
            <div className={styles.material__display__info__left}>
              <h4>Material Code: {item.cut_material_material_code}</h4>
              <h4>Material ID: {item.cut_material_id}</h4>

              {settings.units === "imperial" && (
                <h4>
                  Total Material Used:{" "}
                  {inToFt(item.total_used_length).toFixed(2)}' |{" "}
                  {item.total_used_length.toFixed(2)}"
                </h4>
              )}

              {settings.units === "metric" && (
                <h4>
                  Total Material Used:{" "}
                  {inToM(item.total_used_length).toFixed(2)}M |{" "}
                  {inToMm(item.total_used_length).toFixed(2)}mm
                </h4>
              )}

              {/*<h4>*/}
              {/*  Total Material Used: {(item.total_used_length / 12).toFixed(2)}'*/}
              {/*  | {item.total_used_length.toFixed(2)}"*/}
              {/*</h4>*/}

              {settings.units === "imperial" && (
                <h4>
                  Material Stock Length: {inToFt(item.stock_length).toFixed(2)}'
                  | {item.stock_length.toFixed(2)}
                </h4>
              )}

              {settings.units === "metric" && (
                <h4>
                  Material Stock Length: {inToM(item.stock_length).toFixed(2)}M
                  | {inToMm(item.stock_length).toFixed(2)}mm
                </h4>
              )}

              {/*<h4>*/}
              {/*  Material Stock Length: {(item.stock_length / 12).toFixed(2)}' |{" "}*/}
              {/*  {item.stock_length.toFixed(2)}*/}
              {/*</h4>*/}

              {settings.units === "imperial" && (
                <h4>
                  Material Scrap / Drop:{" "}
                  {inToFt(item.stock_length - item.total_used_length).toFixed(
                    2,
                  )}
                  ' | {(item.stock_length - item.total_used_length).toFixed(2)}"
                </h4>
              )}

              {settings.units === "metric" && (
                <h4>
                  Material Scrap / Drop:{" "}
                  {inToM(item.stock_length - item.total_used_length).toFixed(2)}
                  M |{" "}
                  {inToMm(item.stock_length - item.total_used_length).toFixed(
                    2,
                  )}
                  mm
                </h4>
              )}

              {/*<h4>*/}
              {/*  Material Scrap / Drop:{" "}*/}
              {/*  {((item.stock_length - item.total_used_length) / 12).toFixed(2)}*/}
              {/*  ' | {(item.stock_length - item.total_used_length).toFixed(2)}"*/}
              {/*</h4>*/}
              <h4>Identical Lengths: {item.cut_material_quantity}</h4>
              <h4>Unique Parts: {item.total_parts_cut_on_material}</h4>
            </div>
          </div>
          <div className={styles.material__display__material__section}>
            <div className={styles.material__display__material__visual}>
              {job.job_data_parts
                .filter(
                  (partdata) =>
                    partdata.cut_material_id === item.cut_material_id,
                )
                .map((partdata) =>
                  Array.from({ length: partdata.part_qty }).map((_, index) => (
                    <div
                      key={`${partdata.part_id}-${index}`}
                      className={
                        styles.material__display__material__visual__part
                      }
                      style={{
                        width: `${(partdata.part_cut_length / item.stock_length) * 100}%`,
                      }}
                    >
                      <p>{partdata.part_id}</p>
                    </div>
                  )),
                )}
            </div>
            <div
              className={
                styles.material__display__material__section__part__list
              }
            >
              <div
                className={
                  styles.material__display__material__section__part__list__properties
                }
              >
                <h5
                  className={
                    styles.material__display__material__section__part__list__heading
                  }
                >
                  ID
                </h5>
                {job.job_data_parts
                  .filter(
                    (partdata) =>
                      partdata.cut_material_id === item.cut_material_id,
                  )
                  .map((partdata, i) => (
                    <p key={i}>{partdata.part_id}</p>
                  ))}
              </div>

              <div
                className={
                  styles.material__display__material__section__part__list__properties
                }
              >
                <h5
                  className={
                    styles.material__display__material__section__part__list__heading
                  }
                >
                  Part Number
                </h5>
                {job.job_data_parts
                  .filter(
                    (partdata) =>
                      partdata.cut_material_id === item.cut_material_id,
                  )
                  .map((partdata, i) => (
                    <p key={i}>{partdata.part_number}</p>
                  ))}
              </div>

              <div
                className={
                  styles.material__display__material__section__part__list__properties
                }
              >
                <h5
                  className={
                    styles.material__display__material__section__part__list__heading
                  }
                >
                  Part Length
                </h5>
                {job.job_data_parts
                  .filter(
                    (partdata) =>
                      partdata.cut_material_id === item.cut_material_id,
                  )
                  .map((partdata, i) => (
                    <div key={i}>
                      {settings.units === "imperial" ? (
                        <p>{partdata.part_length.toFixed(2)}"</p>
                      ) : settings.units === "metric" ? (
                        <p>{inToMm(partdata.part_length).toFixed(2)}mm</p>
                      ) : null}
                    </div>
                  ))}
              </div>

              <div
                className={
                  styles.material__display__material__section__part__list__properties
                }
              >
                <h5
                  className={
                    styles.material__display__material__section__part__list__heading
                  }
                >
                  Current Qty
                </h5>
                {job.job_data_parts
                  .filter(
                    (partdata) =>
                      partdata.cut_material_id === item.cut_material_id,
                  )
                  .map((partdata, i) => (
                    <p key={i}>{partdata.part_qty}</p>
                  ))}
              </div>

              <div
                className={
                  styles.material__display__material__section__part__list__properties
                }
              >
                <h5
                  className={
                    styles.material__display__material__section__part__list__heading
                  }
                >
                  Total Qty
                </h5>
                {job.job_data_parts
                  .filter(
                    (partdata) =>
                      partdata.cut_material_id === item.cut_material_id,
                  )
                  .map((partdata, i) => (
                    <p key={i}>{partdata.total_part_qty}</p>
                  ))}
              </div>
            </div>

            {/*MATERIAL DISPLAY*/}
          </div>
        </div>
      ))}

      {/*Visual display of materials*/}
    </div>
  );
};
