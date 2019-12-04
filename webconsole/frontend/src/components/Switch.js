import React from 'react';
import cx from 'classnames';

const SwitchControl = ({
  value,
  onChange,
  onText,
  offText
}) => (
  <div className="switch has-switch">
    <div className={cx("switch-animate", {
      'switch-on': value,
      'switch-off': !value
    })}
      onClick={() => onChange(!value)}>
      <span className="switch-left">{onText}</span>
      <label>&nbsp;</label>
      <span className="switch-right">{offText}</span>
    </div>
  </div>
);

SwitchControl.defaultProps = {
  onText: 'ON',
  offText: 'OFF',
  onChange: () => {}
};

export default SwitchControl;