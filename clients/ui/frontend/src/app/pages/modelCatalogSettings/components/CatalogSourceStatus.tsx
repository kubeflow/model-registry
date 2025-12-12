import * as React from 'react';
import { Button, Label, Spinner, Stack, StackItem, Truncate } from '@patternfly/react-core';
import { CheckCircleIcon, ExclamationCircleIcon, InProgressIcon } from '@patternfly/react-icons';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import { ModelCatalogSettingsContext } from '~/app/context/modelCatalogSettings/ModelCatalogSettingsContext';
import { CatalogSourceStatus as CatalogSourceStatusEnum } from '~/concepts/modelCatalogSettings/const';
import CatalogSourceStatusErrorModal from './CatalogSourceStatusErrorModal';

type CatalogSourceStatusProps = {
  catalogSourceConfig: CatalogSourceConfig;
};

const CatalogSourceStatus: React.FC<CatalogSourceStatusProps> = ({ catalogSourceConfig }) => {
  const { catalogSources, catalogSourcesLoaded, catalogSourcesLoadError } = React.useContext(
    ModelCatalogSettingsContext,
  );
  const [isErrorModalOpen, setIsErrorModalOpen] = React.useState(false);

  // Don't render status for default sources
  if (catalogSourceConfig.isDefault) {
    return <>-</>;
  }

  // If source is disabled, render "-"
  if (!catalogSourceConfig.enabled) {
    return <>-</>;
  }

  // Show loading spinner while fetching sources
  if (!catalogSourcesLoaded) {
    return <Spinner size="md" data-testid={`source-status-loading-${catalogSourceConfig.id}`} />;
  }

  // Find the matching source from the catalog sources list
  const matchingSource = catalogSources?.items?.find(
    (source) => source.id === catalogSourceConfig.id,
  );

  const startingOrUnknownLabel = (
    <Label
      color="grey"
      icon={<InProgressIcon />}
      data-testid={`source-status-${catalogSourcesLoadError ? 'unknown' : 'starting'}-${catalogSourceConfig.id}`}
    >
      {catalogSourcesLoadError ? 'Unknown' : 'Starting'}
    </Label>
  );

  if (!matchingSource || !matchingSource.status) {
    return startingOrUnknownLabel;
  }

  // Render based on status
  switch (matchingSource.status) {
    case CatalogSourceStatusEnum.AVAILABLE:
      return (
        <Label
          color="green"
          icon={<CheckCircleIcon />}
          data-testid={`source-status-connected-${catalogSourceConfig.id}`}
        >
          Connected
        </Label>
      );

    case CatalogSourceStatusEnum.ERROR: {
      const errorMessage = matchingSource.error || 'Unknown error occurred';

      return (
        <>
          <Stack hasGutter>
            <StackItem>
              <Label
                color="red"
                icon={<ExclamationCircleIcon />}
                data-testid={`source-status-failed-${catalogSourceConfig.id}`}
              >
                Failed
              </Label>
            </StackItem>
            <StackItem>
              <Button
                variant="link"
                isInline
                isDanger
                onClick={() => setIsErrorModalOpen(true)}
                data-testid={`source-status-error-link-${catalogSourceConfig.id}`}
              >
                <Truncate content={errorMessage} tooltipProps={{ hidden: true }} />
              </Button>
            </StackItem>
          </Stack>
          <CatalogSourceStatusErrorModal
            isOpen={isErrorModalOpen}
            onClose={() => setIsErrorModalOpen(false)}
            errorMessage={errorMessage}
          />
        </>
      );
    }

    case CatalogSourceStatusEnum.DISABLED:
      return <>-</>;

    default:
      return startingOrUnknownLabel;
  }
};

export default CatalogSourceStatus;
