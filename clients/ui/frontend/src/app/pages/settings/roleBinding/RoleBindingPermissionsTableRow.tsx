import * as React from 'react';
import { ActionsColumn, Tbody, Td, Tr } from '@patternfly/react-table';
import {
  Button,
  Split,
  SplitItem,
  Timestamp,
  TimestampTooltipVariant,
  Truncate,
  Skeleton,
} from '@patternfly/react-core';
import { CheckIcon, TimesIcon, EllipsisVIcon } from '@patternfly/react-icons';
import {
  DashboardPopupIconButton,
  relativeTime,
  RoleBindingKind,
  RoleBindingSubject,
} from 'mod-arch-shared';
import useUser from '~/app/hooks/useUser';
import { useNamespaces } from '~/app/hooks/useNamespaces';
import { NamespaceKind } from '~/app/shared/components/types';
import {
  castRoleBindingPermissionsRoleType,
  roleLabel,
  isCurrentUserChanging,
  displayNameToNamespace,
} from './utils';
import { RoleBindingPermissionsRoleType } from './types';
import { RoleBindingPermissionsNameInput } from './RoleBindingPermissionsNameInput';
import RoleBindingPermissionsPermissionSelection from './RoleBindingPermissionsPermissionSelection';
import RoleBindingPermissionsChangeModal from './RoleBindingPermissionsChangeModal';

type RoleBindingPermissionsTableRowProps = {
  roleBindingObject?: RoleBindingKind;
  subjectKind: RoleBindingSubject['kind'];
  defaultRoleBindingName?: string;
  permissionOptions: {
    type: RoleBindingPermissionsRoleType;
    description: string;
  }[];
  typeAhead?: string[];
  isProjectSubject?: boolean;
  isEditing?: boolean;
  isAdding?: boolean;
  onChange?: (subjectName: string, roleRefName: string) => void;
  onCancel?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
};

const defaultValueName = (
  obj: RoleBindingKind,
  isProjectSubject?: boolean,
  namespaces?: NamespaceKind[],
): string =>
  isProjectSubject && namespaces
    ? namespaces.find(
        (ns) => ns.name === obj.subjects[0]?.name.replace(/^system:serviceaccounts:/, ''),
      )?.displayName ||
      obj.subjects[0]?.name ||
      ''
    : obj.subjects[0]?.name || '';

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
  isProjectSubject,
  onChange,
  onCancel,
  onEdit,
  onDelete,
}) => {
  const [namespaces, namespacesLoaded] = useNamespaces();
  const currentUser = useUser();
  const isCurrentUserBeingChanged = isCurrentUserChanging(obj, currentUser.userId);
  const [roleBindingName, setRoleBindingName] = React.useState(() => {
    if (isAdding || !obj) {
      return '';
    }
    return defaultValueName(obj, isProjectSubject, namespaces);
  });
  const [roleBindingRoleRef, setRoleBindingRoleRef] =
    React.useState<RoleBindingPermissionsRoleType>(() => {
      if (isAdding || !obj) {
        return permissionOptions[0]?.type;
      }
      return defaultValueRole(obj);
    });
  const [isLoading, setIsLoading] = React.useState(false);
  const createdDate = obj?.metadata.creationTimestamp
    ? new Date(obj.metadata.creationTimestamp)
    : null;
  const isDefaultGroup = obj?.metadata.name === defaultRoleBindingName;
  const [showModal, setShowModal] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);

  // Update name when namespaces load or change
  React.useEffect(() => {
    if (namespacesLoaded && obj && !isAdding && isProjectSubject) {
      setRoleBindingName(defaultValueName(obj, isProjectSubject, namespaces));
    }
  }, [namespacesLoaded, obj, isAdding, isProjectSubject, namespaces]);

  //Sync local state with props if exiting edit mode
  React.useEffect(() => {
    if (!isEditing && obj) {
      setRoleBindingName(
        isProjectSubject
          ? defaultValueName(obj, isProjectSubject, namespaces)
          : defaultValueName(obj),
      );
      setRoleBindingRoleRef(defaultValueRole(obj));
    }
  }, [obj, isEditing, isProjectSubject, namespaces]);

  const showLoadingSkeleton = isProjectSubject && !namespacesLoaded;

  return (
    <>
      <Tbody>
        <Tr>
          <Td dataLabel="Username">
            {isEditing || isAdding ? (
              showLoadingSkeleton ? (
                <Skeleton height="36px" />
              ) : (
                <RoleBindingPermissionsNameInput
                  subjectKind={subjectKind}
                  value={roleBindingName}
                  onChange={(selection) => setRoleBindingName(selection)}
                  onClear={() => setRoleBindingName('')}
                  placeholderText={`Type ${
                    isProjectSubject
                      ? 'project name'
                      : subjectKind === 'Group'
                        ? 'group name'
                        : 'username'
                  }`}
                  typeAhead={typeAhead}
                  isProjectSubject={isProjectSubject}
                />
              )
            ) : showLoadingSkeleton ? (
              <Skeleton height="20px" />
            ) : (
              <Truncate content={roleBindingName} />
            )}
          </Td>
          <Td dataLabel="Permission">
            {isEditing || isAdding ? (
              <RoleBindingPermissionsPermissionSelection
                permissionOptions={permissionOptions}
                selection={roleBindingRoleRef}
                onSelect={(roleType) => setRoleBindingRoleRef(roleType)}
              />
            ) : (
              roleLabel(roleBindingRoleRef)
            )}
          </Td>
          <Td dataLabel="Date added">
            {createdDate && !isAdding ? (
              <Timestamp
                date={createdDate}
                tooltip={{
                  variant: TimestampTooltipVariant.default,
                }}
              >
                {relativeTime(Date.now(), createdDate.getTime())}
              </Timestamp>
            ) : null}
          </Td>
          <Td isActionCell>
            {isEditing || isAdding ? (
              <Split hasGutter>
                <SplitItem>
                  <Button
                    data-testid="save-rolebinding-button"
                    aria-label="Save role binding edits"
                    variant="link"
                    icon={<CheckIcon />}
                    onClick={async () => {
                      if (!roleBindingName) {
                        return;
                      }
                      setIsLoading(true);
                      try {
                        await onChange?.(
                          isProjectSubject
                            ? `system:serviceaccounts:${displayNameToNamespace(
                                roleBindingName,
                                namespaces,
                              )}`
                            : roleBindingName,
                          roleBindingRoleRef,
                        );
                      } finally {
                        setIsLoading(false);
                      }
                    }}
                    isLoading={isLoading}
                    isDisabled={isLoading || !roleBindingName || showLoadingSkeleton}
                  />
                </SplitItem>
                <SplitItem>
                  <Button
                    aria-label="Cancel role binding edits"
                    variant="plain"
                    icon={<TimesIcon />}
                    onClick={() => {
                      onCancel?.();
                    }}
                    isDisabled={isLoading}
                  />
                </SplitItem>
              </Split>
            ) : (
              <ActionsColumn
                items={[
                  {
                    title: 'Edit',
                    onClick: () => {
                      if (isCurrentUserBeingChanged) {
                        setIsDeleting(false);
                        setShowModal(true);
                      } else {
                        onEdit?.();
                      }
                    },
                    isDisabled: showLoadingSkeleton,
                  },
                  {
                    isSeparator: true,
                  },
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
                    isDisabled: isDefaultGroup || showLoadingSkeleton,
                  },
                ]}
                rowData={{}}
                actionsToggle={(props) => (
                  <DashboardPopupIconButton
                    icon={<EllipsisVIcon />}
                    aria-label="Role binding table row actions"
                    {...props}
                  />
                )}
              />
            )}
          </Td>
        </Tr>
      </Tbody>
      {showModal && (
        <RoleBindingPermissionsChangeModal
          onClose={() => setShowModal(false)}
          onEdit={() => {
            setShowModal(false);
            const finalSubjectName = isProjectSubject
              ? `system:serviceaccounts:${displayNameToNamespace(roleBindingName, namespaces)}`
              : roleBindingName;
            onChange?.(finalSubjectName, roleBindingRoleRef);
          }}
          onDelete={() => {
            setShowModal(false);
            onDelete?.();
          }}
          isDeleting={isDeleting}
          roleName={roleBindingName}
        />
      )}
    </>
  );
};

export default RoleBindingPermissionsTableRow;
