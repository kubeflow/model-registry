import * as React from 'react';

type UseEmptyCategoryTrackingResult = {
  emptyCategoryLabels: Set<string>;
  reportCategoryEmpty: (label: string, isEmpty: boolean) => void;
};

const useEmptyCategoryTracking = (): UseEmptyCategoryTrackingResult => {
  const [emptyCategoryLabels, setEmptyCategoryLabels] = React.useState<Set<string>>(
    () => new Set<string>(),
  );

  const reportCategoryEmpty = React.useCallback((label: string, isEmpty: boolean) => {
    setEmptyCategoryLabels((prev) => {
      const hasLabel = prev.has(label);
      if (isEmpty && !hasLabel) {
        const next = new Set(prev);
        next.add(label);
        return next;
      }
      if (!isEmpty && hasLabel) {
        const next = new Set(prev);
        next.delete(label);
        return next;
      }
      return prev;
    });
  }, []);

  return { emptyCategoryLabels, reportCategoryEmpty };
};

export default useEmptyCategoryTracking;
