import * as React from 'react';
import {
  Alert,
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
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import { validateYamlContent } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  DESCRIPTION_TEXT,
  ERROR_MESSAGES,
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
  const [isYamlTouched, setIsYamlTouched] = React.useState(false);
  const [filename, setFilename] = React.useState('');
  const [fileUploadError, setFileUploadError] = React.useState<string | undefined>(undefined);
  const isYamlContentValid = validateYamlContent(formData.yamlContent);

  const handleFileChange = (
    _event: React.DragEvent<HTMLElement> | React.ChangeEvent<HTMLInputElement> | Event,
    file: File,
  ) => {
    setFilename(file.name);
    setFileUploadError(undefined);
    const reader = new FileReader();
    reader.onload = () => {
      const text = typeof reader.result === 'string' ? reader.result : '';
      setData('yamlContent', text);
      setIsYamlTouched(true);
    };
    reader.onerror = () => {
      setFileUploadError(ERROR_MESSAGES.FILE_UPLOAD_FAILED_BODY);
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
    setFileUploadError(undefined);
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

  const yamlDescriptionTxtNode = (
    <FormHelperText>
      <HelperText>
        <HelperTextItem>{DESCRIPTION_TEXT.YAML}</HelperTextItem>
      </HelperText>
    </FormHelperText>
  );

  const yamlHelperTxtNode =
    isYamlTouched && !isYamlContentValid ? (
      <FormHelperText>
        <HelperText>
          <HelperTextItem variant="error" data-testid="yaml-content-error">
            {VALIDATION_MESSAGES.YAML_CONTENT_REQUIRED}
          </HelperTextItem>
        </HelperText>
      </FormHelperText>
    ) : undefined;

  return (
    <FormSection data-testid="yaml-section">
      {fileUploadError && (
        <Alert
          variant="danger"
          isInline
          title={ERROR_MESSAGES.FILE_UPLOAD_FAILED}
          className="pf-v6-u-mb-md"
          data-testid="yaml-file-upload-error"
        >
          {fileUploadError}
        </Alert>
      )}
      <FormGroup
        label={
          <Flex
            justifyContent={{ default: 'justifyContentSpaceBetween' }}
            alignItems={{ default: 'alignItemsCenter' }}
          >
            <FlexItem>{FORM_LABELS.YAML_CONTENT}</FlexItem>
            {onToggleExpectedFormatDrawer && (
              <FlexItem>
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
              </FlexItem>
            )}
          </Flex>
        }
        isRequired
        fieldId="yaml-content"
      >
        {yamlDescriptionTxtNode}
        <FormFieldset component={yamlInput} field="YAML" />
        {yamlHelperTxtNode}
      </FormGroup>
    </FormSection>
  );
};

export default YamlSection;
