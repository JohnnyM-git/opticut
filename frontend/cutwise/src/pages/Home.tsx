import React, { FunctionComponent, useState } from "react";
import styles from "../styles/Home.module.css";
import { Badge, Button, List, ListItemIcon, TextField } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import { Add, Segment } from "@mui/icons-material";

export const Home: FunctionComponent = () => {
  interface JobInfo {
    Job: string;
    Customer: string;
  }
  interface Part {
    PartNumber: string;
    MaterialCode: string;
    Length?: number | null;
    Quantity?: number | null;
    CuttingOperation: string;
  }

  interface Material {
    MaterialCode: string;
    Length: number;
    Quantity: number;
  }

  const [jobInfo, setJobInfo] = useState<JobInfo>({ Job: "", Customer: "" });
  const [parts, setParts] = useState<Part[]>([
    {
      PartNumber: "test",
      MaterialCode: "test",
      Length: 10,
      CuttingOperation: "",
      Quantity: 1,
    },
  ]);
  const [material, setMaterial] = useState<Material[]>([
    { MaterialCode: "test", Length: 1, Quantity: 10 },
  ]);
  const [part, setPart] = useState<Part>({
    PartNumber: "",
    MaterialCode: "",
    Length: null,
    Quantity: null,
    CuttingOperation: "",
  });
  return (
    <div>
      <h1>Create New Job</h1>
      <div className={styles.heading}>
        <div className={styles.heading__left}>
          <TextField
            placeholder={"Enter Job / Project #"}
            value={jobInfo?.Job}
            onChange={(e) =>
              setJobInfo((prev) => ({
                ...prev,
                Job: e.target.value.toUpperCase(),
              }))
            }
          />
          <TextField
            placeholder={"Enter Customer"}
            value={jobInfo?.Customer}
            onChange={(e) =>
              setJobInfo((prev) => ({
                ...prev,
                Customer: e.target.value.toUpperCase(),
              }))
            }
          />
        </div>
        <div className={styles.heading__right}>
          <Badge badgeContent={parts.length} color="primary">
            <Button>
              <p>Parts</p>
              <Segment />
            </Button>
          </Badge>
          <Badge badgeContent={material.length} color="primary">
            <Button>
              <p>Materials</p>
              <Segment />
            </Button>
          </Badge>
        </div>
      </div>
      <div className={styles.item__additions}>
        <div className={styles.part__addition}>
          <h4>Add Part to Cut</h4>
          <div className={styles.part__addition__inputs}>
            <TextField
              placeholder={"Part Number"}
              value={part?.PartNumber}
              type={"text"}
              onChange={(e) =>
                setPart((prev) => ({
                  ...prev,
                  PartNumber: e.target.value.toUpperCase(),
                }))
              }
            />
            <TextField
              placeholder={"Material Code"}
              value={part?.MaterialCode}
              type={"text"}
              onChange={(e) =>
                setPart((prev) => ({
                  ...prev,
                  MaterialCode: e.target.value.toUpperCase(),
                }))
              }
            />
            <TextField
              placeholder={"Part Length - inch"}
              value={part?.Length}
              type={"number"}
              onChange={(e) =>
                setPart((prev) => ({
                  ...prev,
                  Length: parseInt(e.target.value),
                }))
              }
            />
            <TextField
              placeholder={"Cutting Operation"}
              value={part?.CuttingOperation}
              type={"text"}
              onChange={(e) =>
                setPart((prev) => ({
                  ...prev,
                  CuttingOperation: e.target.value.toUpperCase(),
                }))
              }
            />
            <TextField
              placeholder={"Part Quantity"}
              onChange={(e) =>
                setPart((prev) => ({
                  ...prev,
                  Quantity: parseInt(e.target.value),
                }))
              }
            />
            <Button>
              <Add onClick={() => setParts((prev) => [...prev, part])} />
            </Button>
          </div>
        </div>

        <div className={styles.material__addition}>
          <h4>Add Material</h4>
          <div className={styles.material__addition__inputs}>
            <TextField placeholder={"Material Code"} />
            <TextField placeholder={"Material Length - inch"} />
            <TextField placeholder={"Material Quantity"} />

            <Button>
              <Add />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
