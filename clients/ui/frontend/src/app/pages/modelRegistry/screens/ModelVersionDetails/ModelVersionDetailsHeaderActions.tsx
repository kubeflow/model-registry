import * as React from 'react';
import {
  Dropdown,
  DropdownList,
  MenuToggle,
  DropdownItem,
  ButtonVariant,
  ActionList,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router';
import { ModelState, ModelVersion } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { ArchiveModelVersionModal } from '~/app/pages/modelRegistry/screens/components/ArchiveModelVersionModal';
import { modelVersionListUrl } from '~/app/pages/modelRegistry/screens/routeUtils';

interface ModelVersionsDetailsHeaderActionsProps {
  mv: ModelVersion;
  hasDeployment?: boolean;
}

const ModelVersionsDetailsHeaderActions: React.FC<ModelVersionsDetailsHeaderActionsProps> = ({
  mv,
  hasDeployment = false,
}) => {
  const { apiState } = React.useContext(ModelRegistryContext);
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const navigate = useNavigate();
  const [isOpenActionDropdown, setOpenActionDropdown] = React.useState(false);
  const [isArchiveModalOpen, setIsArchiveModalOpen] = React.useState(false);
  const tooltipRef = React.useRef<HTMLButtonElement>(null);

  return (
    <ActionList>
      <Dropdown
        isOpen={isOpenActionDropdown}
        onSelect={() => setOpenActionDropdown(false)}
        onOpenChange={(open) => setOpenActionDropdown(open)}
        popperProps={{ position: 'right', appendTo: 'inline' }}
        toggle={(toggleRef) => (
          <MenuToggle
            variant={ButtonVariant.secondary}
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
            isAriaDisabled={hasDeployment}
            id="archive-version-button"
            aria-label="Archive model version"
            key="archive-version-button"
            onClick={() => setIsArchiveModalOpen(true)}
            tooltipProps={
              hasDeployment ? { content: 'Deployed model versions cannot be archived' } : undefined
            }
            ref={tooltipRef}
          >
            Archive model version
          </DropdownItem>
        </DropdownList>
      </Dropdown>
      {isArchiveModalOpen ? (
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
                navigate(modelVersionListUrl(mv.registeredModelId, preferredModelRegistry?.name)),
              )
          }
          modelVersionName={mv.name}
        />
      ) : null}
    </ActionList>
  );
};

export default ModelVersionsDetailsHeaderActions;
