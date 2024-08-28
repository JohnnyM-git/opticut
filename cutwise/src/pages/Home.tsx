import { FunctionComponent, useState } from "react";
import styles from "../styles/Home.module.css";
import { Badge, Button } from "@mui/material";
import { Add, Segment } from "@mui/icons-material";
import { StyledInput } from "../components/StyledInput.tsx";
import { apiUrl } from "../globals.ts";
import { useNavigate } from "react-router-dom";

export const Home: FunctionComponent = () => {
  interface JobInfo {
    Job: string;
    Customer: string;
  }
  interface Part {
    PartNumber: string;
    MaterialCode: string;
    Length: number;
    Quantity: number;
    CuttingOperation: string;
  }

  const partInitialState = {
    PartNumber: "",
    MaterialCode: "",
    Length: 0,
    Quantity: 0,
    CuttingOperation: "",
  };

  interface Material {
    MaterialCode: string;
    Length: number;
    Quantity: number;
  }

  const materialInitialState = {
    MaterialCode: "",
    Length: 0,
    Quantity: 0,
  };

  const [jobInfo, setJobInfo] = useState<JobInfo>({ Job: "", Customer: "" });
  const [parts, setParts] = useState<Part[]>([]);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [part, setPart] = useState<Part>(partInitialState);
  const [material, setMaterial] = useState<Material>(materialInitialState);
  const [partErrorMsg, setPartErrorMsg] = useState<string | JSX.Element>("");
  const [materialErrorMsg, setMaterialErrorMsg] = useState<string>("");
  const [materialQtyDisabled, setMaterialQtyDisabled] = useState(false);
  const navigate = useNavigate();

  async function runProject(): Promise<void> {
    if (checkDataValidity()) {
      const res = await fetch(`${apiUrl}runproject`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          jobInfo: jobInfo,
          parts: parts,
          materials: materials,
        }),
      });
      const data = await res.json();
      console.log(data);
      navigate(`/results/${jobInfo.Job}`);
    }
  }

  function checkDataValidity(): boolean {
    if (
      !jobInfo.Job ||
      !jobInfo.Customer ||
      !parts.length ||
      !materials.length
    ) {
      console.log("Job Info Not Valid");
      return false;
    } else {
      console.log("Running Project");
      return true;
    }
  }

  function updatePart(key: string, value: string): void {
    let processedValue: string | number;

    if (key === "Length") {
      // Convert the value to a float, or default to 0 if invalid or empty
      processedValue = value ? parseFloat(value) || 0 : 0;
    } else if (key === "Quantity") {
      // Convert the value to an integer, or default to 0 if invalid or empty
      processedValue = value ? parseInt(value, 10) || 0 : 0;
    } else if (
      key === "PartNumber" ||
      key === "MaterialCode" ||
      key === "CuttingOperation"
    ) {
      // For string fields, trim and convert to uppercase
      processedValue = value.toUpperCase().trim();
    } else {
      // Handle unexpected keys if necessary
      console.warn(`Unexpected key: ${key}`);
      return;
    }

    setPart((prev) => ({
      ...prev,
      [key]: processedValue,
    }));
  }

  function checkPartValidity(): boolean {
    const errors: string[] = [];
    setPartErrorMsg(""); // Clear previous errors

    // Collect errors based on the current state
    if (part.PartNumber.trim() === "") {
      errors.push("Part number");
    }
    if (part.Length <= 0) {
      errors.push("Length");
    }
    if (part.Quantity <= 0) {
      errors.push("Quantity");
    }
    if (part.MaterialCode.trim() === "") {
      errors.push("Material Code");
    }

    // If there are errors, join them into a single string and set the error message
    if (errors.length !== 0) {
      setPartErrorMsg("Required Part Fields are Invalid: " + errors.join(", "));
      return false;
    }

    // If no errors, return true
    return true;
  }

  function addToParts(e: any): void {
    e.preventDefault();
    if (checkPartValidity()) {
      setParts((prev) => [...prev, part]);
    }
  }

  function updateMaterial(key: keyof Material, value: string | number): void {
    setMaterial((prev) => ({
      ...prev,
      [key]: !value && value !== 0 ? materialInitialState[key] : value,
    }));
  }

  function checkMaterialValidity(): boolean {
    const errors: string[] = [];
    // Check if MaterialCode is a non-empty string
    if (material.MaterialCode.trim() === "") {
      errors.push("Material Code");
    }

    // Check if Length is greater than 0
    if (material.Length <= 0) {
      errors.push("Length");
    }

    // Check if Quantity is greater than 0
    if (material.Quantity <= 0) {
      errors.push("Quantity");
    }

    if (errors.length !== 0) {
      setMaterialErrorMsg(
        "Required Material Fields are Invalid: " + errors.join(", "),
      );
      return false;
    }

    // All checks passed, return true
    return true;
  }

  function addToMaterials(e: any): void {
    e.preventDefault();
    if (checkMaterialValidity()) {
      // Update state to add the new material
      setMaterials((prev) => [...prev, material]);
      // Clear any previous error message
      setMaterialErrorMsg("");
    } else {
      // Set an error message to indicate invalid materials
      setMaterialErrorMsg("Materials must be valid");
    }
  }

  function updateMaterialQty(): void {
    if (material.Quantity === 9999) {
      console.log("Q in If", material.Quantity);
      updateMaterial("Quantity", 0);
    }
    if (material.Quantity === 0) {
      console.log("Q in If 2", material.Quantity);
      updateMaterial("Quantity", 9999);
    }

    setMaterialQtyDisabled(!materialQtyDisabled);
  }

  return (
    <div>
      <h1>Create New Job</h1>
      <div className={styles.heading}>
        <div className={styles.heading__left}>
          <StyledInput
            type={"text"}
            placeholder={"Enter Project #"}
            value={jobInfo?.Job}
            onChange={(e) =>
              setJobInfo((prev) => ({
                ...prev,
                Job: e.target.value.toUpperCase(),
              }))
            }
          />

          <StyledInput
            type={"text"}
            placeholder={"Enter Customer"}
            value={jobInfo?.Customer}
            onChange={(e) =>
              setJobInfo((prev) => ({
                ...prev,
                Customer: e.target.value.toUpperCase(),
              }))
            }
          />
          <Button onClick={runProject}>Run Project</Button>
        </div>
        <div className={styles.heading__right}>
          <Badge badgeContent={parts.length} color="primary">
            <Button>
              <p>Parts</p>
              <Segment />
            </Button>
          </Badge>
          <Badge badgeContent={materials.length} color="primary">
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
            <form onSubmit={addToParts}>
              <div className={styles.part__addition__inputs__form}>
                <StyledInput
                  type={"text"}
                  placeholder="Part Number"
                  value={part.PartNumber}
                  onChange={(e) => updatePart("PartNumber", e.target.value)}
                  // onFocus={handleFocus}
                  // onBlur={handleBlur}
                />

                <StyledInput
                  type={"text"}
                  placeholder="Material Code"
                  value={part.MaterialCode}
                  onChange={(e) => updatePart("MaterialCode", e.target.value)}
                />

                <StyledInput
                  type={"number"}
                  placeholder="Part Length"
                  value={part.Length !== 0 ? part.Length : ""}
                  onChange={(e) => updatePart("Length", e.target.value)}
                />

                <StyledInput
                  type={"text"}
                  placeholder={"Cutting Operation"}
                  value={part.CuttingOperation}
                  onChange={(e) =>
                    updatePart("CuttingOperation", e.target.value)
                  }
                />
                <StyledInput
                  type={"number"}
                  placeholder="Part Quantity"
                  value={part.Quantity !== 0 ? part.Quantity : ""}
                  onChange={(e) => updatePart("Quantity", e.target.value)}
                />
              </div>
              <Button type="submit">
                Add Part <Add />{" "}
              </Button>
            </form>
          </div>
          {partErrorMsg && <div className="error">{partErrorMsg}</div>}
        </div>

        <div className={styles.material__addition}>
          <h4>Add Material</h4>
          <div className={styles.material__addition__inputs}>
            <form onSubmit={addToMaterials}>
              <div className={styles.material__addition__inputs__form}>
                <StyledInput
                  type={"text"}
                  placeholder={"Material Code"}
                  value={material.MaterialCode}
                  onChange={(e) =>
                    updateMaterial(
                      "MaterialCode",
                      e.target.value.toUpperCase().trim(),
                    )
                  }
                />
                <StyledInput
                  type={"number"}
                  placeholder={"Material Length - inch"}
                  value={material.Length !== 0 ? material.Length : ""}
                  onChange={(e) =>
                    updateMaterial("Length", parseFloat(e.target.value))
                  }
                />
                <StyledInput
                  type={"text"}
                  placeholder={"Material Quantity"}
                  // value={materialQtyDisabled}
                  value={
                    material.Quantity === 9999
                      ? "Unlimited"
                      : material.Quantity !== 0
                        ? material.Quantity
                        : ""
                  }
                  onChange={(e) =>
                    updateMaterial("Quantity", parseInt(e.target.value))
                  }
                  disabled={materialQtyDisabled}
                />
                <div className={styles.checkbox}>
                  <input
                    type={"checkbox"}
                    onClick={() => updateMaterialQty()}
                  />
                  <p>Check for unlimited material</p>
                </div>
              </div>
              <Button type="submit">
                Add Material
                <Add />
              </Button>
            </form>
          </div>
          {materialErrorMsg && <div className="error">{materialErrorMsg}</div>}
        </div>
      </div>
    </div>
  );
};
