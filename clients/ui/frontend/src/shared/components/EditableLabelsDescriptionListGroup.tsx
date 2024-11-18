import * as React from 'react';
import {
  Button,
  Form,
  FormGroup,
  FormHelperText,
  HelperText,
  HelperTextItem,
  Label,
  LabelGroup,
  TextInput,
} from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import { ExclamationCircleIcon } from '@patternfly/react-icons';
import DashboardDescriptionListGroup, {
  DashboardDescriptionListGroupProps,
} from '~/shared/components/DashboardDescriptionListGroup';

type EditableTextDescriptionListGroupProps = Partial<
  Pick<DashboardDescriptionListGroupProps, 'title' | 'contentWhenEmpty'>
> & {
  labels: string[];
  saveEditedLabels: (labels: string[]) => Promise<unknown>;
  allExistingKeys?: string[];
  isArchive?: boolean;
};

const EditableLabelsDescriptionListGroup: React.FC<EditableTextDescriptionListGroupProps> = ({
  title = 'Labels',
  contentWhenEmpty = 'No labels',
  labels,
  saveEditedLabels,
  isArchive,
  allExistingKeys = labels,
}) => {
  const [isEditing, setIsEditing] = React.useState(false);
  const [unsavedLabels, setUnsavedLabels] = React.useState(labels);
  const [isSavingEdits, setIsSavingEdits] = React.useState(false);

  const editUnsavedLabel = (newText: string, index: number) => {
    if (isSavingEdits) {
      return;
    }
    const copy = [...unsavedLabels];
    copy[index] = newText;
    setUnsavedLabels(copy);
  };
  const removeUnsavedLabel = (text: string) => {
    if (isSavingEdits) {
      return;
    }
    setUnsavedLabels(unsavedLabels.filter((label) => label !== text));
  };
  const addUnsavedLabel = (text: string) => {
    if (isSavingEdits) {
      return;
    }
    setUnsavedLabels([...unsavedLabels, text]);
  };

  // Don't allow a label that matches a non-label property key or another label (as they stand before saving)
  // Note that this means if you remove a label and add it back before saving, that is valid
  const reservedKeys = [
    ...allExistingKeys.filter((key) => !labels.includes(key)),
    ...unsavedLabels,
  ];

  const [isAddLabelModalOpen, setIsAddLabelModalOpen] = React.useState(false);
  const [addLabelInputValue, setAddLabelInputValue] = React.useState('');
  const addLabelInputRef = React.useRef<HTMLInputElement>(null);
  let addLabelValidationError: string | null = null;
  if (reservedKeys.includes(addLabelInputValue)) {
    addLabelValidationError = 'Label must not match an existing label or property key';
  } else if (addLabelInputValue.length > 63) {
    addLabelValidationError = "Label text can't exceed 63 characters";
  }

  const toggleAddLabelModal = () => {
    setAddLabelInputValue('');
    setIsAddLabelModalOpen(!isAddLabelModalOpen);
  };
  React.useEffect(() => {
    if (isAddLabelModalOpen && addLabelInputRef.current) {
      addLabelInputRef.current.focus();
    }
  }, [isAddLabelModalOpen]);

  const addLabelModalSubmitDisabled = !addLabelInputValue || !!addLabelValidationError;
  const submitAddLabelModal = (event?: React.FormEvent) => {
    event?.preventDefault();
    if (!addLabelModalSubmitDisabled) {
      addUnsavedLabel(addLabelInputValue);
      toggleAddLabelModal();
    }
  };

  return (
    <>
      <DashboardDescriptionListGroup
        title={title}
        isEmpty={labels.length === 0}
        contentWhenEmpty={contentWhenEmpty}
        isEditable={!isArchive}
        isEditing={isEditing}
        isSavingEdits={isSavingEdits}
        contentWhenEditing={
          <LabelGroup
            data-testid="label-group"
            isEditable={!isSavingEdits}
            numLabels={unsavedLabels.length}
            addLabelControl={
              !isSavingEdits && (
                <Label
                  textMaxWidth="40ch"
                  color="blue"
                  variant="overflow"
                  onClick={toggleAddLabelModal}
                >
                  Add label
                </Label>
              )
            }
          >
            {unsavedLabels.map((label, index) => (
              <Label
                key={label}
                color="blue"
                data-testid="label"
                isEditable={!isSavingEdits}
                editableProps={{ 'aria-label': `Editable label with text ${label}` }}
                onClose={() => removeUnsavedLabel(label)}
                closeBtnProps={{ isDisabled: isSavingEdits }}
                onEditComplete={(_event, newText) => {
                  if (!reservedKeys.includes(newText) && newText.length <= 63) {
                    editUnsavedLabel(newText, index);
                  }
                }}
              >
                {label}
              </Label>
            ))}
          </LabelGroup>
        }
        onEditClick={() => {
          setUnsavedLabels(labels);
          setIsEditing(true);
        }}
        onSaveEditsClick={async () => {
          setIsSavingEdits(true);
          try {
            await saveEditedLabels(unsavedLabels);
          } finally {
            setIsSavingEdits(false);
          }
          setIsEditing(false);
        }}
        onDiscardEditsClick={() => {
          setUnsavedLabels(labels);
          setIsEditing(false);
        }}
      >
        <LabelGroup data-testid="label-group">
          {labels.map((label) => (
            <Label textMaxWidth="40ch" key={label} color="blue" data-testid="label">
              {label}
            </Label>
          ))}
        </LabelGroup>
      </DashboardDescriptionListGroup>
      <Modal
        variant="small"
        title="Add label"
        isOpen={isAddLabelModalOpen}
        onClose={toggleAddLabelModal}
        actions={[
          <Button
            key="save"
            variant="primary"
            form="add-label-form"
            onClick={submitAddLabelModal}
            isDisabled={addLabelModalSubmitDisabled}
          >
            Save
          </Button>,
          <Button key="cancel" variant="link" onClick={toggleAddLabelModal}>
            Cancel
          </Button>,
        ]}
      >
        <Form id="add-label-form" onSubmit={submitAddLabelModal}>
          <FormGroup label="Label text" fieldId="add-label-form-label-text" isRequired>
            <TextInput
              type="text"
              id="add-label-form-label-text"
              name="add-label-form-label-text"
              value={addLabelInputValue}
              onChange={(_event: React.FormEvent<HTMLInputElement>, value: string) =>
                setAddLabelInputValue(value)
              }
              ref={addLabelInputRef}
              isRequired
              validated={addLabelValidationError ? 'error' : 'default'}
            />
            {addLabelValidationError && (
              <FormHelperText>
                <HelperText>
                  <HelperTextItem icon={<ExclamationCircleIcon />} variant="error">
                    {addLabelValidationError}
                  </HelperTextItem>
                </HelperText>
              </FormHelperText>
            )}
          </FormGroup>
        </Form>
      </Modal>
    </>
  );
};

export default EditableLabelsDescriptionListGroup;
