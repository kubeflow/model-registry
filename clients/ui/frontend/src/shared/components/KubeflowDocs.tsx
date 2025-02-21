import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon } from '@patternfly/react-icons';

type Props = {
  buttonLabel?: string;
  dockLinks?: string;
  isInline?: boolean;
  linkTestId?: string;
};

const KubeflowDocs: React.FC<Props> = ({
  buttonLabel = 'Kubeflow Docs',
  dockLinks = 'https://www.kubeflow.org/docs/components/model-registry/installation/#installing-model-registry',
  isInline,
  linkTestId,
}) => (
  <Button
    isInline={isInline}
    variant="link"
    icon={<OutlinedQuestionCircleIcon />}
    data-testid={linkTestId}
    onClick={() => {
      window.open(dockLinks);
    }}
  >
    {buttonLabel}
  </Button>
);

export default KubeflowDocs;
