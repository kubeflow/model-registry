import { Timestamp, TimestampTooltipVariant } from '@patternfly/react-core';
import React from 'react';
import { relativeTime } from 'mod-arch-shared';

type ModelTimestampProps = {
  timeSinceEpoch?: string;
};

const ModelTimestamp: React.FC<ModelTimestampProps> = ({ timeSinceEpoch }) => {
  if (!timeSinceEpoch) {
    return '--';
  }

  const time = new Date(parseInt(timeSinceEpoch)).getTime();

  if (Number.isNaN(time)) {
    return '--';
  }

  return (
    <Timestamp
      date={new Date(parseInt(timeSinceEpoch))}
      tooltip={{
        variant: TimestampTooltipVariant.default,
      }}
    >
      {relativeTime(Date.now(), time)}
    </Timestamp>
  );
};

export default ModelTimestamp;
