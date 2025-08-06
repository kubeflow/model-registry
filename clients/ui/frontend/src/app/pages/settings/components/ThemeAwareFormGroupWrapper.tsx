import * as React from 'react';
import { FormGroup } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

// Props required by this wrapper component
type ThemeAwareFormGroupWrapperProps = {
  children: React.ReactNode; // The input component
  label: string;
  fieldId: string;
  isRequired?: boolean;
  helperTextNode?: React.ReactNode; // The pre-rendered HelperText component or null
  className?: string; // Optional className for the outer FormGroup
};

const ThemeAwareFormGroupWrapper: React.FC<ThemeAwareFormGroupWrapperProps> = ({
  children,
  label,
  fieldId,
  isRequired,
  helperTextNode,
  className,
}) => {
  const { isMUITheme } = useThemeContext();
  const hasError = !!helperTextNode; // Determine error state based on helper text presence

  if (isMUITheme) {
    // For MUI theme, render FormGroup -> FormFieldset -> Input
    // Helper text is rendered *after* the FormGroup wrapper
    return (
      <>
        <FormGroup
          className={`${className || ''} ${hasError ? 'pf-m-error' : ''}`.trim()} // Apply className and error state class
          label={label}
          isRequired={isRequired}
          fieldId={fieldId}
        >
          <FormFieldset component={children} field={label} />
        </FormGroup>
        {helperTextNode}
      </>
    );
  }

  // For PF theme, render standard FormGroup
  return (
    <>
      <FormGroup
        className={`${className || ''} ${hasError ? 'pf-m-error' : ''}`.trim()} // Apply className and error state class
        label={label}
        isRequired={isRequired}
        fieldId={fieldId}
      >
        {children}
        {helperTextNode}
      </FormGroup>
    </>
  );
};

export default ThemeAwareFormGroupWrapper;
