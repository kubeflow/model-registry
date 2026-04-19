import * as React from 'react';
import {
  Button,
  Flex,
  FlexItem,
  FormGroup,
  FileUpload,
  FormHelperText,
  HelperText,
  HelperTextItem,
} from '@patternfly/react-core';
import { OpenDrawerRightIcon } from '@patternfly/react-icons';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { validateYamlContent } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  HELP_TEXT,
  EXPECTED_YAML_FORMAT_LABEL,
} from '~/app/pages/modelCatalogSettings/constants';

type YamlSectionProps = {
  formData: ManageSourceFormData;
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
  onToggleExpectedFormatDrawer?: () => void;
};

const YamlSection: React.FC<YamlSectionProps> = ({
  formData,
  setData,
  onToggleExpectedFormatDrawer,
}) => {
  const { isMUITheme } = useThemeContext();
  const [isYamlTouched, setIsYamlTouched] = React.useState(false);
  const [filename, setFilename] = React.useState('');
  const isYamlContentValid = validateYamlContent(formData.yamlContent);

  const handleFileChange = (
    _event: React.DragEvent<HTMLElement> | React.ChangeEvent<HTMLInputElement> | Event,
    file: File,
  ) => {
    setFilename(file.name);
    const reader = new FileReader();
    reader.onload = () => {
      const text = typeof reader.result === 'string' ? reader.result : '';
      setData('yamlContent', text);
      setIsYamlTouched(true);
    };
    reader.readAsText(file);
  };

  const handleTextChange = (_event: React.ChangeEvent<HTMLTextAreaElement>, value: string) => {
    setData('yamlContent', value);
  };

  const handleClear = () => {
    setFilename('');
    setData('yamlContent', '');
    setIsYamlTouched(true);
  };

  const yamlInput = (
    <div data-testid="yaml-content-input">
      <FileUpload
        id="yaml-content"
        type="text"
        value={formData.yamlContent}
        filename={filename}
        filenamePlaceholder="Drag and drop a YAML file or upload one"
        onFileInputChange={handleFileChange}
        onTextChange={handleTextChange}
        onClearClick={handleClear}
        onBlur={() => setIsYamlTouched(true)}
        validated={isYamlTouched && !isYamlContentValid ? 'error' : 'default'}
        browseButtonText="Upload"
        allowEditingUploadedText
        dropzoneProps={{
          accept: { 'text/yaml': ['.yaml', '.yml'] },
        }}
      />
    </div>
  );

  const expectedFormatButton = onToggleExpectedFormatDrawer ? (
    <Button
      variant="link"
      isInline
      onClick={onToggleExpectedFormatDrawer}
      data-testid="view-expected-yaml-format-link"
      icon={<OpenDrawerRightIcon />}
      iconPosition="end"
    >
      {EXPECTED_YAML_FORMAT_LABEL}
    </Button>
  ) : null;

  const descriptionTextNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>{HELP_TEXT.YAML}</HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  const hasError = isYamlTouched && !isYamlContentValid;
  const helperTextNode = hasError ? (
    <FormHelperText>
      <HelperText>
        <HelperTextItem variant="error" data-testid="yaml-content-error">
          {VALIDATION_MESSAGES.YAML_CONTENT_REQUIRED}
        </HelperTextItem>
      </HelperText>
    </FormHelperText>
  ) : undefined;

  if (isMUITheme) {
    return (
      <FormSection data-testid="yaml-section">
        {expectedFormatButton && (
          <Flex justifyContent={{ default: 'justifyContentFlexEnd' }}>
            <FlexItem>{expectedFormatButton}</FlexItem>
          </Flex>
        )}
        {descriptionTextNode}
        <FormGroup
          className={hasError ? 'pf-m-error' : undefined}
          label={FORM_LABELS.YAML_CONTENT}
          isRequired
          fieldId="yaml-content"
        >
          <FormFieldset component={yamlInput} field={FORM_LABELS.YAML_CONTENT} />
        </FormGroup>
        {helperTextNode}
      </FormSection>
    );
  }

  return (
    <FormSection data-testid="yaml-section">
      <FormGroup
        label={
          <Flex
            justifyContent={{ default: 'justifyContentSpaceBetween' }}
            alignItems={{ default: 'alignItemsCenter' }}
          >
            <FlexItem>{FORM_LABELS.YAML_CONTENT}</FlexItem>
            {expectedFormatButton && <FlexItem>{expectedFormatButton}</FlexItem>}
          </Flex>
        }
        isRequired
        fieldId="yaml-content"
      >
        <FormFieldset component={yamlInput} field="YAML" />
        {descriptionTextNode}
        {helperTextNode}
      </FormGroup>
    </FormSection>
  );
};

export default YamlSection;
