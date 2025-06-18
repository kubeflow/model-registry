import * as React from 'react';
import * as _ from 'lodash-es';

export const useDeepCompareMemoize = <T>(value: T): T => {
  const ref = React.useRef<T>(value);

  if (!_.isEqual(value, ref.current)) {
    ref.current = value;
  }

  return ref.current;
};
