import * as React from 'react';

type ScrollViewOnMountProps = {
  shouldScroll: boolean;
  scrollToTop?: boolean;
};

const scrollAllAncestorsToTop = (node: HTMLElement | null) => {
  if (!node) {
    return;
  }
  node.scrollTo(0, 0);
  scrollAllAncestorsToTop(node.parentElement);
};

const ScrollViewOnMount: React.FC<ScrollViewOnMountProps> = ({
  shouldScroll,
  scrollToTop = false,
}) => {
  const ref = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (shouldScroll && ref.current) {
      if (scrollToTop) {
        scrollAllAncestorsToTop(ref.current);
      } else {
        ref.current.scrollIntoView();
      }
    }
  }, [scrollToTop, shouldScroll]);

  return <div ref={ref} />;
};

export default ScrollViewOnMount;
