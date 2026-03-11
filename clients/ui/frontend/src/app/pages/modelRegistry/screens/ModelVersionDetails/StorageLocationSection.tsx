import React from 'react';
import {
  DescriptionList,
  ExpandableSection,
  Content,
  ContentVariants,
  Spinner,
  Alert,
  Popover,
  Button,
  Icon,
  List,
  ListItem,
  Title,
  Bullseye,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon, FolderIcon } from '@patternfly/react-icons';
import { DashboardDescriptionListGroup, InlineTruncatedClipboardCopy } from 'mod-arch-shared';
import { ModelTransferJob } from '~/app/types';
import { useNamespaces } from '~/app/hooks/useNamespaces';
import { FindAdministratorOptions } from '~/app/utilities/const';
import {
  StorageType,
  getStorageTypeLabel,
  getModelUriLabel,
  getModelUriPopoverContent,
  getSourceLabel,
  getDestinationUri,
  getSourcePath,
} from '~/app/utils';

type StorageLocationSectionProps = {
  fallbackNamespace: string;
  transferJob: ModelTransferJob | null;
  transferJobLoaded: boolean;
  transferJobError: Error | undefined;
  onRetry: () => void;
};

const NoAccessPopover: React.FC<{ namespace: string }> = ({ namespace }) => (
  <Popover
    headerContent={
      <>
        You don&apos;t have access to the <strong>{namespace}</strong> namespace.
      </>
    }
    bodyContent={
      <>
        <Content component={ContentVariants.p}>
          To request access to a new or existing namespace, contact your administrator.
        </Content>
        <Content component={ContentVariants.p}>Your administrator might be:</Content>
        <List>
          {FindAdministratorOptions.map((option) => (
            <ListItem key={option}>{option}</ListItem>
          ))}
        </List>
      </>
    }
  >
    <Button
      variant="plain"
      aria-label="More info about namespace access"
      data-testid="no-access-popover-button"
    >
      <OutlinedQuestionCircleIcon />
    </Button>
  </Popover>
);

const SecretDisplay: React.FC<{
  secretName?: string;
  storageType?: StorageType;
  namespace: string;
  hasAccess?: boolean;
}> = ({ secretName, storageType, namespace, hasAccess = true }) => (
  <Content component={ContentVariants.p} data-testid="storage-secret-display">
    {(secretName || storageType) && (
      <>
        <strong>
          {secretName || (storageType && `${getStorageTypeLabel(storageType)} storage`)}
        </strong>{' '}
        secret in{' '}
      </>
    )}
    <Icon size="sm" isInline>
      <FolderIcon />
    </Icon>{' '}
    {namespace}
    {!hasAccess && <NoAccessPopover namespace={namespace} />}
  </Content>
);

const StorageLocationTitle: React.FC = () => (
  <Title headingLevel={ContentVariants.h3} data-testid="storage-location-title">
    Storage location{' '}
    <Popover
      headerContent="Storage location"
      bodyContent={
        <Content component={ContentVariants.p}>
          The secret currently being used to store this model and its artifacts. This is the{' '}
          <strong>Model destination location</strong> specified during registration.
        </Content>
      }
    >
      <Button
        variant="plain"
        aria-label="More info about storage location"
        data-testid="storage-location-popover-button"
      >
        <OutlinedQuestionCircleIcon />
      </Button>
    </Popover>
  </Title>
);

const StorageLocationSection: React.FC<StorageLocationSectionProps> = ({
  fallbackNamespace,
  transferJob,
  transferJobLoaded,
  transferJobError,
  onRetry,
}) => {
  const [isSourceDetailsExpanded, setIsSourceDetailsExpanded] = React.useState(false);
  const [namespaces, namespacesLoaded] = useNamespaces();

  const userHasNamespaceAccess =
    namespacesLoaded && namespaces.some((ns) => ns.name === fallbackNamespace);
  const hasAccessError =
    namespacesLoaded && !userHasNamespaceAccess && transferJobLoaded && !transferJob;

  if (transferJobError && !hasAccessError) {
    return (
      <Alert
        variant="danger"
        isInline
        title="Failed to load storage location"
        actionLinks={
          <Button variant="link" onClick={onRetry}>
            Retry
          </Button>
        }
      >
        {transferJobError.message}
      </Alert>
    );
  }

  if (!transferJobLoaded || !namespacesLoaded) {
    return (
      <Bullseye>
        <Spinner size="lg" />
      </Bullseye>
    );
  }

  if (hasAccessError) {
    return (
      <Stack hasGutter data-testid="storage-location-section">
        <StackItem>
          <StorageLocationTitle />
        </StackItem>
        <StackItem>
          <SecretDisplay namespace={fallbackNamespace} hasAccess={false} />
        </StackItem>
      </Stack>
    );
  }

  if (!transferJob) {
    return (
      <Stack hasGutter data-testid="storage-location-section">
        <StackItem>
          <StorageLocationTitle />
        </StackItem>
        <StackItem>
          <Content component={ContentVariants.p}>Storage information unavailable.</Content>
        </StackItem>
      </Stack>
    );
  }

  const destType = transferJob.destination.type;
  const sourceType = transferJob.source.type;
  const namespace = transferJob.namespace || '';
  const destinationUri = getDestinationUri(transferJob);
  const sourcePath = getSourcePath(transferJob);
  const sourceLabel = getSourceLabel(sourceType);

  return (
    <Stack hasGutter data-testid="storage-location-section">
      <StackItem>
        <StorageLocationTitle />
      </StackItem>

      <StackItem>
        <SecretDisplay
          secretName={transferJob.destSecretName}
          storageType={destType}
          namespace={namespace}
          hasAccess
        />
      </StackItem>

      <StackItem>
        <DescriptionList>
          <DashboardDescriptionListGroup
            title={getModelUriLabel(destType)}
            popover={getModelUriPopoverContent(destType)}
            isEmpty={!destinationUri}
            contentWhenEmpty="No URI"
          >
            <InlineTruncatedClipboardCopy
              testId="storage-location-uri"
              textToCopy={destinationUri}
            />
          </DashboardDescriptionListGroup>
        </DescriptionList>
      </StackItem>

      <StackItem>
        <ExpandableSection
          toggleText="Storage source details"
          onToggle={(_e, expanded) => setIsSourceDetailsExpanded(expanded)}
          isExpanded={isSourceDetailsExpanded}
          data-testid="storage-source-details"
        >
          <Stack hasGutter>
            <StackItem>
              <Content component={ContentVariants.p}>
                Details of the secret used to store the model before it was registered.
              </Content>
            </StackItem>

            <StackItem>
              <DescriptionList>
                <DashboardDescriptionListGroup
                  title="Model origin location"
                  popover="The secret that was used to store the model at the time it was registered."
                >
                  <SecretDisplay
                    secretName={transferJob.sourceSecretName}
                    storageType={sourceType}
                    namespace={namespace}
                  />
                </DashboardDescriptionListGroup>
              </DescriptionList>
            </StackItem>

            <StackItem>
              <DescriptionList>
                <DashboardDescriptionListGroup
                  title={sourceLabel}
                  isEmpty={!sourcePath}
                  contentWhenEmpty={`No ${sourceLabel.toLowerCase()}`}
                >
                  <InlineTruncatedClipboardCopy
                    testId="storage-source-path"
                    textToCopy={sourcePath}
                  />
                </DashboardDescriptionListGroup>
              </DescriptionList>
            </StackItem>
          </Stack>
        </ExpandableSection>
      </StackItem>
    </Stack>
  );
};

export default StorageLocationSection;
