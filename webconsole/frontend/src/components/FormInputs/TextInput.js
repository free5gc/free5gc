import React from 'react';
import cx from 'classnames';

const TextInput = ({
  input,
  label,
  type,
  meta: { touched, error, warning },
  inputClassName,
  placeholder,
  helpText,
  disabled
}) => (
  <div>
    <input
      {...input}
      type={type}
      className={cx(inputClassName, 'form-control', {
        error: !!error
      })}
      placeholder={placeholder}
      disabled={disabled} />

    { touched && error &&
      <label className="error" htmlFor={input.name}>{error}</label>
    }

    { helpText &&
      <span className="help-block">{helpText}</span>
    }
  </div>
);

export default TextInput;