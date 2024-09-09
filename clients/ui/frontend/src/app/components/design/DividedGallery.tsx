import * as React from 'react';
import { Button, Gallery, GalleryProps } from '@patternfly/react-core';
import { css } from '@patternfly/react-styles';
import { TimesIcon } from '@patternfly/react-icons';

type DividedGalleryProps = Omit<GalleryProps, 'minWidths' | 'maxWidths'> & {
  minSize: string;
  itemCount: number;
  showClose?: boolean;
  closeAlt?: string;
  onClose?: () => void;
  closeTestId?: string;
};

import './DividedGallery.scss';

const DividedGallery: React.FC<DividedGalleryProps> = ({
  minSize,
  itemCount,
  showClose,
  closeAlt,
  onClose,
  children,
  className,
  closeTestId,
  ...rest
}) => (
  <div className={css('kubeflowdivided-gallery', className)} {...rest}>
    <Gallery
      minWidths={{ default: minSize, md: minSize }}
      maxWidths={{ default: '100%', md: `${100 / itemCount}%` }}
    >
      <div className="kubeflowdivided-gallery__border" />
      {children}
      {showClose ? (
        <div className="kubeflowdivided-gallery__close">
          <Button
            data-testid={closeTestId}
            aria-label={closeAlt || 'close'}
            isInline
            variant="plain"
            onClick={onClose}
          >
            <TimesIcon alt={`close ${closeAlt}`} />
          </Button>
        </div>
      ) : null}
    </Gallery>
  </div>
);

export default DividedGallery;
