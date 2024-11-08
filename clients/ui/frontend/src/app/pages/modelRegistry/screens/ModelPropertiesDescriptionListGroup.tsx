import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { Table, Tbody, Th, Thead, Tr } from '@patternfly/react-table';
import { PlusCircleIcon } from '@patternfly/react-icons';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import DashboardDescriptionListGroup from '~/shared/components/DashboardDescriptionListGroup';
import { getProperties, mergeUpdatedProperty } from '~/app/pages/modelRegistry/screens/utils';
import { ModelRegistryCustomProperties } from '~/app/types';
import ModelPropertiesTableRow from '~/app/pages/modelRegistry/screens/ModelPropertiesTableRow';

type ModelPropertiesDescriptionListGroupProps = {
  customProperties: ModelRegistryCustomProperties;
  isArchive?: boolean;
  saveEditedCustomProperties: (properties: ModelRegistryCustomProperties) => Promise<unknown>;
};

const ModelPropertiesDescriptionListGroup: React.FC<ModelPropertiesDescriptionListGroupProps> = ({
  customProperties = {},
  isArchive,
  saveEditedCustomProperties,
}) => {
  const [editingPropertyKeys, setEditingPropertyKeys] = React.useState<string[]>([]);
  const setIsEditingKey = (key: string, isEditing: boolean) =>
    setEditingPropertyKeys([
      ...editingPropertyKeys.filter((k) => k !== key),
      ...(isEditing ? [key] : []),
    ]);
  const [isAdding, setIsAdding] = React.useState(false);
  const isEditingSomeRow = isAdding || editingPropertyKeys.length > 0;

  const [isSavingEdits, setIsSavingEdits] = React.useState(false);

  // We only show string properties with a defined value (no labels or other property types)
  const filteredProperties = getProperties(customProperties);

  const [isShowingMoreProperties, setIsShowingMoreProperties] = React.useState(false);
  const keys = Object.keys(filteredProperties);
  const needExpandControl = keys.length > 5;
  const shownKeys = isShowingMoreProperties ? keys : keys.slice(0, 5);
  const numHiddenKeys = keys.length - shownKeys.length;

  // Includes keys reserved by non-string properties and labels
  const allExistingKeys = Object.keys(customProperties);

  const requiredAsterisk = (
    <span aria-hidden="true" className={text.textColorStatusDanger}>
      {' *'}
    </span>
  );

  return (
    <DashboardDescriptionListGroup
      title="Properties"
      action={
        !isArchive && (
          <Button
            isInline
            variant="link"
            icon={<PlusCircleIcon />}
            iconPosition="start"
            isDisabled={isAdding || isSavingEdits}
            onClick={() => setIsAdding(true)}
          >
            Add property
          </Button>
        )
      }
      isEmpty={!isAdding && keys.length === 0}
      contentWhenEmpty="No properties"
    >
      <Table aria-label="Properties table" data-testid="properties-table" variant="compact">
        <Thead>
          <Tr>
            <Th>Key {isEditingSomeRow && requiredAsterisk}</Th>
            <Th>Value {isEditingSomeRow && requiredAsterisk}</Th>
            <Th screenReaderText="Actions" />
          </Tr>
        </Thead>
        <Tbody>
          {shownKeys.map((key) => (
            <ModelPropertiesTableRow
              key={key}
              isArchive={isArchive}
              keyValuePair={{ key, value: filteredProperties[key].string_value }}
              allExistingKeys={allExistingKeys}
              isEditing={editingPropertyKeys.includes(key)}
              setIsEditing={(isEditing) => setIsEditingKey(key, isEditing)}
              isSavingEdits={isSavingEdits}
              setIsSavingEdits={setIsSavingEdits}
              saveEditedProperty={(oldKey, newPair) =>
                saveEditedCustomProperties(
                  mergeUpdatedProperty({ customProperties, op: 'update', oldKey, newPair }),
                )
              }
              deleteProperty={(oldKey) =>
                saveEditedCustomProperties(
                  mergeUpdatedProperty({ customProperties, op: 'delete', oldKey }),
                )
              }
            />
          ))}
          {isAdding && (
            <ModelPropertiesTableRow
              isAddRow
              allExistingKeys={allExistingKeys}
              setIsEditing={setIsAdding}
              isSavingEdits={isSavingEdits}
              setIsSavingEdits={setIsSavingEdits}
              saveEditedProperty={(_oldKey, newPair) =>
                saveEditedCustomProperties(
                  mergeUpdatedProperty({ customProperties, op: 'create', newPair }),
                )
              }
            />
          )}
        </Tbody>
      </Table>
      {needExpandControl && (
        <Button
          variant="link"
          className={spacing.mtSm}
          onClick={() => setIsShowingMoreProperties(!isShowingMoreProperties)}
        >
          {isShowingMoreProperties
            ? 'Show fewer properties'
            : `Show ${numHiddenKeys} more ${numHiddenKeys === 1 ? 'property' : 'properties'}`}
        </Button>
      )}
    </DashboardDescriptionListGroup>
  );
};

export default ModelPropertiesDescriptionListGroup;
