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
<<<<<<< HEAD
  Badge,
  Flex,
  FlexItem,
=======
  Button,
  Label,
>>>>>>> 7ad51d9 (improving the version selector in version details)
} from '@patternfly/react-core';
import { ModelVersion } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
<<<<<<< HEAD
import ViewAllVersionsButton from '~/app/pages/modelRegistry/screens/components/ViewAllVersionsButton';
=======
import { modelVersionListUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { useNavigate } from 'react-router-dom';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
>>>>>>> 7ad51d9 (improving the version selector in version details)

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
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const [modelVersions] = useModelVersionsByRegisteredModel(rmId);
  const liveModelVersions = filterLiveVersions(modelVersions.items);
<<<<<<< HEAD
  const latestVersion = liveModelVersions.reduce<ModelVersion | null>((latest, current) => {
    if (
      latest === null ||
      Number(current.createTimeSinceEpoch) > Number(latest.createTimeSinceEpoch)
    ) {
      return current;
    }
    return latest;
  }, null);
=======
  const latestVersion = modelVersions.items.reduce((latest, current) =>
    Number(current.createTimeSinceEpoch) > Number(latest.createTimeSinceEpoch) ? current : latest,
    modelVersions.items[0]
  );
>>>>>>> 7ad51d9 (improving the version selector in version details)

  const menuListItems = liveModelVersions
    .filter((item) => input === '' || item.name.toLowerCase().includes(input.toLowerCase()))
    .map((mv, index) => (
      <MenuItem isSelected={mv.id === selection.id} itemId={mv.id} key={index}>
<<<<<<< HEAD
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>{mv.name}</FlexItem>
          <FlexItem>
            {latestVersion && mv.id === latestVersion.id && <Badge color="blue">Latest</Badge>}
          </FlexItem>
        </Flex>
=======
        {mv.name}
        {mv.id === latestVersion.id && (
          <Label color="blue" style={{ marginLeft: 8 }}>
            Latest
          </Label>
        )}
>>>>>>> 7ad51d9 (improving the version selector in version details)
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
<<<<<<< HEAD
        <MenuList data-testid="model-version-selector-list">
          {menuListItems}
          <MenuItem>
            <ViewAllVersionsButton rmId={rmId} totalVersions={modelVersions.items.length} />
          </MenuItem>
        </MenuList>
=======
        <MenuSearch>
          <MenuSearchInput>
            <SearchInput
              data-testid="search-input"
              value={input}
              aria-label="Filter menu items"
              onChange={(_event, value) => setInput(value)}
            />
          </MenuSearchInput>
        </MenuSearch>
        <Divider />
        <MenuList data-testid="model-version-selector-list">{menuListItems}</MenuList>
        <MenuItem>
          <Button
            variant="link"
            isInline
            style={{ textTransform: 'none' }}
            icon={<ArrowRightIcon />}
            iconPosition="right"
            onClick={() => {
              setOpen(false);
              navigate(modelVersionListUrl(rmId, preferredModelRegistry?.name));
            }}
            data-testid="view-all-versions-link"
          >
            {`View all ${modelVersions.items.length} versions`}
          </Button>
        </MenuItem>
>>>>>>> 7ad51d9 (improving the version selector in version details)
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
<<<<<<< HEAD
=======
      popperProps={{ minWidth: '250px', maxWidth: 'none' }}
>>>>>>> 7ad51d9 (improving the version selector in version details)
      onOpenChange={(open) => setOpen(open)}
    />
  );
};

export default ModelVersionSelector;
