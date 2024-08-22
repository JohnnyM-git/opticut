import React, { FunctionComponent, useEffect } from "react";
import styles from "../styles/StyledInput.module.css";

interface Props {
  type: "text" | "number"; // Specifies the type of the input
  placeholder?: string; // Optional placeholder text
  value?: string | number; // Optional value (should be a string)
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; // Optional change handler
  onFocus?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional focus handler
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional blur handler
}

export const StyledInput: FunctionComponent<Props> = ({
  type,
  placeholder,
  value = "",
  onChange,
  onFocus,
  onBlur,
}) => {
  return (
    <input
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      onFocus={onFocus}
      onBlur={onBlur}
      className={styles.styled__input} // Example className, apply your styles here
    />
  );
};
