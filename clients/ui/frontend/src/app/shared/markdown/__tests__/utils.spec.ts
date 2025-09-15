import { ContentVariants } from '@patternfly/react-core';
import { shiftHeadingLevel } from '~/app/shared/markdown/utils';

describe('shiftHeadingLevel', () => {
  it('returns correct heading when within range', () => {
    expect(shiftHeadingLevel(1, 1)).toBe(ContentVariants.h1);
    expect(shiftHeadingLevel(2, 1)).toBe(ContentVariants.h2);
    expect(shiftHeadingLevel(3, 2)).toBe(ContentVariants.h4);
  });

  it('caps heading at h6 when adjusted value exceeds 6', () => {
    expect(shiftHeadingLevel(5, 3)).toBe(ContentVariants.h6);
    expect(shiftHeadingLevel(6, 2)).toBe(ContentVariants.h6);
  });
});
