import * as React from 'react';

type ScrollViewOnMountProps = {
  shouldScroll: boolean;
  scrollToTop?: boolean;
};

const ScrollViewOnMount: React.FC<ScrollViewOnMountProps> = ({
  shouldScroll,
  scrollToTop = false,
}) => {
  const ref = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (shouldScroll && ref.current) {
      if (scrollToTop) {
        ref.current.scrollIntoView();
      } else {
        window.scrollTo(0, 0);
      }
    }
  }, [scrollToTop, shouldScroll]);

  return <div ref={ref} />;
};

export default ScrollViewOnMount;
