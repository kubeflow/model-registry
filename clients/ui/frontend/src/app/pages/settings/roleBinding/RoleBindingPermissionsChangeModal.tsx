import * as React from 'react';
import DeleteModal from '~/app/shared/components/DeleteModal';

type RoleBindingPermissionsChangeModalProps = {
  onClose: () => void;
  onEdit: () => void;
  onDelete: () => void;
  isDeleting: boolean;
  roleName?: string;
};

const RoleBindingPermissionsChangeModal: React.FC<RoleBindingPermissionsChangeModalProps> = ({
  onClose,
  onEdit,
  onDelete,
  isDeleting,
  roleName,
}) => {
  const textToShow = isDeleting ? 'Delete' : 'Edit';
  const [submitted, setSubmitted] = React.useState(false);
  return (
    <DeleteModal
      title={`Confirm ${textToShow.toLowerCase()}`}
      onClose={onClose}
      deleting={submitted}
      onDelete={() => {
        setSubmitted(true);
        if (isDeleting) {
          onDelete();
        } else {
          onEdit();
        }
      }}
      deleteName={roleName || 'delete'}
      submitButtonLabel={isDeleting ? 'Delete' : 'Save'}
      genericLabel
    >
      Are you sure you want to {isDeleting ? 'delete' : 'edit'} permissions for{' '}
      <strong>{roleName || 'this role binding'}</strong>? {isDeleting ? 'Deleting' : 'Editing'}{' '}
      these permissions may result in loss of access to this resource.
    </DeleteModal>
  );
};

export default RoleBindingPermissionsChangeModal;
