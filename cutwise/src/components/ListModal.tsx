import styles from "./styles/ListModal.module.css";
import { StyledInput } from "./StyledInput.tsx";
import { useEffect } from "react";

interface props {
  listName: string;
  list: any;
}

export const ListModal = (props: any) => {
  useEffect(() => {
    console.log(props.list);
  }, []);

  return (
    <div className={styles.list__modal__backdrop}>
      <div className={styles.list__modal}>
        <h1>{props.listName === "parts" ? "Parts" : "Materials"}</h1>
      </div>
      <div className={styles.list__modal__parts}>
        {props?.list === "parts"
          ? props?.list.map((item: any) => (
              <div key={item.PartNumber}>
                <StyledInput type={"text"} value={item.PartNumber} />
              </div>
            ))
          : ""}
      </div>
    </div>
  );
};
