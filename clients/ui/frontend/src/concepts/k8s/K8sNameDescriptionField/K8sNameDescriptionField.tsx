import * as React from 'react';
import {
  Button,
  FormGroup,
  HelperText,
  HelperTextItem,
  TextArea,
  TextInput,
} from '@patternfly/react-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import ResourceNameDefinitionTooltip from '~/concepts/k8s/ResourceNameDefinitionTootip';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import ResourceNameField from './ResourceNameField';

// TODO: replace with the actual call once we have the endpoint

/** Companion data hook */
// export const useK8sNameDescriptionFieldData = (
//   configuration: UseK8sNameDescriptionDataConfiguration = {},
// ): UseK8sNameDescriptionFieldData => {
//   const [data, setData] = React.useState<K8sNameDescriptionFieldData>(() =>
//     setupDefaults(configuration),
//   );

//   const onDataChange = React.useCallback<K8sNameDescriptionFieldUpdateFunction>((key, value) => {
//     setData((currentData) => handleUpdateLogic(currentData)(key, value));
//   }, []);

//   return { data, onDataChange };
// };

type NameDescType = {
  name: string;
  description: string;
};

type K8sNameDescriptionFieldProps = {
  data: NameDescType;
  onDataChange: (data: NameDescType) => void;
  dataTestId: string;
  descriptionLabel?: string;
  nameLabel?: string;
  nameHelperText?: React.ReactNode;
  hideDescription?: boolean;
};

/**
 * Use in place of any K8s Resource creation / edit.
 * @see useK8sNameDescriptionFieldData
 */
const K8sNameDescriptionField: React.FC<K8sNameDescriptionFieldProps> = ({
  data,
  onDataChange,
  dataTestId,
  descriptionLabel = 'Description',
  nameLabel = 'Name',
  nameHelperText,
  hideDescription,
}) => {
  const [showK8sField, setShowK8sField] = React.useState(false);
  const { isMUITheme } = useThemeContext();

  const nameInput = (
    <TextInput
      data-testid={`${dataTestId}-name`}
      id={`${dataTestId}-name`}
      name={`${dataTestId}-name`}
      value={data.name}
      onChange={(_e, value) => onDataChange({ ...data, name: value })}
      isRequired
    />
  );

  const nameFormGroup = (
    <>
      <FormGroup label={nameLabel} isRequired fieldId={`${dataTestId}-name`}>
        <FormFieldset component={nameInput} field="Name" />
      </FormGroup>
      {nameHelperText || !showK8sField ? (
        //  &&
        // !k8sName.state.immutable
        <HelperText>
          {nameHelperText && <HelperTextItem>{nameHelperText}</HelperTextItem>}
          {!showK8sField && (
            // !k8sName.state.immutable
            // &&
            <>
              {/* {k8sName.value && (
                <HelperTextItem>
                  The resource name will be <b>{k8sName.value}</b>.
                </HelperTextItem>
              )} */}
              <HelperTextItem>
                <Button
                  data-testid={`${dataTestId}-editResourceLink`}
                  variant="link"
                  isInline
                  onClick={() => setShowK8sField(true)}
                >
                  Edit resource name
                </Button>{' '}
                <ResourceNameDefinitionTooltip />
              </HelperTextItem>
            </>
          )}
        </HelperText>
      ) : null}
    </>
  );

  const descriptionTextInput = (
    <TextInput
      data-testid={`${dataTestId}-description`}
      id={`${dataTestId}-description`}
      name={`${dataTestId}-description`}
      type="text"
      value={data.description}
      onChange={(_e, value) => onDataChange({ ...data, description: value })}
    />
  );

  const descriptionTextArea = (
    <TextArea
      data-testid={`${dataTestId}-description`}
      id={`${dataTestId}-description`}
      name={`${dataTestId}-description`}
      type="text"
      value={data.description}
      onChange={(_e, value) => onDataChange({ ...data, description: value })}
      resizeOrientation="vertical"
    />
  );

  const descriptionFormGroup = (
    <FormGroup label={descriptionLabel} fieldId={`${dataTestId}-description`}>
      <FormFieldset component={descriptionTextInput} field="Description" />
    </FormGroup>
  );

  return (
    <>
      {isMUITheme ? (
        nameFormGroup
      ) : (
        <>
          <FormGroup label={nameLabel} isRequired fieldId={`${dataTestId}-name`}>
            {nameInput}
            {nameHelperText}
          </FormGroup>
        </>
      )}

      <ResourceNameField allowEdit={showK8sField} dataTestId={dataTestId} />

      {!hideDescription && isMUITheme ? (
        descriptionFormGroup
      ) : (
        <FormGroup label={descriptionLabel} fieldId={`${dataTestId}-description`}>
          {descriptionTextArea}
        </FormGroup>
      )}
    </>
  );
};

export default K8sNameDescriptionField;
