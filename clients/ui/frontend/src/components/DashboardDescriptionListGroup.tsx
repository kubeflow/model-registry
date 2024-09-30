import * as React from 'react';
import {
  ActionList,
  ActionListItem,
  Button,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Flex,
  FlexItem,
  Split,
  SplitItem,
} from '@patternfly/react-core';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import { CheckIcon, PencilAltIcon, TimesIcon } from '@patternfly/react-icons';

import '~/components/DashboardDescriptionListGroup.scss';

type EditableProps = {
  isEditing: boolean;
  contentWhenEditing: React.ReactNode;
  isSavingEdits?: boolean;
  onEditClick: () => void;
  onSaveEditsClick: () => void;
  onDiscardEditsClick: () => void;
};

export type DashboardDescriptionListGroupProps = {
  title: React.ReactNode;
  tooltip?: React.ReactNode;
  action?: React.ReactNode;
  isEmpty?: boolean;
  contentWhenEmpty?: React.ReactNode;
  children: React.ReactNode;
} & (({ isEditable: true } & EditableProps) | ({ isEditable?: false } & Partial<EditableProps>));

const DashboardDescriptionListGroup: React.FC<DashboardDescriptionListGroupProps> = (props) => {
  const {
    title,
    tooltip,
    action,
    isEmpty,
    contentWhenEmpty,
    isEditable = false,
    isEditing,
    contentWhenEditing,
    isSavingEdits = false,
    onEditClick,
    onSaveEditsClick,
    onDiscardEditsClick,
    children,
  } = props;
  return (
    <DescriptionListGroup>
      {action || isEditable ? (
        <DescriptionListTerm className="kubeflow-custom-description-list-term-with-action">
          <Split>
            <SplitItem isFilled>{title}</SplitItem>
            <SplitItem>
              {action ||
                (isEditing ? (
                  <ActionList isIconList>
                    <ActionListItem>
                      <Button
                        data-testid={`save-edit-button-${title}`}
                        aria-label={`Save edits to ${title}`}
                        variant="link"
                        onClick={onSaveEditsClick}
                        isDisabled={isSavingEdits}
                      >
                        <CheckIcon />
                      </Button>
                    </ActionListItem>
                    <ActionListItem>
                      <Button
                        data-testid={`discard-edit-button-${title}`}
                        aria-label={`Discard edits to ${title} `}
                        variant="plain"
                        onClick={onDiscardEditsClick}
                        isDisabled={isSavingEdits}
                      >
                        <TimesIcon />
                      </Button>
                    </ActionListItem>
                  </ActionList>
                ) : (
                  <Button
                    data-testid={`edit-button-${title}`}
                    aria-label={`Edit ${title}`}
                    isInline
                    variant="link"
                    icon={<PencilAltIcon />}
                    iconPosition="end"
                    onClick={onEditClick}
                  >
                    Edit
                  </Button>
                ))}
            </SplitItem>
          </Split>
        </DescriptionListTerm>
      ) : (
        <DescriptionListTerm>
          <Flex
            spaceItems={{ default: 'spaceItemsSm' }}
            alignItems={{ default: 'alignItemsCenter' }}
          >
            <FlexItem>{title}</FlexItem>
            {tooltip}
          </Flex>
        </DescriptionListTerm>
      )}
      <DescriptionListDescription
        className={isEmpty && !isEditing ? text.textColorDisabled : ''}
        aria-disabled={!!(isEmpty && !isEditing)}
      >
        {isEditing ? contentWhenEditing : isEmpty ? contentWhenEmpty : children}
      </DescriptionListDescription>
    </DescriptionListGroup>
  );
};

export default DashboardDescriptionListGroup;
