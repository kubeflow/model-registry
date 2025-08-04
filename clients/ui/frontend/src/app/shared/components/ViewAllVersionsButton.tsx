import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { ArrowRightIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { modelVersionListUrl } from '~/app/pages/modelRegistry/screens/routeUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';

type ViewAllVersionsButtonProps = {
  rmId?: string;
  totalVersions: number;
  onClose?: () => void;
  className?: string;
};

const ViewAllVersionsButton: React.FC<ViewAllVersionsButtonProps> = ({
  rmId,
  totalVersions,
  onClose,
  className,
}) => {
  const navigate = useNavigate();
  const { preferredModelRegistry } = React.useContext(ModelRegistrySelectorContext);

  const handleClick = () => {
    onClose?.();
    navigate(modelVersionListUrl(rmId, preferredModelRegistry?.name));
  };

  return (
    <Button
      variant="link"
      isInline
      style={{ textTransform: 'none' }}
      icon={<ArrowRightIcon />}
      iconPosition="right"
      onClick={handleClick}
      data-testid="view-all-versions-link"
      className={className}
    >
      {`View all ${totalVersions} versions`}
    </Button>
  );
};

export default ViewAllVersionsButton;
