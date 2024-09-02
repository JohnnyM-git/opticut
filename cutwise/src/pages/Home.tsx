import { FunctionComponent, useState } from "react";
import styles from "../styles/Home.module.css";
import { Badge, Button } from "@mui/material";
import { Add, Segment } from "@mui/icons-material";
import { StyledInput } from "../components/StyledInput.tsx";
import { apiUrl } from "../globals.ts";
import { useNavigate } from "react-router-dom";
import { useSettings } from "../SettingsContext.tsx";
import { mmToIn } from "../functions/unitConverter.ts";
import { ListModal } from "../components/ListModal.tsx";
import { dialog } from "@tauri-apps/api";
// import { dialog } from "@tauri-apps/api";
// import { invoke } from "@tauri-apps/api/tauri";
// import axios from "axios";
// import { readBinaryFile } from "tauri/api/fs";
// import { appWindow } from '@tauri-apps/api/window';

interface Props {
  toggleModal: () => void;
}

export const Home: FunctionComponent<Props> = ({ toggleModal }) => {
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

  interface ModalState {
    Open: boolean;
    list: string;
  }

  const initialModalState = {
    Open: false,
    list: "parts",
  };

  const [jobInfo, setJobInfo] = useState<JobInfo>({ Job: "", Customer: "" });
  const [parts, setParts] = useState<Part[]>([]);
  const [materials, setMaterials] = useState<Material[]>([]);
  const [part, setPart] = useState<Part>(partInitialState);
  const [material, setMaterial] = useState<Material>(materialInitialState);
  const [partErrorMsg, setPartErrorMsg] = useState<string | null>("");
  const [materialErrorMsg, setMaterialErrorMsg] = useState<string>("");
  const [materialQtyDisabled, setMaterialQtyDisabled] = useState(false);
  const navigate = useNavigate();
  const { settings } = useSettings();
  const [modal, setModal] = useState<ModalState>(initialModalState);
  // const [file, setFile] = useState()

  async function runProject(): Promise<void> {
    if (checkDataValidity()) {
      if (settings.units === "metric") {
        await updateUnits();
      }
      console.log("PARTS:", parts);

      console.log("Materials:", materials);
      const res = await fetch(`${apiUrl}run-project`, {
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

    if (
      parts.some((existingPart) => existingPart.PartNumber === part.PartNumber)
    ) {
      setPartErrorMsg("Part already exists and cannot be added again");
      return false;
    }

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

  async function updateUnits() {
    parts.forEach((part) => {
      part.Length = mmToIn(part.Length);
    });

    materials.forEach((material) => {
      material.Length = mmToIn(material.Length);
    });
  }

  function checkMaterialValidity(): boolean {
    const errors: string[] = [];

    if (
      materials.some(
        (existingMaterial) =>
          existingMaterial.MaterialCode === material.MaterialCode &&
          existingMaterial.Length === material.Length,
      )
    ) {
      setMaterialErrorMsg(
        "Current material at specified length already exists. If you need to add more please" +
          " update the quantity",
      );
      return false;
    }

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
    }
  }

  function updateMaterialQty(): void {
    updateMaterial("Quantity", 1);
    if (!materialQtyDisabled) {
      updateMaterial("Quantity", 9999);
    }
    setMaterialQtyDisabled(!materialQtyDisabled);
    console.log(material.Quantity);
  }

  function openModal(list: string) {
    console.log(list);
    toggleModal();
    setModal((prevModal) => ({
      ...prevModal,
      Open: !prevModal.Open,
      list: list,
    }));
  }
  function closeModal() {
    setModal((prevModal) => ({
      ...prevModal,
      Open: false,
    }));
    toggleModal();
  }

  function deletePart(partIdentifier: string) {
    setParts((prevState) => {
      // Determine if partIdentifier is for Part or Material
      return prevState.filter(
        (part) =>
          part.PartNumber !== partIdentifier &&
          part.MaterialCode !== partIdentifier,
      );
    });
  }

  function updatePartInArray(partIdentifier: string, newState: Partial<Part>) {
    setParts((prevState) =>
      prevState.map((part) =>
        part.PartNumber === partIdentifier ||
        part.MaterialCode === partIdentifier
          ? { ...part, ...newState }
          : part,
      ),
    );
  }

  function deleteMaterial(materialIndex: number) {
    setMaterials((prevMaterials) =>
      prevMaterials.filter((_, index) => index !== materialIndex),
    );
  }

  function updateMaterialInArray(
    materialIndex: number,
    newState: Partial<Material>,
  ) {
    setMaterials((prevMaterials) => {
      // Create a new array with the updated material
      return prevMaterials.map((material, index) =>
        index === materialIndex ? { ...material, ...newState } : material,
      );
    });
  }

  // async function selectAndUploadFile() {
  //     // console.log('App window label:', window.__TAURI__.appWindow.label);
  //     console.log(file)
  //
  //     const formData = new FormData();
  //
  //        formData.append(
  //            "file",
  //            file.name
  //           )
  //
  //    console.log(formData)
  //
  //
  //     // Open file dialog to select a file
  //     // const selected = await dialog.open({
  //     //     filters: [{
  //     //         name: 'Excel Files',
  //     //         extensions: ['xlsx', 'xls'],
  //     //     }]
  //     // });
  //     //
  //     // if (typeof selected === 'string') {
  //     //     try {
  //     //         console.log('Selected file:', selected);
  //     //
  //     //         // Read the file as binary data
  //     //         const fileData = await readBinaryFile(selected);
  //     //         console.log(fileData)
  //     //
  //     //         // Prepare the file for upload
  //     //         const formData = new FormData();
  //     //         // @ts-ignore
  //     //         const blob = new Blob([fileData], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  //     //         formData.append('file', blob, 'file.xlsx');
  //     //
  //     //         // Send the file to the Go backend
  //     //         const response = await axios.post('http://localhost:2828/file-upload', formData, {
  //     //             headers: {
  //     //                 'Content-Type': 'multipart/form-data',
  //     //             },
  //     //         });
  //     //
  //     //         console.log(response.data); // Process and display the response data
  //     //     } catch (e) {
  //     //         console.error('Error uploading file:', e);
  //     //     }
  //     // }
  // }
  //
  // // Helper function to read file as binary
  // async function readBinaryFile(filePath: string) {
  //     const response = await fetch(filePath);
  //     return response.arrayBuffer();
  // }

  // async function testFileSelection() {
  //     const selected = await dialog.open({
  //         filters: [{
  //             name: 'Excel Files',
  //             extensions: ['xlsx', 'xls'],
  //         }]
  //     });
  //     if (typeof selected === 'string') {
  //         console.log('Selected file:', selected);
  //         const file = await readBinaryFile(selected);
  //         console.log('File content:', file);
  //     }
  // }

  async function runUploadedFiles() {
    try {
      setParts([]);
      setMaterials([]);
      const filePath = await dialog.open({
        filters: [
          {
            name: "Excel Files",
            extensions: ["xlsx", "xls"],
          },
        ],
      });

      console.log(filePath);

      const res = await fetch(`${apiUrl}file-upload?filePath=${filePath}`);

      // const res = await fetch(`${apiUrl}fileupload`, {});
      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      // const data: JobsResponse = await res.json();
      const data = await res.json();
      console.log(data);
      data.parts.forEach((part: Part) => {
        setParts((prevState) => [...prevState, part]);
      });

      data.materials.forEach((material: Material) => {
        setMaterials((prevState) => [...prevState, material]);
      });

      setJobInfo((prevState) => ({
        ...prevState,
        Job: data.Job.Job,
      }));

      setJobInfo((prevState) => ({
        ...prevState,
        Customer: data.Job.Customer,
      }));

      // setJobs(data.JobsList); // Update state with JobsList
    } catch (error) {
      console.error("Uploading File", error);
    }
  }

  async function batchProcess() {
    const startTime = performance.now(); // Start timer

    const res = await fetch(`${apiUrl}batch-run-files`);
    const data = await res.json();

    const endTime = performance.now(); // End timer

    console.log(data);

    // Calculate the time taken
    const timeTaken = endTime - startTime;
    console.log(
      `Batch processing completed in ${timeTaken.toFixed(2)} milliseconds`,
    );

    // if (data.errors.length > 0) {
    //   console.log(data.errors);
    //   return;
    // }
    // navigate(`/localjobs`);
    console.log(data);
  }

  return (
    <>
      {modal.Open && (
        <div className={styles.modal}>
          <ListModal
            listName={modal.list}
            list={modal.list === "parts" ? parts : materials}
            deletePart={deletePart}
            deleteMaterial={deleteMaterial}
            updatePart={updatePartInArray}
            updateMaterialInArray={updateMaterialInArray}
            job={jobInfo?.Job}
            closeModal={closeModal}
          />
        </div>
      )}

      <div className={`${modal.Open ? styles.modal__open : styles.home}`}>
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
            {/*  <input type={"file"} onChange={(e) => setFile(e.target.files[0])}/>*/}
            <Button onClick={runUploadedFiles}>Upload File</Button>
            <Button onClick={batchProcess}>Batch Process</Button>
          </div>
          <div className={styles.heading__right}>
            <Badge badgeContent={parts.length} color="primary">
              <Button onClick={() => openModal("parts")}>
                <p>Parts</p>
                <Segment />
              </Button>
            </Badge>
            <Badge badgeContent={materials.length} color="primary">
              <Button onClick={() => openModal("materials")}>
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
              <form className={styles.form} onSubmit={addToParts}>
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
                    placeholder={
                      settings.units === "imperial"
                        ? "Part Length - Inch"
                        : "Part Length - mm"
                    }
                    // value={part.Length}
                    value={part.Length !== 0 ? part.Length : ""}
                    onChange={(e) => updatePart("Length", e.target.value)}
                    step={"0.01"}
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
              <form className={styles.form} onSubmit={addToMaterials}>
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
                    placeholder={
                      settings.units === "imperial"
                        ? "Material Length - Inch"
                        : "Material Length - mm"
                    }
                    value={material.Length !== 0 ? material.Length : ""}
                    // value={material.Length}
                    onChange={(e) =>
                      updateMaterial("Length", parseFloat(e.target.value))
                    }
                    step={"0.0001"}
                  />
                  <StyledInput
                    type={"text"}
                    placeholder={
                      materialQtyDisabled
                        ? "Material Qty Unlimited"
                        : "Material Quantity"
                    }
                    // value={material.Quantity}
                    value={
                      material.Quantity !== 0 && !materialQtyDisabled
                        ? material.Quantity
                        : ""
                    }
                    onChange={(e) =>
                      updateMaterial("Quantity", parseInt(e.target.value))
                    }
                    disabled={materialQtyDisabled}
                  />
                  <div
                    className={styles.checkbox}
                    onClick={() => updateMaterialQty()}
                  >
                    <input
                      type={"checkbox"}
                      checked={materialQtyDisabled}
                      onChange={() => updateMaterialQty()}
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
            {materialErrorMsg && (
              <div className="error">{materialErrorMsg}</div>
            )}
          </div>
        </div>
      </div>
    </>
  );
};
