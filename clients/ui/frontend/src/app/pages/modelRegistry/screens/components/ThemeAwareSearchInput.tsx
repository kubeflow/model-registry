import * as React from 'react';
import { SearchInput, SearchInputProps, TextInput } from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

type ThemeAwareSearchInputProps = Omit<SearchInputProps, 'onChange' | 'onClear'> & {
  onChange: (value: string) => void; // Simplified onChange signature
  onClear?: () => void; // Simplified optional onClear signature
  fieldLabel?: string; // Additional prop for MUI FormFieldset label
  'data-testid'?: string;
};

const ThemeAwareSearchInput: React.FC<ThemeAwareSearchInputProps> = ({
  value,
  onChange,
  onClear,
  fieldLabel,
  placeholder,
  isDisabled,
  className,
  style,
  'aria-label': ariaLabel = 'Search',
  'data-testid': dataTestId,
  ...rest
}) => {
  const { isMUITheme } = useThemeContext();

  if (isMUITheme) {
    // Render MUI version using TextInput + FormFieldset
    return (
      <FormFieldset
        className={className}
        field={fieldLabel}
        component={
          <TextInput
            value={value}
            type="text"
            onChange={(_event, newValue) => onChange(newValue)} // Adapt signature
            isDisabled={isDisabled}
            aria-label={ariaLabel}
            data-testid={dataTestId}
            style={style}
          />
        }
      />
    );
  }

  // Render PF version using SearchInput
  return (
    <SearchInput
      {...rest} // Pass all other applicable SearchInputProps
      className={className}
      style={style}
      placeholder={placeholder}
      value={value}
      isDisabled={isDisabled}
      aria-label={ariaLabel}
      data-testid={dataTestId}
      onChange={(_event, newValue) => onChange(newValue)} // Adapt signature
      onClear={(event) => {
        event.stopPropagation();
        onChange('');
        onClear?.(); // Adapt signature
      }}
    />
  );
};

export default ThemeAwareSearchInput;
