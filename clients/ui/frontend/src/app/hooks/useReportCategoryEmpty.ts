import * as React from 'react';

const useReportCategoryEmpty = (
  reportCategoryEmpty: (label: string, isEmpty: boolean) => void,
  label: string,
  isLoaded: boolean,
  itemCount: number,
  searchTerm: string,
): void => {
  React.useEffect(() => {
    if (!isLoaded || searchTerm) {
      return undefined;
    }
    const timer = setTimeout(() => {
      reportCategoryEmpty(label, itemCount === 0);
    }, 100);
    return () => clearTimeout(timer);
  }, [isLoaded, itemCount, label, searchTerm, reportCategoryEmpty]);
};

export default useReportCategoryEmpty;
