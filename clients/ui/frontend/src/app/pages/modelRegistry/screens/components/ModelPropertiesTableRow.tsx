import * as React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import {
  ActionList,
  ActionListItem,
  Button,
  ExpandableSection,
  FormHelperText,
  HelperText,
  HelperTextItem,
  TextInput,
  Truncate,
} from '@patternfly/react-core';
import { CheckIcon, ExternalLinkAltIcon, TimesIcon } from '@patternfly/react-icons';
import { KeyValuePair, EitherNotBoth } from 'mod-arch-core';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import { isValidHttpUrl } from '~/app/pages/modelRegistry/screens/utils';

type ModelPropertiesTableRowProps = {
  allExistingKeys: string[];
  setIsEditing: (isEditing: boolean) => void;
  isSavingEdits: boolean;
  isArchive?: boolean;
  setIsSavingEdits: (isSaving: boolean) => void;
  saveEditedProperty: (oldKey: string, newPair: KeyValuePair) => Promise<unknown>;
} & EitherNotBoth<
  { isAddRow: true },
  {
    isEditing: boolean;
    keyValuePair: KeyValuePair;
    deleteProperty: (key: string) => Promise<unknown>;
  }
>;

const ModelPropertiesTableRow: React.FC<ModelPropertiesTableRowProps> = ({
  isAddRow,
  isEditing = isAddRow,
  keyValuePair = { key: '', value: '' },
  deleteProperty = () => Promise.resolve(),
  allExistingKeys,
  setIsEditing,
  isSavingEdits,
  setIsSavingEdits,
  isArchive,
  saveEditedProperty,
}) => {
  const { key, value } = keyValuePair;

  const [unsavedKey, setUnsavedKey] = React.useState(key);
  const [unsavedValue, setUnsavedValue] = React.useState(value);

  const [isValueExpanded, setIsValueExpanded] = React.useState(true);

  let keyValidationError: string | null = null;
  if (unsavedKey !== key && allExistingKeys.includes(unsavedKey)) {
    keyValidationError = 'Key must not match an existing property key or label';
  } else if (unsavedKey.length > 63) {
    keyValidationError = "Key text can't exceed 63 characters";
  }

  const clearUnsavedInputs = () => {
    setUnsavedKey(key);
    setUnsavedValue(value);
  };

  const onEditClick = () => {
    clearUnsavedInputs();
    setIsEditing(true);
  };

  const onDeleteClick = async () => {
    setIsSavingEdits(true);
    try {
      await deleteProperty(key);
    } finally {
      setIsSavingEdits(false);
    }
  };

  const onSaveEditsClick = async () => {
    setIsSavingEdits(true);
    try {
      await saveEditedProperty(key, { key: unsavedKey, value: unsavedValue });
    } finally {
      setIsSavingEdits(false);
    }
    setIsEditing(false);
  };

  const onDiscardEditsClick = () => {
    clearUnsavedInputs();
    setIsEditing(false);
  };

  const propertyKeyInput = (
    <TextInput
      data-testid={isAddRow ? `add-property-key-input` : `edit-property-key-input ${key}`}
      aria-label={
        isAddRow ? 'Key input for new property' : `Key input for editing property with key ${key}`
      }
      isRequired
      type="text"
      autoFocus
      value={unsavedKey}
      onChange={(_event, str) => setUnsavedKey(str)}
      validated={keyValidationError ? 'error' : 'default'}
    />
  );

  const propertyValueInput = (
    <TextInput
      data-testid={isAddRow ? `add-property-value-input` : `edit-property-value-input ${value}`}
      aria-label={
        isAddRow
          ? 'Value input for new property'
          : `Value input for editing property with key ${key}`
      }
      isRequired
      type="text"
      value={unsavedValue}
      onChange={(_event, str) => setUnsavedValue(str)}
    />
  );

  return (
    <Tr>
      <Td dataLabel="Key" width={45} modifier="breakWord">
        {isEditing ? (
          <>
            <FormFieldset className="tr-fieldset-wrapper" component={propertyKeyInput} />

            {keyValidationError && (
              <FormHelperText>
                <HelperText>
                  <HelperTextItem variant="error">{keyValidationError}</HelperTextItem>
                </HelperText>
              </FormHelperText>
            )}
          </>
        ) : (
          key
        )}
      </Td>
      <Td dataLabel="Value" width={45} modifier="breakWord">
        {isEditing ? (
          <FormFieldset className="tr-fieldset-wrapper" component={propertyValueInput} />
        ) : (
          <ExpandableSection
            variant="truncate"
            truncateMaxLines={3}
            toggleText={isValueExpanded ? 'Show less' : 'Show more'}
            onToggle={(_event, isExpanded) => setIsValueExpanded(isExpanded)}
            isExpanded={isValueExpanded}
          >
            {isValidHttpUrl(value) ? (
              <Button
                variant="link"
                icon={<ExternalLinkAltIcon />}
                iconPosition="end"
                component="a"
                href={value}
                target="_blank"
                isInline
              >
                <Truncate content={value} />
              </Button>
            ) : (
              value
            )}
          </ExpandableSection>
        )}
      </Td>
      {!isArchive && (
        <Td isActionCell width={10}>
          {isEditing ? (
            <ActionList isIconList>
              <ActionListItem>
                <Button
                  data-testid="save-edit-button-property"
                  icon={<CheckIcon />}
                  aria-label={`Save edits to property with key ${key}`}
                  variant="link"
                  onClick={onSaveEditsClick}
                  isDisabled={isSavingEdits || !unsavedKey || !unsavedValue || !!keyValidationError}
                />
              </ActionListItem>
              <ActionListItem>
                <Button
                  data-testid="discard-edit-button-property"
                  icon={<TimesIcon />}
                  aria-label={`Discard edits to property with key ${key}`}
                  variant="plain"
                  onClick={onDiscardEditsClick}
                  isDisabled={isSavingEdits}
                />
              </ActionListItem>
            </ActionList>
          ) : (
            <ActionsColumn
              isDisabled={isSavingEdits}
              popperProps={{ direction: 'up' }}
              items={[
                { title: 'Edit', onClick: onEditClick, isDisabled: isSavingEdits },
                { isSeparator: true },
                { title: 'Delete', onClick: onDeleteClick, isDisabled: isSavingEdits },
              ]}
            />
          )}
        </Td>
      )}
    </Tr>
  );
};

export default ModelPropertiesTableRow;
