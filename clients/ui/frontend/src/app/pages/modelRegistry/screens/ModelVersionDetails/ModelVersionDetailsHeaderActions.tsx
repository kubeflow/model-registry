import * as React from 'react';
import { Dropdown, DropdownList, MenuToggle, DropdownItem } from '@patternfly/react-core';
import { useNavigate } from 'react-router';
import { ModelState, ModelVersion } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ArchiveModelVersionModal } from '~/app/pages/modelRegistry/screens/components/ArchiveModelVersionModal';
import { modelVersionArchiveDetailsUrl } from '~/app/pages/modelRegistry/screens/routeUtils';

interface ModelVersionsDetailsHeaderActionsProps {
  mv: ModelVersion;
  refresh: () => void;
}

const ModelVersionsDetailsHeaderActions: React.FC<ModelVersionsDetailsHeaderActionsProps> = ({
  mv,
}) => {
  const { apiState } = React.useContext(ModelRegistryContext);
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const navigate = useNavigate();
  const [isOpenActionDropdown, setOpenActionDropdown] = React.useState(false);
  const [isArchiveModalOpen, setIsArchiveModalOpen] = React.useState(false);
  const tooltipRef = React.useRef<HTMLButtonElement>(null);

  return (
    <>
      <Dropdown
        isOpen={isOpenActionDropdown}
        onSelect={() => setOpenActionDropdown(false)}
        onOpenChange={(open) => setOpenActionDropdown(open)}
        popperProps={{ position: 'right' }}
        toggle={(toggleRef) => (
          <MenuToggle
            variant="primary"
            ref={toggleRef}
            onClick={() => setOpenActionDropdown(!isOpenActionDropdown)}
            isExpanded={isOpenActionDropdown}
            aria-label="Model version details action toggle"
            data-testid="model-version-details-action-button"
          >
            Actions
          </MenuToggle>
        )}
      >
        <DropdownList>
          <DropdownItem
            id="archive-version-button"
            aria-label="Archive version"
            key="archive-version-button"
            onClick={() => setIsArchiveModalOpen(true)}
            ref={tooltipRef}
          >
            Archive version
          </DropdownItem>
        </DropdownList>
      </Dropdown>
      <ArchiveModelVersionModal
        onCancel={() => setIsArchiveModalOpen(false)}
        onSubmit={() =>
          apiState.api
            .patchModelVersion(
              {},
              {
                state: ModelState.ARCHIVED,
              },
              mv.id,
            )
            .then(() =>
              navigate(
                modelVersionArchiveDetailsUrl(
                  mv.id,
                  mv.registeredModelId,
                  preferredModelRegistry?.name,
                ),
              ),
            )
        }
        isOpen={isArchiveModalOpen}
        modelVersionName={mv.name}
      />
    </>
  );
};

export default ModelVersionsDetailsHeaderActions;
