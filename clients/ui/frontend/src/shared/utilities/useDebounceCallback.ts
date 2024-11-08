import * as _ from 'lodash-es';
import * as React from 'react';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const useDebounceCallback = <T extends (...args: any) => any>(
  fn: T,
  delay = 250,
): ReturnType<typeof _.debounce<T>> => React.useMemo(() => _.debounce<T>(fn, delay), [fn, delay]);

export default useDebounceCallback;
