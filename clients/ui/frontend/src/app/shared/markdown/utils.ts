import { ContentVariants } from '@patternfly/react-core';

const headingMap: Record<number, ContentVariants> = {
  1: ContentVariants.h1,
  2: ContentVariants.h2,
  3: ContentVariants.h3,
  4: ContentVariants.h4,
  5: ContentVariants.h5,
  6: ContentVariants.h6,
};

export const shiftHeadingLevel = (level: number, maxHeading: number): ContentVariants => {
  const adjusted = Math.min(level + maxHeading - 1, 6);
  return headingMap[adjusted];
};
