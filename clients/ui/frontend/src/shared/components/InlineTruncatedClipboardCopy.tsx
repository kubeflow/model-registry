import { ClipboardCopy, ClipboardCopyVariant, Truncate } from '@patternfly/react-core';
import * as React from 'react';

type Props = {
  textToCopy: string;
  truncatePosition?: 'middle' | 'end';
  testId?: string;
  maxWidth?: number;
};

/** Hopefully PF will add some flexibility with ClipboardCopy
 *  in the future and this will not be necessary
 * https://github.com/patternfly/patternfly-react/issues/10890
 **/

// TODO: Fix this when PF 6 supports a ReactNode as a child for the ClipboardCopy component
const InlineTruncatedClipboardCopy: React.FC<Props> = ({
  textToCopy,
  testId,
  maxWidth,
  truncatePosition,
}) => (
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  <ClipboardCopy
    variant={ClipboardCopyVariant.inlineCompact}
    style={{ display: 'inline-flex', maxWidth }}
    hoverTip="Copy"
    clickTip="Copied"
    data-testid={testId}
    onCopy={() => {
      navigator.clipboard.writeText(textToCopy);
    }}
  >
    <Truncate content={textToCopy} position={truncatePosition} />
  </ClipboardCopy>
);

export default InlineTruncatedClipboardCopy;
