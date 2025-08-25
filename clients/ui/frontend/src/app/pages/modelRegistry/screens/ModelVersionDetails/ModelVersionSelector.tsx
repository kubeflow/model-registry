import * as React from 'react';
import {
  Menu,
  MenuContainer,
  MenuContent,
  MenuItem,
  MenuList,
  MenuSearch,
  MenuSearchInput,
  MenuToggle,
  SearchInput,
  Divider,
  Badge,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { ModelVersion } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
import ViewAllVersionsButton from '~/app/pages/modelRegistry/screens/components/ViewAllVersionsButton';

type ModelVersionSelectorProps = {
  rmId?: string;
  selection: ModelVersion;
  onSelect: (versionId: string) => void;
};

const ModelVersionSelector: React.FC<ModelVersionSelectorProps> = ({
  rmId,
  selection,
  onSelect,
}) => {
  const [isOpen, setOpen] = React.useState(false);
  const [input, setInput] = React.useState('');

  const toggleRef = React.useRef(null);
  const menuRef = React.useRef(null);
  const [modelVersions] = useModelVersionsByRegisteredModel(rmId);
  const liveModelVersions = filterLiveVersions(modelVersions.items);
  const latestVersion = liveModelVersions.reduce<ModelVersion | null>((latest, current) => {
    if (
      latest === null ||
      Number(current.createTimeSinceEpoch) > Number(latest.createTimeSinceEpoch)
    ) {
      return current;
    }
    return latest;
  }, null);

  const menuListItems = liveModelVersions
    .filter((item) => input === '' || item.name.toLowerCase().includes(input.toLowerCase()))
    .toSorted((a, b) => Number(b.createTimeSinceEpoch) - Number(a.createTimeSinceEpoch)) // Sort by creation time, newest first
    .map((mv, index) => (
      <MenuItem isSelected={mv.id === selection.id} itemId={mv.id} key={index}>
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>{mv.name}</FlexItem>
          <FlexItem>
            {latestVersion && mv.id === latestVersion.id && <Badge color="blue">Latest</Badge>}
          </FlexItem>
        </Flex>
      </MenuItem>
    ));

  if (input.length > 0 && liveModelVersions.length === 0) {
    menuListItems.push(
      <MenuItem isDisabled key="no result">
        No results found
      </MenuItem>,
    );
  }

  const menu = (
    <Menu
      onSelect={(_e, itemId) => {
        if (typeof itemId === 'string') {
          onSelect(itemId);
          setOpen(false);
        }
      }}
      data-id="model-version-selector-menu"
      ref={menuRef}
      isScrollable
      activeItemId={selection.id}
    >
      <MenuSearch>
        <MenuSearchInput>
          <SearchInput
            data-testid="search-input"
            value={input}
            aria-label="Filter menu items"
            placeholder="Find by version name"
            onChange={(_event, value) => setInput(value)}
          />
        </MenuSearchInput>
      </MenuSearch>
      <Divider />
      <MenuContent>
        <MenuList data-testid="model-version-selector-list">
          {menuListItems}
          <MenuItem>
            <ViewAllVersionsButton rmId={rmId} totalVersions={modelVersions.items.length} />
          </MenuItem>
        </MenuList>
      </MenuContent>
    </Menu>
  );

  return (
    <MenuContainer
      isOpen={isOpen}
      toggleRef={toggleRef}
      toggle={
        <MenuToggle
          id="model-version-selector"
          ref={toggleRef}
          onClick={() => setOpen(!isOpen)}
          isExpanded={isOpen}
          isFullWidth
          data-testid="model-version-toggle-button"
        >
          <Flex spaceItems={{ default: 'spaceItemsSm' }}>
            <FlexItem>{selection.name}</FlexItem>
            <FlexItem>
              {latestVersion && selection.id === latestVersion.id && (
                <Badge color="blue">Latest</Badge>
              )}
            </FlexItem>
          </Flex>
        </MenuToggle>
      }
      menu={menu}
      menuRef={menuRef}
      onOpenChange={(open) => setOpen(open)}
    />
  );
};

export default ModelVersionSelector;
