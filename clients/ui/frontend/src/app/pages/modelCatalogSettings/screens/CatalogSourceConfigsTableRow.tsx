import * as React from 'react';
import { ActionsColumn, Td, Tr } from '@patternfly/react-table';
import { Button, Label, Switch } from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { manageSourceUrl } from '~/app/routes/modelCatalogSettings/modelCatalogSettings';
import {
  CATALOG_SOURCE_TYPE_LABELS,
  ModelVisibilityBadgeColor,
} from '~/concepts/modelCatalogSettings/const';
import { hasSourceFilters, getOrganizationDisplay } from '~/concepts/modelCatalogSettings/utils';
import DeleteModal from '~/app/shared/components/DeleteModal';
import { useNotification } from '~/app/hooks/useNotification';
import CatalogSourceStatus from '~/app/pages/modelCatalogSettings/components/CatalogSourceStatus';

type CatalogSourceConfigsTableRowProps = {
  catalogSourceConfig: CatalogSourceConfig;
  onDeleteSource: (sourceId: string) => Promise<void>;
  isUpdatingToggle: boolean;
  onToggleUpdate: (checked: boolean, sourceConfig: CatalogSourceConfig) => void;
};

const CatalogSourceConfigsTableRow: React.FC<CatalogSourceConfigsTableRowProps> = ({
  catalogSourceConfig,
  onDeleteSource,
  isUpdatingToggle,
  onToggleUpdate,
}) => {
  const navigate = useNavigate();
  const notification = useNotification();
  const [isDeleteModalOpen, setIsDeleteModalOpen] = React.useState(false);
  const [isDeleting, setIsDeleting] = React.useState(false);
  const [deleteError, setDeleteError] = React.useState<Error | undefined>();

  const isDefault = catalogSourceConfig.isDefault ?? false;
  const isEnabled = catalogSourceConfig.enabled ?? true;

  const hasFilters = React.useMemo(
    () => hasSourceFilters(catalogSourceConfig),
    [catalogSourceConfig],
  );

  const handleEnableToggle = (checked: boolean) => {
    onToggleUpdate(checked, catalogSourceConfig);
  };

  const handleManageSource = () => {
    navigate(manageSourceUrl(catalogSourceConfig.id));
  };

  const handleDeleteClick = () => {
    setDeleteError(undefined);
    setIsDeleteModalOpen(true);
  };

  const handleDeleteConfirm = async () => {
    setIsDeleting(true);
    setDeleteError(undefined);

    try {
      await onDeleteSource(catalogSourceConfig.id);
      setIsDeleteModalOpen(false);
      notification.success(`${catalogSourceConfig.name} deleted successfully`);
    } catch (error) {
      setDeleteError(error instanceof Error ? error : new Error('Failed to delete source'));
    } finally {
      setIsDeleting(false);
    }
  };

  const handleCloseDeleteModal = () => {
    if (!isDeleting) {
      setIsDeleteModalOpen(false);
    }
  };

  const organizationValue = getOrganizationDisplay(catalogSourceConfig, isDefault);

  return (
    <>
      <Tr>
        <Td dataLabel="Name">
          <span data-testid={`source-name-${catalogSourceConfig.id}`}>
            {catalogSourceConfig.name}
          </span>
        </Td>
        <Td dataLabel="Organization">
          <span data-testid={`source-organization-${catalogSourceConfig.id}`}>
            {organizationValue}
          </span>
        </Td>
        <Td dataLabel="Model visibility">
          {hasFilters ? (
            <Label
              color={ModelVisibilityBadgeColor.FILTERED}
              data-testid={`model-visibility-filtered-${catalogSourceConfig.id}`}
            >
              Filtered
            </Label>
          ) : (
            <Label
              color={ModelVisibilityBadgeColor.UNFILTERED}
              data-testid={`model-visibility-unfiltered-${catalogSourceConfig.id}`}
            >
              Unfiltered
            </Label>
          )}
        </Td>
        <Td dataLabel="Source type">
          <span data-testid={`source-type-${catalogSourceConfig.id}`}>
            {CATALOG_SOURCE_TYPE_LABELS[catalogSourceConfig.type]}
          </span>
        </Td>
        <Td dataLabel="Enable">
          <Switch
            data-testid={`enable-toggle-${catalogSourceConfig.id}`}
            id={`enable-toggle-${catalogSourceConfig.id}`}
            aria-label={`Enable ${catalogSourceConfig.name}`}
            isChecked={isEnabled}
            isDisabled={isUpdatingToggle}
            onChange={(_event, checked) => handleEnableToggle(checked)}
          />
        </Td>
        <Td dataLabel="Validation status">
          <CatalogSourceStatus catalogSourceConfig={catalogSourceConfig} />
        </Td>
        <Td dataLabel="Actions">
          <Button
            variant="link"
            onClick={handleManageSource}
            data-testid={`manage-source-button-${catalogSourceConfig.id}`}
          >
            Manage source
          </Button>
        </Td>
        <Td isActionCell>
          {!isDefault && (
            <ActionsColumn
              items={[
                {
                  title: 'Delete source',
                  onClick: handleDeleteClick,
                },
              ]}
              data-testid={`source-actions-${catalogSourceConfig.id}`}
            />
          )}
        </Td>
      </Tr>
      {isDeleteModalOpen && (
        <DeleteModal
          title="Delete a source"
          testId="delete-source-modal"
          onClose={handleCloseDeleteModal}
          deleting={isDeleting}
          onDelete={handleDeleteConfirm}
          deleteName={catalogSourceConfig.name}
          error={deleteError}
        >
          The <strong>{catalogSourceConfig.name}</strong> repository will be deleted, and its models
          will be removed from the model catalog.
        </DeleteModal>
      )}
    </>
  );
};

export default CatalogSourceConfigsTableRow;
