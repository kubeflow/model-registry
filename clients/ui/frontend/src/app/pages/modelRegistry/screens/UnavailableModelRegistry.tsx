import * as React from 'react';
import { EmptyStateErrorMessage } from 'mod-arch-shared';
import AdminHelpAction from './components/AdminHelpAction';

const UNAVAILABLE_TITLE = 'Model registry unavailable';

const getUnavailableBodyText = (registryDisplayName: string): string =>
  `The ${registryDisplayName} registry is currently unavailable. It might still be starting up, or there might be a configuration error. Wait a few minutes and try again. If the problem persists, contact your administrator.`;

type UnavailableModelRegistryProps = {
  /** Display name of the selected registry, shown in the message (e.g. "Unavailable Registry Example"). */
  registryDisplayName: string;
};

const UnavailableModelRegistry: React.FC<UnavailableModelRegistryProps> = ({
  registryDisplayName,
}) => (
  <div data-testid="unavailable-model-registry">
    <EmptyStateErrorMessage
      title={UNAVAILABLE_TITLE}
      bodyText={getUnavailableBodyText(registryDisplayName)}
    >
      <AdminHelpAction />
    </EmptyStateErrorMessage>
  </div>
);

export default UnavailableModelRegistry;
