import * as React from 'react';
import { Button } from '@patternfly/react-core';
<<<<<<< HEAD
import { ArrowRightIcon } from '@patternfly/react-icons';
import {
  modelVersionListUrl,
  archiveModelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
=======
import {
  archiveModelVersionListUrl,
  modelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistry } from '~/app/types';
>>>>>>> 77db33d (added versions card to model details (#1392))

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
<<<<<<< HEAD
  isArchiveModel?: boolean;
  showIcon?: boolean;
=======
  preferredModelRegistry?: ModelRegistry;
  isArchiveModel?: boolean;
  icon?: React.ReactNode;
>>>>>>> 77db33d (added versions card to model details (#1392))
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
<<<<<<< HEAD
  isArchiveModel,
  showIcon = false,
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
      icon={showIcon ? <ArrowRightIcon /> : undefined}
      iconPosition={showIcon ? 'right' : undefined}
    >
      {`View all ${totalVersions} versions`}
    </Button>
  );
};
=======
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
>>>>>>> 77db33d (added versions card to model details (#1392))

export default ViewAllVersionsButton;
