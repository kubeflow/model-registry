import * as React from 'react';
import { Td } from '@patternfly/react-table';
import { Checkbox, Tooltip } from '@patternfly/react-core';

type CheckboxTrProps = {
  id: string;
  isChecked: boolean | null;
  onToggle: () => void;
  isDisabled?: boolean;
  tooltip?: string;
} & React.ComponentProps<typeof Td>;

const CheckboxTd: React.FC<CheckboxTrProps> = ({
  id,
  isChecked,
  onToggle,
  isDisabled,
  tooltip,
  ...props
}) => {
  let content = (
    <Checkbox
      aria-label="Checkbox"
      id={`${id}-checkbox`}
      isChecked={isChecked}
      onChange={() => onToggle()}
      isDisabled={isDisabled}
    />
  );

  if (tooltip) {
    content = <Tooltip content={tooltip}>{content}</Tooltip>;
  }

  return (
    <Td dataLabel="Checkbox" {...props}>
      {content}
    </Td>
  );
};

export default CheckboxTd;
