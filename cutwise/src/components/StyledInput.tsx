import React, { FunctionComponent } from "react";
import styles from "../styles/StyledInput.module.css";

interface Props {
  type: "text" | "number"; // Specifies the type of the input
  placeholder?: string; // Optional placeholder text
  value?: string | number | boolean; // Optional value (should be a string)
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; // Optional change handler
  onFocus?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional focus handler
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional blur handler
  disabled?: boolean;
}

export const StyledInput: FunctionComponent<Props> = ({
  type,
  placeholder,
  value = "",
  onChange,
  onFocus,
  onBlur,
  disabled,
}) => {
  return (
    <input
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      onFocus={onFocus}
      onBlur={onBlur}
      disabled={disabled}
      className={styles.styled__input} // Example className, apply your styles here
    />
  );
};
