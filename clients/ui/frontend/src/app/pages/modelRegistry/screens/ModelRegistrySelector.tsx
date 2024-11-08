import * as React from 'react';
import {
  Bullseye,
  Button,
  DescriptionList,
  DescriptionListDescription,
  DescriptionListGroup,
  DescriptionListTerm,
  Divider,
  Flex,
  FlexItem,
  Icon,
  MenuToggle,
  Popover,
  Select,
  SelectGroup,
  SelectList,
  SelectOption,
  Tooltip,
} from '@patternfly/react-core';
import truncateStyles from '@patternfly/react-styles/css/components/Truncate/truncate';
import { InfoCircleIcon, BlueprintIcon } from '@patternfly/react-icons';
import { useBrowserStorage } from '~/shared/components/browserStorage';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ModelRegistry } from '~/app/types';

const MODEL_REGISTRY_FAVORITE_STORAGE_KEY = 'kubeflow.dashboard.model.registry.favorite';

type ModelRegistrySelectorProps = {
  modelRegistry: string;
  onSelection: (modelRegistry: string) => void;
  primary?: boolean;
};

const ModelRegistrySelector: React.FC<ModelRegistrySelectorProps> = ({
  modelRegistry,
  onSelection,
  primary,
}) => {
  const { modelRegistries, updatePreferredModelRegistry } = React.useContext(
    ModelRegistrySelectorContext,
  );

  const selection = modelRegistries.find((mr) => mr.name === modelRegistry);
  const [isOpen, setIsOpen] = React.useState(false);
  const [favorites, setFavorites] = useBrowserStorage<string[]>(
    MODEL_REGISTRY_FAVORITE_STORAGE_KEY,
    [],
  );

  const selectionDisplayName = selection ? selection.displayName : modelRegistry;

  const toggleLabel = modelRegistries.length === 0 ? 'No model registries' : selectionDisplayName;

  const getMRSelectDescription = (mr: ModelRegistry) => {
    const desc = mr.description || mr.name;
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

  const options = [
    <SelectGroup label="Select a model registry" key="all">
      <SelectList>
        {modelRegistries.map((mr) => (
          <SelectOption
            id={mr.name}
            key={mr.name}
            value={mr.name}
            description={getMRSelectDescription(mr)}
            isFavorited={favorites.includes(mr.name)}
          >
            {mr.displayName}
          </SelectOption>
        ))}
      </SelectList>
    </SelectGroup>,
  ];

  const createFavorites = (favIds: string[]) => {
    const favorite: JSX.Element[] = [];

    options.forEach((item) => {
      if (item.type === SelectList) {
        item.props.children.filter(
          (child: JSX.Element) => favIds.includes(child.props.value) && favorite.push(child),
        );
      } else if (item.type === SelectGroup) {
        item.props.children.props.children.filter(
          (child: JSX.Element) => favIds.includes(child.props.value) && favorite.push(child),
        );
      } else if (favIds.includes(item.props.value)) {
        favorite.push(item);
      }
    });

    return favorite;
  };

  const selector = (
    <Select
      toggle={(toggleRef) => (
        <MenuToggle
          ref={toggleRef}
          data-testid="model-registry-selector-dropdown"
          aria-label="Model registry toggle"
          id="download-steps-logs-toggle"
          onClick={() => setIsOpen(!isOpen)}
          isExpanded={isOpen}
          isDisabled={modelRegistries.length === 0}
        >
          {toggleLabel}
        </MenuToggle>
      )}
      onSelect={(_e, value) => {
        setIsOpen(false);
        updatePreferredModelRegistry(modelRegistries.find((obj) => obj.name === value));
        if (typeof value === 'string') {
          onSelection(value);
        }
      }}
      selected={toggleLabel}
      onOpenChange={(open) => setIsOpen(open)}
      isOpen={isOpen}
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
    >
      {favorites.length > 0 && (
        <React.Fragment key="favorites-group">
          <SelectGroup label="Favorites">
            <SelectList>{createFavorites(favorites)}</SelectList>
          </SelectGroup>
          <Divider />
        </React.Fragment>
      )}
      {options}
    </Select>
  );

  if (primary) {
    return selector;
  }

  return (
    <Flex spaceItems={{ default: 'spaceItemsXs' }} alignItems={{ default: 'alignItemsCenter' }}>
      <Icon>
        <BlueprintIcon />
      </Icon>
      <Flex spaceItems={{ default: 'spaceItemsSm' }} alignItems={{ default: 'alignItemsCenter' }}>
        <FlexItem>
          <Bullseye>Model registry</Bullseye>
        </FlexItem>
        <FlexItem>{selector}</FlexItem>
        {selection && selection.description && (
          <FlexItem>
            <Popover
              aria-label="Model registry description popover"
              headerContent={selection.displayName}
              bodyContent={selection.description}
            >
              <Button variant="link" icon={<InfoCircleIcon />}>
                View Description
              </Button>
            </Popover>
          </FlexItem>
        )}
      </Flex>
    </Flex>
  );
};

export default ModelRegistrySelector;
