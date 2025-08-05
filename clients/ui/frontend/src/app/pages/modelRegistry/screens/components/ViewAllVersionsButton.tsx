import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
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
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
  preferredModelRegistry,
  isArchiveModel,
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
    icon={<ArrowRightIcon />}
    iconPosition="right"
  >
    {`View all ${totalVersions} versions`}
  </Button>
);

export default ViewAllVersionsButton;
