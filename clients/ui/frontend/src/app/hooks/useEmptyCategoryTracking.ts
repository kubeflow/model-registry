import * as React from 'react';

type UseEmptyCategoryTrackingResult = {
  emptyCategoryLabels: Set<string>;
  reportCategoryEmpty: (label: string, isEmpty: boolean) => void;
};

const useEmptyCategoryTracking = (): UseEmptyCategoryTrackingResult => {
  const [emptyCategoryLabels, setEmptyCategoryLabels] = React.useState<Set<string>>(
    () => new Set<string>(),
  );
  const emptyCategoryLabelsRef = React.useRef(emptyCategoryLabels);
  emptyCategoryLabelsRef.current = emptyCategoryLabels;

  const reportCategoryEmpty = React.useCallback((label: string, isEmpty: boolean) => {
    const { current } = emptyCategoryLabelsRef;
    const hasLabel = current.has(label);
    if (isEmpty && !hasLabel) {
      const next = new Set(current);
      next.add(label);
      setEmptyCategoryLabels(next);
    } else if (!isEmpty && hasLabel) {
      const next = new Set(current);
      next.delete(label);
      setEmptyCategoryLabels(next);
    }
  }, []);

  return { emptyCategoryLabels, reportCategoryEmpty };
};

export default useEmptyCategoryTracking;
