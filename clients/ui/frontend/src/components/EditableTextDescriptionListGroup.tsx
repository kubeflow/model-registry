import * as React from 'react';
import { ExpandableSection, TextArea } from '@patternfly/react-core';
import DashboardDescriptionListGroup, {
  DashboardDescriptionListGroupProps,
} from '~/components/DashboardDescriptionListGroup';

type EditableTextDescriptionListGroupProps = Pick<
  DashboardDescriptionListGroupProps,
  'title' | 'contentWhenEmpty'
> & {
  value: string;
  saveEditedValue: (value: string) => Promise<void>;
  testid?: string;
};

const EditableTextDescriptionListGroup: React.FC<EditableTextDescriptionListGroupProps> = ({
  title,
  contentWhenEmpty,
  value,
  saveEditedValue,
  testid,
}) => {
  const [isEditing, setIsEditing] = React.useState(false);
  const [unsavedValue, setUnsavedValue] = React.useState(value);
  const [isSavingEdits, setIsSavingEdits] = React.useState(false);
  const [isTextExpanded, setIsTextExpanded] = React.useState(false);
  return (
    <DashboardDescriptionListGroup
      title={title}
      isEmpty={!value}
      contentWhenEmpty={contentWhenEmpty}
      isEditable
      isEditing={isEditing}
      isSavingEdits={isSavingEdits}
      contentWhenEditing={
        <TextArea
          data-testid={`edit-text-area-${title}`}
          aria-label={`Text box for editing ${title}`}
          value={unsavedValue}
          onChange={(_event, v) => setUnsavedValue(v)}
          isDisabled={isSavingEdits}
          rows={24}
          resizeOrientation="vertical"
        />
      }
      onEditClick={() => {
        setUnsavedValue(value);
        setIsEditing(true);
      }}
      onSaveEditsClick={async () => {
        setIsSavingEdits(true);
        try {
          await saveEditedValue(unsavedValue);
        } finally {
          setIsSavingEdits(false);
        }
        setIsEditing(false);
      }}
      onDiscardEditsClick={() => {
        setUnsavedValue(value);
        setIsEditing(false);
      }}
    >
      <ExpandableSection
        data-testid={testid}
        variant="truncate"
        truncateMaxLines={12}
        toggleText={isTextExpanded ? 'Show less' : 'Show more'}
        onToggle={(_event, isExpanded) => setIsTextExpanded(isExpanded)}
        isExpanded={isTextExpanded}
      >
        {value}
      </ExpandableSection>
    </DashboardDescriptionListGroup>
  );
};

export default EditableTextDescriptionListGroup;
