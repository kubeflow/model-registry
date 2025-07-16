import * as React from 'react';
import {
  Dropdown,
  DropdownItem,
  DropdownList,
  Flex,
  MenuToggle,
  MenuToggleAction,
  MenuToggleElement,
  Toolbar,
  ToolbarContent,
  ToolbarItem,
  ToolbarToggleGroup,
} from '@patternfly/react-core';
import { EllipsisVIcon, FilterIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import {
  registeredModelArchiveUrl,
  registerModelUrl,
  registerVersionUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';

type RegisteredModelsTableToolbarProps = {
  toggleGroupItems?: React.ReactNode;
  onClearAllFilters?: () => void;
};

const RegisteredModelsTableToolbar: React.FC<RegisteredModelsTableToolbarProps> = ({
  toggleGroupItems: tableToggleGroupItems,
  onClearAllFilters,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);
  const [isRegisterNewVersionOpen, setIsRegisterNewVersionOpen] = React.useState(false);
  const [isArchivedModelKebabOpen, setIsArchivedModelKebabOpen] = React.useState(false);

  const tooltipRef = React.useRef<HTMLButtonElement>(null);

  return (
    <Toolbar data-testid="registered-models-table-toolbar" clearAllFilters={onClearAllFilters}>
      <ToolbarContent>
        {/* TODO: Remove this Flex after the ToolbarContent can center the children elements */}
        <Flex>
          <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
            {tableToggleGroupItems}
          </ToolbarToggleGroup>
          <ToolbarItem>
            <Dropdown
              isOpen={isRegisterNewVersionOpen}
              onSelect={() => setIsRegisterNewVersionOpen(false)}
              onOpenChange={(isOpen) => setIsRegisterNewVersionOpen(isOpen)}
              toggle={(toggleRef) => (
                <MenuToggle
                  isFullWidth
                  variant="primary"
                  ref={toggleRef}
                  onClick={() => setIsRegisterNewVersionOpen(!isRegisterNewVersionOpen)}
                  isExpanded={isRegisterNewVersionOpen}
                  splitButtonItems={[
                    <MenuToggleAction
                      id="register-model-button"
                      key="register-model-button"
                      data-testid="register-model-button"
                      aria-label="Register model"
                      onClick={() => navigate(registerModelUrl(preferredModelRegistry?.name))}
                    >
                      Register model
                    </MenuToggleAction>,
                  ]}
                  aria-label="Register model toggle"
                  data-testid="register-model-split-button"
                />
              )}
            >
              <DropdownList>
                <DropdownItem
                  id="register-new-version-button"
                  aria-label="Register new version"
                  key="register-new-version-button"
                  onClick={() => {
                    navigate(registerVersionUrl(preferredModelRegistry?.name));
                  }}
                  ref={tooltipRef}
                >
                  Register new version
                </DropdownItem>
              </DropdownList>
            </Dropdown>
          </ToolbarItem>
          <ToolbarItem>
            <Dropdown
              isOpen={isArchivedModelKebabOpen}
              onSelect={() => setIsArchivedModelKebabOpen(false)}
              onOpenChange={(isOpen: boolean) => setIsArchivedModelKebabOpen(isOpen)}
              toggle={(tr: React.Ref<MenuToggleElement>) => (
                <MenuToggle
                  data-testid="registered-models-table-kebab-action"
                  ref={tr}
                  variant="plain"
                  onClick={() => setIsArchivedModelKebabOpen(!isArchivedModelKebabOpen)}
                  isExpanded={isArchivedModelKebabOpen}
                  aria-label="View archived models"
                >
                  <EllipsisVIcon />
                </MenuToggle>
              )}
              shouldFocusToggleOnSelect
            >
              <DropdownList>
                <DropdownItem
                  onClick={() => navigate(registeredModelArchiveUrl(preferredModelRegistry?.name))}
                >
                  View archived models
                </DropdownItem>
              </DropdownList>
            </Dropdown>
          </ToolbarItem>
        </Flex>
      </ToolbarContent>
    </Toolbar>
  );
};

export default RegisteredModelsTableToolbar;
