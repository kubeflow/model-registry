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
  descriptionTextNode?: React.ReactNode;    // Always-visible help text
  helperTextNode?: React.ReactNode; // Error-only helper text
  className?: string; // Optional className for the outer FormGroup
  labelHelp?: React.ReactElement;
  'data-testid'?: string;
};

const ThemeAwareFormGroupWrapper: React.FC<ThemeAwareFormGroupWrapperProps> = ({
  children,
  label,
  fieldId,
  isRequired,
  descriptionTextNode,
  helperTextNode,
  className,
  labelHelp,
  'data-testid': dataTestId,
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
          labelHelp={labelHelp}
          data-testid={dataTestId}
        >
          <FormFieldset component={children} field={label} />
        </FormGroup>
        {descriptionTextNode}
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
        labelHelp={labelHelp}
        data-testid={dataTestId}
      >
        {children}
        {descriptionTextNode}
        {helperTextNode}
      </FormGroup>
    </>
  );
};

export default ThemeAwareFormGroupWrapper;
