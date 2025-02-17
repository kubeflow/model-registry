import * as React from 'react';
import { Button } from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon } from '@patternfly/react-icons';

type Props = {
  buttonLabel?: string;
  docsLink?: string;
  isInline?: boolean;
  linkTestId?: string;
};

const KubeflowDocs: React.FC<Props> = ({
  buttonLabel = 'Kubeflow Docs',
  docsLink = 'https://www.kubeflow.org/docs/components/model-registry/installation/#installing-model-registry',
  isInline,
  linkTestId,
}) => (
  <Button
    isInline={isInline}
    variant="link"
    icon={<OutlinedQuestionCircleIcon />}
    data-testid={linkTestId}
    onClick={() => {
      window.open(docsLink);
    }}
  >
    {buttonLabel}
  </Button>
);

export default KubeflowDocs;
