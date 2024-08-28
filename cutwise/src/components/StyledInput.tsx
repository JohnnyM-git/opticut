import React, { FunctionComponent } from "react";
import styles from "../styles/StyledInput.module.css";

interface Props {
  type: "text" | "number" | "radio"; // Specifies the type of the input
  placeholder?: string; // Optional placeholder text
  value?: string | number; // Optional value (should be a string)
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void; // Optional change handler
  onFocus?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional focus handler
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void; // Optional blur handler
  disabled?: boolean;
  step?: string;
  name?: string;
  checked?: boolean;
}

export const StyledInput: FunctionComponent<Props> = ({
  type,
  placeholder,
  value = "",
  onChange,
  onFocus,
  onBlur,
  disabled,
  step,
  name,
  checked,
}) => {
  return (
    <div className={styles.inputWrapper}>
      <input
        type={type}
        // placeholder={placeholder}
        value={value}
        onChange={onChange}
        onFocus={onFocus}
        onBlur={onBlur}
        disabled={disabled}
        name={name}
        checked={checked}
        step={type === "number" ? step : undefined}
        className={`${styles.styledInput} ${value || placeholder ? styles.hasValue : ""}`}
      />
      {placeholder && (
        <label className={styles.floatingLabel}>{placeholder}</label>
      )}
    </div>
  );
};
