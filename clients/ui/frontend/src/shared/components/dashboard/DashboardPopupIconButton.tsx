import React from 'react';
import { Button, ButtonProps, Icon, IconComponentProps } from '@patternfly/react-core';

type DashboardPopupIconButtonProps = Omit<ButtonProps, 'variant' | 'isInline'> & {
  icon: React.ReactNode;
  iconProps?: Omit<IconComponentProps, 'isInline'>;
};

/**
 * Overriding PF's button styles to allow for a11y in opening tooltips or popovers on a single item
 */
const DashboardPopupIconButton = ({
  icon,
  iconProps,
  ...props
}: DashboardPopupIconButtonProps): React.JSX.Element => (
  <Button
    icon={
      <Icon isInline {...iconProps}>
        {icon}
      </Icon>
    }
    variant="plain"
    isInline
    {...props}
  />
);

export default DashboardPopupIconButton;
