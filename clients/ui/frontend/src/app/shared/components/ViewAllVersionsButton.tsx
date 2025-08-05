import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
import {
  modelVersionListUrl,
  archiveModelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
  isArchiveModel?: boolean;
  className?: string;
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
  isArchiveModel,
  className,
}) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  return (
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
      style={{ textTransform: 'none' }}
      className={className}
    >
      {`View all ${totalVersions} versions`}
    </Button>
  );
};

export default ViewAllVersionsButton;
