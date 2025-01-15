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
  Popover,
  Split,
  SplitItem,
} from '@patternfly/react-core';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import {
  CheckIcon,
  PencilAltIcon,
  TimesIcon,
  OutlinedQuestionCircleIcon,
} from '@patternfly/react-icons';

import './DashboardDescriptionListGroup.scss';
import DashboardPopupIconButton from '~/shared/components/dashboard/DashboardPopupIconButton';

type EditableProps = {
  isEditing: boolean;
  contentWhenEditing: React.ReactNode;
  isSavingEdits?: boolean;
  onEditClick: () => void;
  onSaveEditsClick: () => void;
  onDiscardEditsClick: () => void;
  editButtonTestId?: string;
  saveButtonTestId?: string;
  cancelButtonTestId?: string;
  discardButtonTestId?: string;
};

export type DashboardDescriptionListGroupProps = {
  title: string;
  popover?: React.ReactNode;
  action?: React.ReactNode;
  isEmpty?: boolean;
  contentWhenEmpty?: React.ReactNode;
  children: React.ReactNode;
  groupTestId?: string;
  isSaveDisabled?: boolean;
} & (({ isEditable: true } & EditableProps) | ({ isEditable?: false } & Partial<EditableProps>));

const DashboardDescriptionListGroup: React.FC<DashboardDescriptionListGroupProps> = (props) => {
  const {
    title,
    popover,
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
    groupTestId,
    editButtonTestId,
    saveButtonTestId,
    cancelButtonTestId,
    isSaveDisabled,
  } = props;
  return (
    <DescriptionListGroup data-testid={groupTestId}>
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
                        data-testid={saveButtonTestId}
                        icon={<CheckIcon />}
                        aria-label={`Save edits to ${title}`}
                        variant="link"
                        onClick={onSaveEditsClick}
                        isDisabled={isSavingEdits || isSaveDisabled}
                      />
                    </ActionListItem>
                    <ActionListItem>
                      <Button
                        data-testid={cancelButtonTestId}
                        icon={<TimesIcon />}
                        aria-label={`Discard edits to ${title} `}
                        variant="plain"
                        onClick={onDiscardEditsClick}
                        isDisabled={isSavingEdits}
                      />
                    </ActionListItem>
                  </ActionList>
                ) : (
                  <Button
                    data-testid={editButtonTestId}
                    aria-label={`Edit ${title}`}
                    variant="link"
                    icon={<PencilAltIcon />}
                    iconPosition="start"
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
            spaceItems={{ default: 'spaceItemsNone' }}
            alignItems={{ default: 'alignItemsCenter' }}
          >
            <FlexItem>{title}</FlexItem>
            {popover && (
              <Popover bodyContent={popover}>
                <DashboardPopupIconButton
                  icon={<OutlinedQuestionCircleIcon />}
                  aria-label="More info"
                />
              </Popover>
            )}
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
