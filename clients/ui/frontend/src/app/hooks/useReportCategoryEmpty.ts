import * as React from 'react';

const useReportCategoryEmpty = (
  reportCategoryEmpty: (label: string, isEmpty: boolean) => void,
  label: string,
  isLoaded: boolean,
  itemCount: number,
  searchTerm: string,
): void => {
  const reportTimerRef = React.useRef<ReturnType<typeof setTimeout>>();

  React.useEffect(() => {
    clearTimeout(reportTimerRef.current);
    if (isLoaded && !searchTerm) {
      reportTimerRef.current = setTimeout(() => {
        reportCategoryEmpty(label, itemCount === 0);
      }, 100);
    }
    return () => clearTimeout(reportTimerRef.current);
  }, [isLoaded, itemCount, label, searchTerm, reportCategoryEmpty]);
};

export default useReportCategoryEmpty;
