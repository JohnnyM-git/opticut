import styles from "./styles/ListModal.module.css";
import { StyledInput } from "./StyledInput";
import { useEffect } from "react";
import {Close, Delete, Save} from "@mui/icons-material";

interface Part {
    PartNumber: string;
    MaterialCode: string;
    Length: number;
    Quantity: number;
    CuttingOperation: string;
}

interface Material {
    MaterialCode: string;
    Length: number;
    Quantity: number;
}

interface Props {
    listName: string;
    list: Part[] | Material[];
    deletePart: (partNumber: string) => void;
    deleteMaterial: (materialIndex: number) => void;
    updatePart: (partNumber: string, newState: Partial<Part> | Partial<Material>) => void;
    updateMaterialInArray: (materialIndex: number, newState: Partial<Material>) => void;
    job: string;
    closeModal: () => void;
}

export const ListModal = (props: Props) => {
    useEffect(() => {
        console.log(props.list);
    }, [props.list]);

    const handlePartInputChange = (e: React.ChangeEvent<HTMLInputElement>, partNumber: string, field: string) => {
        const value = e.target.value;
        props.updatePart(partNumber, { [field]: value });
    };

    const handleMaterialInputChange = (e: React.ChangeEvent<HTMLInputElement>, index: number, field: string) => {
        const value = e.target.value;
        props.updateMaterialInArray(index, { [field]: value });
    };

    return (

        <div className={styles.list__modal}>
            <div className={styles.modal__heading}>

                <h1>{props.job} - {props.listName?.toUpperCase()}</h1>
                <button onClick={() => props.closeModal()}>
                    <Close/>
                </button>
            </div>
            {/*<h1>{props.listName === "parts" ? "Parts" : "Materials"}</h1>*/}

            <div className={styles.list__modal__parts}>
                {props.listName === "parts"
                    ? (props.list as Part[]).map((item: Part) => (
                        <div key={item.PartNumber} className={styles.list__modal__parts__part}>
                            <StyledInput
                                type="text"
                                value={item.PartNumber}
                                onChange={(e) => handlePartInputChange(e, item.PartNumber, 'PartNumber')}
                            />
                            <StyledInput
                                type="text"
                                value={item.MaterialCode}
                                onChange={(e) => handlePartInputChange(e, item.PartNumber, 'MaterialCode')}
                            />
                            <StyledInput
                                type="number"
                                value={item.Length}
                                onChange={(e) => handlePartInputChange(e, item.PartNumber, 'Length')}
                            />
                            <StyledInput
                                type="text"
                                value={item.CuttingOperation}
                                onChange={(e) => handlePartInputChange(e, item.PartNumber, 'CuttingOperation')}
                            />
                            <StyledInput
                                type="number"
                                value={item.Quantity}
                                onChange={(e) => handlePartInputChange(e, item.PartNumber, 'Quantity')}
                            />
                            <button onClick={() => props.deletePart(item.PartNumber)}>
                                <Delete/>
                            </button>
                            <button onClick={() => props.updatePart(item.PartNumber, item)}>
                                <Save/>
                            </button>
                        </div>
                    ))
                    : (props.list as Material[]).map((item: Material, index: number) => (
                        <div key={index} className={styles.list__modal__parts__part}>
                            <StyledInput
                                type="text"
                                value={item.MaterialCode}
                                onChange={(e) => handleMaterialInputChange(e, index, 'MaterialCode')}
                            />
                            <StyledInput
                                type="number"
                                value={item.Length}
                                onChange={(e) => handleMaterialInputChange(e, index, 'Length')}
                            />
                            <StyledInput
                                type="number"
                                value={item.Quantity}
                                onChange={(e) => handleMaterialInputChange(e, index, 'Quantity')}
                            />
                            <button onClick={() => props.deleteMaterial(index)}>
                                <Delete/>
                            </button>
                            <button onClick={() => props.updateMaterialInArray(index, item)}>
                                <Save/>
                            </button>
                        </div>
                    ))
                }
            </div>
        </div>

    );
};
