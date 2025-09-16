import * as React from 'react';

type ScrollViewOnMountProps = {
  shouldScroll: boolean;
};

const ScrollViewOnMount: React.FC<ScrollViewOnMountProps> = ({ shouldScroll }) => {
  const ref = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (shouldScroll && ref.current) {
      ref.current.scrollIntoView();
    }
  }, [shouldScroll]);

  return <div ref={ref} />;
};

export default ScrollViewOnMount;
