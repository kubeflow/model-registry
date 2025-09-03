import React from 'react';
import { Badge, Button, ExpandableSection } from '@patternfly/react-core';
import { AddCircleOIcon } from '@patternfly/react-icons';
import { Table, Tbody, Th, Thead, Tr } from '@patternfly/react-table';
import spacing from '@patternfly/react-styles/css/utilities/Spacing/spacing';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import ModelPropertiesTableRow from '~/app/pages/modelRegistry/screens/components/ModelPropertiesTableRow';
import { getProperties, mergeUpdatedProperty } from '~/app/pages/modelRegistry/screens/utils';
import { ModelRegistryCustomProperties } from '~/app/types';

type ModelPropertiesExpandableSectionProps = {
  customProperties?: ModelRegistryCustomProperties;
  isArchive?: boolean;
  saveEditedCustomProperties: (properties: ModelRegistryCustomProperties) => Promise<unknown>;
  isExpandedByDefault?: boolean;
};

const ModelPropertiesExpandableSection: React.FC<ModelPropertiesExpandableSectionProps> = ({
  customProperties = {},
  isArchive,
  saveEditedCustomProperties,
  isExpandedByDefault = false,
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

  const [isExpanded, setIsExpanded] = React.useState(isExpandedByDefault);

  return (
    <ExpandableSection
      isExpanded={isExpanded}
      onToggle={() => setIsExpanded(!isExpanded)}
      toggleContent={
        <>
          Properties <Badge isRead>{keys.length}</Badge>
        </>
      }
    >
      {keys.length > 0 && (
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
                keyValuePair={{ key, value: filteredProperties[key].string_value || '' }}
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
      )}
      {!isArchive && (
        <Button
          variant="link"
          data-testid="add-property-button"
          icon={<AddCircleOIcon />}
          isDisabled={isAdding || isSavingEdits}
          onClick={() => {
            setIsShowingMoreProperties(true);
            setIsAdding(true);
          }}
        >
          Add property
        </Button>
      )}
      {needExpandControl && (
        <Button
          variant="link"
          className={spacing.mtSm}
          data-testid="expand-control-button"
          onClick={() => setIsShowingMoreProperties(!isShowingMoreProperties)}
        >
          {isShowingMoreProperties
            ? 'Show fewer properties'
            : `Show ${numHiddenKeys} more ${numHiddenKeys === 1 ? 'property' : 'properties'}`}
        </Button>
      )}
    </ExpandableSection>
  );
};

export default ModelPropertiesExpandableSection;
