import * as React from 'react';

type ScrolledGalleryProps = {
  count: number;
  childWidth: string;
} & React.HTMLProps<HTMLDivElement>;

const ScrolledGallery: React.FC<ScrolledGalleryProps> = ({
  children,
  count,
  childWidth,
  ...rest
}) => {
  let gridTemplateColumns = childWidth;
  for (let i = 1; i < count; i++) {
    gridTemplateColumns = `${gridTemplateColumns} ${childWidth}`;
  }
  return (
    <div
      style={{
        gridTemplateColumns,
        display: 'grid',
        gridAutoFlow: 'column',
        overflowY: 'auto',
        gap: 'var(--pf-v6-global--spacer--md)',
        paddingBottom: 'var(--pf-v6-global--spacer--sm)',
      }}
      {...rest}
    >
      {children}
    </div>
  );
};

export default ScrolledGallery;
