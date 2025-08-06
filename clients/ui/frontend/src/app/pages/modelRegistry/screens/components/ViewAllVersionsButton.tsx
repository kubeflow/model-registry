import * as React from 'react';
import { Button } from '@patternfly/react-core';
<<<<<<< HEAD
<<<<<<< HEAD
import { ArrowRightIcon } from '@patternfly/react-icons';
import {
  modelVersionListUrl,
  archiveModelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
=======
=======
import { ArrowRightIcon } from '@patternfly/react-icons';
>>>>>>> fe57b68 (moved the veiw button to screens)
import {
  modelVersionListUrl,
  archiveModelVersionListUrl,
} from '~/app/pages/modelRegistry/screens/routeUtils';
<<<<<<< HEAD
import { ModelRegistry } from '~/app/types';
>>>>>>> 77db33d (added versions card to model details (#1392))
=======
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
>>>>>>> fe57b68 (moved the veiw button to screens)

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
<<<<<<< HEAD
<<<<<<< HEAD
  isArchiveModel?: boolean;
  showIcon?: boolean;
=======
  preferredModelRegistry?: ModelRegistry;
  isArchiveModel?: boolean;
  icon?: React.ReactNode;
>>>>>>> 77db33d (added versions card to model details (#1392))
=======
  isArchiveModel?: boolean;
  showIcon?: boolean;
>>>>>>> fe57b68 (moved the veiw button to screens)
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
<<<<<<< HEAD
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
=======
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
>>>>>>> fe57b68 (moved the veiw button to screens)

export default ViewAllVersionsButton;
