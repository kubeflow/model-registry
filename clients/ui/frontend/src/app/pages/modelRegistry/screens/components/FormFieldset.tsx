import { useThemeContext } from 'mod-arch-kubeflow';
import React, { CSSProperties, ReactNode } from 'react';
import './FormFieldset.scss';

interface FormFieldsetProps {
  component: ReactNode;
  field?: string;
  className?: string;
  fieldsetStyle?: CSSProperties;
}

const FormFieldset: React.FC<FormFieldsetProps> = ({
  component,
  field,
  className,
  fieldsetStyle,
}) => {
  const { isMUITheme } = useThemeContext();

  if (!isMUITheme) {
    return <>{component}</>;
  }

  return (
    <div className={className ?? 'form-fieldset-wrapper'}>
      {component}
      <fieldset aria-hidden="true" className="form-fieldset" style={fieldsetStyle}>
        {field && (
          <legend className="form-fieldset-legend">
            <span>{field}</span>
          </legend>
        )}
      </fieldset>
    </div>
  );
};

export default FormFieldset;
