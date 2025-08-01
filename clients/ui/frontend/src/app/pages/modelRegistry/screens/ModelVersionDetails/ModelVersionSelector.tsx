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
  Button,
  Label,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router-dom';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { ModelVersion } from '~/app/types';
import useModelVersionsByRegisteredModel from '~/app/hooks/useModelVersionsByRegisteredModel';
import { filterLiveVersions } from '~/app/utils';
import { modelVersionListUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';

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
  const latestVersion = modelVersions.items.reduce(
    (latest, current) =>
      Number(current.createTimeSinceEpoch) > Number(latest.createTimeSinceEpoch) ? current : latest,
    modelVersions.items[0],
  );

  const menuListItems = liveModelVersions
    .filter((item) => !input || item.name.toLowerCase().includes(input.toString().toLowerCase()))
    .map((mv, index) => (
      <MenuItem isSelected={mv.id === selection.id} itemId={mv.id} key={index}>
        <Flex spaceItems={{ default: 'spaceItemsSm' }}>
          <FlexItem>{mv.name}</FlexItem>
          <FlexItem>{mv.id === latestVersion.id && <Label color="blue">Latest</Label>}</FlexItem>
        </Flex>
      </MenuItem>
    ));

  if (input && liveModelVersions.length === 0) {
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
      <MenuContent>
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
          {selection.name}
        </MenuToggle>
      }
      menu={menu}
      menuRef={menuRef}
      onOpenChange={(open) => setOpen(open)}
    />
  );
};

export default ModelVersionSelector;
