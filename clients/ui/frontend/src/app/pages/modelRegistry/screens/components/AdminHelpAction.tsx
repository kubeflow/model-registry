import * as React from 'react';
import { WhosMyAdministrator, KubeflowDocs } from 'mod-arch-shared';
import { useThemeContext } from 'mod-arch-kubeflow';
import { PopoverPosition } from '@patternfly/react-core';

type AdminHelpActionProps = {
  buttonLabel?: string;
  linkTestId?: string;
  headerContent?: string;
  leadText?: string;
  contentTestId?: string;
  popoverPosition?: PopoverPosition;
};

const AdminHelpAction: React.FC<AdminHelpActionProps> = ({
  buttonLabel = "Who's my administrator?",
  linkTestId = 'whos-my-admin-link',
  headerContent = "Who's my administrator?",
  leadText = 'To request access to a new or existing model registry, contact your administrator.',
  contentTestId = 'whos-my-admin-content',
  popoverPosition = PopoverPosition.left,
}) => {
  const { isMUITheme } = useThemeContext();

  if (isMUITheme) {
    return <KubeflowDocs buttonLabel={buttonLabel} linkTestId={linkTestId} />;
  }
  return (
    <WhosMyAdministrator
      buttonLabel={buttonLabel}
      headerContent={headerContent}
      leadText={leadText}
      contentTestId={contentTestId}
      linkTestId={linkTestId}
      popoverPosition={popoverPosition}
    />
  );
};

export default AdminHelpAction;
