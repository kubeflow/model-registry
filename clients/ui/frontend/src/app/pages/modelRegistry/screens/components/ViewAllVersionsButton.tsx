import * as React from 'react';
import { Button } from '@patternfly/react-core';
import {
  archiveModelVersionListUrl,
  modelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistry } from '~/app/types';

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
  preferredModelRegistry?: ModelRegistry;
  isArchiveModel?: boolean;
  icon?: React.ReactNode;
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
  preferredModelRegistry,
  isArchiveModel,
  icon,
}) => (
  <Button
    component="a"
    isInline
    data-testid="versions-route-link"
    href={
      isArchiveModel
        ? archiveModelVersionListUrl(rmId, preferredModelRegistry?.name)
        : modelVersionListUrl(rmId, preferredModelRegistry?.name)
    }
    variant="link"
    icon={icon}
    iconPosition={icon ? 'right' : undefined}
  >
    {`View all ${totalVersions} versions`}
  </Button>
);

export default ViewAllVersionsButton;
