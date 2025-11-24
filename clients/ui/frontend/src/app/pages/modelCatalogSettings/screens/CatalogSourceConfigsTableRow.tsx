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

type CatalogSourceConfigsTableRowProps = {
  catalogSourceConfig: CatalogSourceConfig;
  onDelete?: (config: CatalogSourceConfig) => void;
};

const CatalogSourceConfigsTableRow: React.FC<CatalogSourceConfigsTableRowProps> = ({
  catalogSourceConfig,
  onDelete,
}) => {
  const navigate = useNavigate();
  const isDefault = catalogSourceConfig.isDefault ?? false;
  const isEnabled = catalogSourceConfig.enabled ?? true;

  const hasFilters = React.useMemo(
    () => hasSourceFilters(catalogSourceConfig),
    [catalogSourceConfig],
  );

  const handleEnableToggle = (checked: boolean) => {
    // TODO: Implement actual enable/disable functionality
    window.alert(
      `Toggle clicked! "${catalogSourceConfig.name}" will be ${checked ? 'enabled' : 'disabled'} when functionality is implemented.`,
    );
  };

  const handleManageSource = () => {
    navigate(manageSourceUrl(catalogSourceConfig.id));
  };

  const handleDeleteSource = () => {
    // TODO: - Implement actual delete functionality
    onDelete?.(catalogSourceConfig);
  };

  const organizationValue = getOrganizationDisplay(catalogSourceConfig, isDefault);

  return (
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
        {!isDefault && (
          <Switch
            data-testid={`enable-toggle-${catalogSourceConfig.id}`}
            id={`enable-toggle-${catalogSourceConfig.id}`}
            aria-label={`Enable ${catalogSourceConfig.name}`}
            isChecked={isEnabled}
            onChange={(_event, checked) => handleEnableToggle(checked)}
          />
        )}
      </Td>
      <Td dataLabel="Validation status">{/* TODO: Status implementation */}</Td>
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
        <ActionsColumn
          items={[
            {
              title: 'Delete source',
              onClick: handleDeleteSource,
              isDisabled: isDefault,
            },
          ]}
          data-testid={`source-actions-${catalogSourceConfig.id}`}
        />
      </Td>
    </Tr>
  );
};

export default CatalogSourceConfigsTableRow;
