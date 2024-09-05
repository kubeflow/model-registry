import * as React from 'react';
import { GalleryItemProps } from '@patternfly/react-core';
import { css } from '@patternfly/react-styles';

import './DividedGallery.scss';

const DividedGalleryItem: React.FC<GalleryItemProps> = ({ className, ...rest }) => (
  <div className={css('kubeflowdivided-gallery__item', className)} {...rest} />
);

export default DividedGalleryItem;
