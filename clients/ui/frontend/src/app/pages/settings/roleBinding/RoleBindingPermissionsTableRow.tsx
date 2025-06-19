import * as React from 'react';
import { ActionsColumn, Tbody, Td, Tr } from '@patternfly/react-table';
import {
  Button,
  Popover,
  Split,
  SplitItem,
  Content,
  Timestamp,
  TimestampTooltipVariant,
  Tooltip,
  Truncate,
} from '@patternfly/react-core';
import {
  CheckIcon,
  OutlinedQuestionCircleIcon,
  TimesIcon,
  EllipsisVIcon,
} from '@patternfly/react-icons';
import {
  DashboardPopupIconButton,
  relativeTime,
  RoleBindingKind,
  RoleBindingSubject,
} from 'mod-arch-shared';
import useUser from '~/app/hooks/useUser';
import {
  castRoleBindingPermissionsRoleType,
  firstSubject,
  roleLabel,
  isCurrentUserChanging,
} from './utils';
import { RoleBindingPermissionsRoleType } from './types';
import RoleBindingPermissionsNameInput from './RoleBindingPermissionsNameInput';
import RoleBindingPermissionsPermissionSelection from './RoleBindingPermissionsPermissionSelection';
import RoleBindingPermissionsChangeModal from './RoleBindingPermissionsChangeModal';

type RoleBindingPermissionsTableRowProps = {
  roleBindingObject?: RoleBindingKind;
  subjectKind: RoleBindingSubject['kind'];
  isEditing: boolean;
  isAdding: boolean;
  defaultRoleBindingName?: string;
  permissionOptions: {
    type: RoleBindingPermissionsRoleType;
    description: string;
  }[];
  typeAhead?: string[];
  onChange: (name: string, roleType: RoleBindingPermissionsRoleType) => void;
  onCancel: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
};

const defaultValueName = (obj: RoleBindingKind) => firstSubject(obj);
const defaultValueRole = (obj: RoleBindingKind) =>
  castRoleBindingPermissionsRoleType(obj.roleRef.name);

const RoleBindingPermissionsTableRow: React.FC<RoleBindingPermissionsTableRowProps> = ({
  roleBindingObject: obj,
  subjectKind,
  isEditing,
  isAdding,
  defaultRoleBindingName,
  permissionOptions,
  typeAhead,
  onChange,
  onCancel,
  onEdit,
  onDelete,
}) => {
  // TODO: We don't have project context yet and need to add logic to show projects permission tab under manage permissions of MR - might need to move the project-context part to shared library
  const currentUser = useUser();
  const isCurrentUserBeingChanged = isCurrentUserChanging(obj, currentUser.userId);
  const [roleBindingName, setRoleBindingName] = React.useState(() => {
    if (isAdding || !obj) {
      return '';
    }
    return defaultValueName(obj);
  });
  const [roleBindingRoleRef, setRoleBindingRoleRef] =
    React.useState<RoleBindingPermissionsRoleType>(() => {
      if (isAdding || !obj) {
        return permissionOptions[0]?.type;
      }
      return defaultValueRole(obj);
    });
  const [isLoading, setIsLoading] = React.useState(false);
  const createdDate = new Date(obj?.metadata.creationTimestamp ?? '');
  const isDefaultGroup = obj?.metadata.name === defaultRoleBindingName;
  const [showModal, setShowModal] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);

  //Sync local state with props if exiting edit mode
  React.useEffect(() => {
    if (!isEditing && obj) {
      setRoleBindingName(defaultValueName(obj));
      setRoleBindingRoleRef(defaultValueRole(obj));
    }
  }, [obj, isEditing]);

  return (
    <>
      <Tbody>
        <Tr>
          <Td dataLabel="Username">
            {isEditing || isAdding ? (
              <RoleBindingPermissionsNameInput
                subjectKind={subjectKind}
                value={roleBindingName}
                onChange={(selection: React.SetStateAction<string>) => {
                  setRoleBindingName(selection);
                }}
                onClear={() => setRoleBindingName('')}
                placeholderText="Select a group"
                typeAhead={typeAhead}
              />
            ) : (
              <Content component="p">
                <Truncate content={roleBindingName} />
                {` `}
                {isDefaultGroup && (
                  <Popover
                    bodyContent={
                      <div>
                        This group is created by default. You can add users to this group in
                        OpenShift user management, or ask the cluster admin to do so.
                      </div>
                    }
                  >
                    <DashboardPopupIconButton
                      icon={<OutlinedQuestionCircleIcon />}
                      aria-label="More info"
                    />
                  </Popover>
                )}
              </Content>
            )}
          </Td>
          <Td dataLabel="Permission">
            {(isEditing || isAdding) && permissionOptions.length > 1 ? (
              <RoleBindingPermissionsPermissionSelection
                permissionOptions={permissionOptions}
                selection={roleBindingRoleRef}
                onSelect={(selection) => {
                  setRoleBindingRoleRef(selection);
                }}
              />
            ) : (
              <Content component="p">{roleLabel(roleBindingRoleRef)}</Content>
            )}
          </Td>
          <Td dataLabel="Date added">
            {!isEditing && !isAdding && (
              <Content component="p">
                <Timestamp
                  date={createdDate}
                  tooltip={{ variant: TimestampTooltipVariant.default }}
                >
                  {relativeTime(Date.now(), createdDate.getTime())}
                </Timestamp>
              </Content>
            )}
          </Td>
          <Td isActionCell modifier="nowrap" style={{ textAlign: 'right' }}>
            {isEditing || isAdding ? (
              <Split>
                <SplitItem>
                  <Button
                    data-testid={isAdding ? `save-new-button` : `save-button ${roleBindingName}`}
                    data-id="save-rolebinding-button"
                    aria-label="Save role binding"
                    variant="link"
                    icon={<CheckIcon />}
                    isDisabled={isLoading || !roleBindingName || !roleBindingRoleRef}
                    onClick={() => {
                      if (isCurrentUserBeingChanged) {
                        setIsDeleting(false);
                        setShowModal(true);
                      } else {
                        setIsLoading(true);
                        onChange(roleBindingName, roleBindingRoleRef);
                        setIsLoading(false);
                      }
                    }}
                  />
                </SplitItem>
                <SplitItem>
                  <Button
                    data-id="cancel-rolebinding-button"
                    aria-label="Cancel role binding"
                    variant="plain"
                    isDisabled={isLoading}
                    icon={<TimesIcon />}
                    onClick={() => {
                      onCancel();
                    }}
                  />
                </SplitItem>
              </Split>
            ) : isDefaultGroup ? (
              <Tooltip content="The default group cannot be edited or deleted. The group's members can be managed via the API.">
                <Button
                  icon={<EllipsisVIcon />}
                  variant="plain"
                  isAriaDisabled
                  aria-label="The default group always has access to model registry."
                />
              </Tooltip>
            ) : (
              <ActionsColumn
                items={[
                  {
                    title: 'Edit',
                    onClick: () => {
                      onEdit?.();
                    },
                  },
                  { isSeparator: true },
                  {
                    title: 'Delete',
                    onClick: () => {
                      if (isCurrentUserBeingChanged) {
                        setIsDeleting(true);
                        setShowModal(true);
                      } else {
                        onDelete?.();
                      }
                    },
                  },
                ]}
              />
            )}
          </Td>
        </Tr>
      </Tbody>
      {showModal && (
        <RoleBindingPermissionsChangeModal
          roleName={currentUser.userId}
          onClose={() => {
            setShowModal(false);
            if (isEditing) {
              onCancel();
            }
          }}
          onEdit={() => {
            setIsLoading(true);
            onChange(roleBindingName, roleBindingRoleRef);
            setIsLoading(false);
            setShowModal(false);
          }}
          onDelete={() => onDelete?.()}
          isDeleting={isDeleting}
        />
      )}
    </>
  );
};

export default RoleBindingPermissionsTableRow;
