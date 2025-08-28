import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { Link } from 'react-router-dom';
import {
  modelVersionListUrl,
  archiveModelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
  isArchiveModel?: boolean;
  showIcon?: boolean;
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
  isArchiveModel,
  showIcon = false,
}) => {
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  return (
    <Link
      to={
        isArchiveModel
          ? archiveModelVersionListUrl(rmId, preferredModelRegistry?.name)
          : modelVersionListUrl(rmId, preferredModelRegistry?.name)
      }
      style={{ textDecoration: 'none' }}
    >
      <Button
        isInline
        data-testid="versions-route-link"
        variant="link"
        icon={showIcon ? <ArrowRightIcon /> : undefined}
        iconPosition={showIcon ? 'right' : undefined}
      >
        {`View all ${totalVersions} versions`}
      </Button>
    </Link>
  );
};

export default ViewAllVersionsButton;
