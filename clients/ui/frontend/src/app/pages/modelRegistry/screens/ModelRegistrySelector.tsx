import * as React from 'react';
import {
  Bullseye,
  Button,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Flex,
  FlexItem,
  Icon,
  Popover,
  PopoverPosition,
  Tooltip,
} from '@patternfly/react-core';
import text from '@patternfly/react-styles/css/utilities/Text/text';
import truncateStyles from '@patternfly/react-styles/css/components/Truncate/truncate';
import { InfoCircleIcon, BlueprintIcon } from '@patternfly/react-icons';
import {
  WhosMyAdministrator,
  KubeflowDocs,
  SimpleSelect,
  InlineTruncatedClipboardCopy,
} from 'mod-arch-shared';
import { useBrowserStorage } from 'mod-arch-core';
import { useThemeContext } from 'mod-arch-kubeflow';
import { SimpleSelectOption } from 'mod-arch-shared/dist/components/SimpleSelect';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistry } from '~/app/types';
import { getServerAddress } from './utils';

const MODEL_REGISTRY_FAVORITE_STORAGE_KEY = 'kubeflow.dashboard.model.registry.favorite';

type ModelRegistrySelectorProps = {
  modelRegistry: string;
  onSelection: (modelRegistry: string) => void;
  primary?: boolean;
  isFullWidth?: boolean;
  hasError?: boolean;
};

const ModelRegistrySelector: React.FC<ModelRegistrySelectorProps> = ({
  modelRegistry,
  onSelection,
  primary,
  isFullWidth,
  hasError,
}) => {
  const { modelRegistries, updatePreferredModelRegistry } = React.useContext(
    ModelRegistrySelectorContext,
  );
  const { isMUITheme } = useThemeContext();

  const selection = modelRegistries.find((mr) => mr.name === modelRegistry);
  const [favorites, setFavorites] = useBrowserStorage<string[]>(
    MODEL_REGISTRY_FAVORITE_STORAGE_KEY,
    [],
  );

  const selectionDisplayName = selection ? selection.displayName : modelRegistry;

  const toggleLabel = modelRegistries.length === 0 ? 'No model registries' : selectionDisplayName;

  const getMRSelectDescription = (mr: ModelRegistry) => {
    const desc = mr.description || '';
    if (!desc) {
      return;
    }
    const tooltipContent = (
      <DescriptionList>
        <DescriptionListGroup>
          <DescriptionListTerm>{`${mr.displayName} description`}</DescriptionListTerm>
          <DescriptionListDescription>{desc}</DescriptionListDescription>
        </DescriptionListGroup>
      </DescriptionList>
    );
    return (
      <Tooltip content={tooltipContent} isContentLeftAligned>
        <span className={truncateStyles.truncate}>
          <span className={truncateStyles.truncateStart}>{desc}</span>
        </span>
      </Tooltip>
    );
  };

  const allOptions: SimpleSelectOption[] = modelRegistries.map((mr) => ({
    key: mr.name,
    label: mr.name,
    dropdownLabel: mr.displayName,
    description: getMRSelectDescription(mr),
    isFavorited: favorites.includes(mr.name),
  }));

  const favoriteOptions = (favIds: string[]) =>
    allOptions.filter((option) => favIds.includes(option.key));

  const selector = (
    <SimpleSelect
      isScrollable
      placeholder="Select a model registry"
      dataTestId="model-registry-selector-dropdown"
      toggleProps={{ id: 'download-steps-logs-toggle', status: hasError ? 'danger' : undefined }}
      toggleLabel={toggleLabel}
      aria-label="Model registry toggle"
      previewDescription={false}
      onChange={(key) => {
        updatePreferredModelRegistry(modelRegistries.find((obj) => obj.name === key));
        onSelection(key);
      }}
      isFullWidth={isFullWidth}
      maxMenuHeight="300px"
      popperProps={{ maxWidth: '400px' }}
      value={selection?.name}
      groupedOptions={[
        ...(favorites.length > 0
          ? [
              {
                key: 'favorites-group',
                label: 'Favorites',
                options: favoriteOptions(favorites),
              },
            ]
          : []),
        {
          key: 'all',
          label: 'All model registries',
          options: allOptions,
        },
      ]}
      onActionClick={(event: React.MouseEvent, value: string, actionId: string) => {
        event.stopPropagation();
        if (actionId === 'fav') {
          const isFavorited = favorites.includes(value);
          if (isFavorited) {
            setFavorites(favorites.filter((id) => id !== value));
          } else {
            setFavorites([...favorites, value]);
          }
        }
      }}
    />
  );

  if (primary) {
    return selector;
  }

  return (
    <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>
        <Icon>
          <BlueprintIcon />
        </Icon>
      </FlexItem>
      <FlexItem>
        <Bullseye>Model registry</Bullseye>
      </FlexItem>
      <FlexItem>{selector}</FlexItem>
      {selection && (
        <FlexItem>
          <Popover
            aria-label="Model registry description popover"
            data-testid="mr-details-popover"
            position="right"
            headerContent={`${selection.displayName} details`}
            bodyContent={
              <DescriptionList>
                <DescriptionListGroup>
                  <DescriptionListTerm>Description</DescriptionListTerm>
                  <DescriptionListDescription
                    className={!selection.description ? text.textColorDisabled : ''}
                  >
                    {selection.description || 'No description'}
                  </DescriptionListDescription>
                </DescriptionListGroup>
                <DescriptionListGroup>
                  <DescriptionListTerm>Server URL</DescriptionListTerm>
                  <DescriptionListDescription>
                    <InlineTruncatedClipboardCopy
                      textToCopy={`https://${getServerAddress(selection)}`}
                    />
                  </DescriptionListDescription>
                </DescriptionListGroup>
              </DescriptionList>
            }
          >
            <Button variant="link" icon={<InfoCircleIcon />} data-testid="view-details-button">
              View details
            </Button>
          </Popover>
        </FlexItem>
      )}
      <FlexItem align={{ default: 'alignRight' }}>
        {isMUITheme ? (
          <KubeflowDocs
            buttonLabel="Need another registry?"
            linkTestId="model-registry-help-button"
          />
        ) : (
          <WhosMyAdministrator
            buttonLabel="Need another registry?"
            headerContent="Need another registry?"
            leadText="To request access to a new or existing model registry, contact your administrator."
            contentTestId="model-registry-help-content"
            linkTestId="model-registry-help-button"
            popoverPosition={PopoverPosition.left}
          />
        )}
      </FlexItem>
    </Flex>
  );
};

export default ModelRegistrySelector;
