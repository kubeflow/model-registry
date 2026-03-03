import * as React from 'react';
import { EmptyStateErrorMessage, WhosMyAdministrator, KubeflowDocs } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { PopoverPosition } from '@patternfly/react-core';

const UNAVAILABLE_TITLE = 'Model registry unavailable';

const getUnavailableBodyText = (registryDisplayName: string): string =>
  `The ${registryDisplayName} registry is currently unavailable. It might still be starting up, or there might be a configuration error. Wait a few minutes and try again. If the problem persists, contact your administrator.`;

type UnavailableModelRegistryProps = {
  /** Display name of the selected registry (e.g. "unavailable-registry-example") */
  registryDisplayName: string;
};

const UnavailableModelRegistry: React.FC<UnavailableModelRegistryProps> = ({
  registryDisplayName,
}) => {
  const { isMUITheme } = useThemeContext();

  return (
    <div data-testid="unavailable-model-registry">
      <EmptyStateErrorMessage
        title={UNAVAILABLE_TITLE}
        bodyText={getUnavailableBodyText(registryDisplayName)}
      >
        {isMUITheme ? (
          <KubeflowDocs buttonLabel="Who's my administrator?" linkTestId="whos-my-admin-link" />
        ) : (
          <WhosMyAdministrator
            buttonLabel="Who's my administrator?"
            headerContent="Who's my administrator?"
            leadText="To request access to a new or existing model registry, contact your administrator."
            contentTestId="whos-my-admin-content"
            linkTestId="whos-my-admin-link"
            popoverPosition={PopoverPosition.left}
          />
        )}
      </EmptyStateErrorMessage>
    </div>
  );
};

export default UnavailableModelRegistry;
