import * as React from 'react';
import { Popover } from '@patternfly/react-core';
import { OutlinedQuestionCircleIcon } from '@patternfly/react-icons';
import DashboardPopupIconButton from '~/shared/components/dashboard/DashboardPopupIconButton';

type FieldGroupHelpLabelIconProps = {
  content: React.ComponentProps<typeof Popover>['bodyContent'];
};

const FieldGroupHelpLabelIcon: React.FC<FieldGroupHelpLabelIconProps> = ({ content }) => (
  <Popover bodyContent={content}>
    <DashboardPopupIconButton icon={<OutlinedQuestionCircleIcon />} aria-label="More info" />
  </Popover>
);

export default FieldGroupHelpLabelIcon;
