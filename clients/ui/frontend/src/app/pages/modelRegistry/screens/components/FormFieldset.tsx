import { useThemeContext } from 'mod-arch-kubeflow';
import React, { ReactNode } from 'react';

interface FormFieldsetProps {
  component: ReactNode;
  field?: string;
  className?: string;
}

const FormFieldset: React.FC<FormFieldsetProps> = ({ component, field, className }) => {
  const { isMUITheme } = useThemeContext();

  if (!isMUITheme) {
    return <>{component}</>;
  }

  return (
    <div className={className ?? 'form-fieldset-wrapper'}>
      {component}
      <fieldset aria-hidden="true" className="form-fieldset">
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
