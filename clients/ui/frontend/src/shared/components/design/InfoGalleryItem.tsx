import * as React from 'react';
import {
  Button,
  ButtonVariant,
  Flex,
  FlexItem,
  GalleryItemProps,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { SectionType, sectionTypeBackgroundColor } from '~/shared/components/design/utils';
import DividedGalleryItem from '~/shared/components/design/DividedGalleryItem';

const HEADER_ICON_SIZE = 40;
const HEADER_ICON_PADDING = 2;

type InfoGalleryItemProps = {
  title: string;
  imgSrc: string;
  sectionType: SectionType;
  description: React.ReactNode;
  isOpen: boolean;
  onClick?: () => void;
  testId?: string;
} & GalleryItemProps;

const InfoGalleryItem: React.FC<InfoGalleryItemProps> = ({
  title,
  imgSrc,
  sectionType,
  description,
  isOpen,
  onClick,
  testId,
  ...rest
}) => (
  <DividedGalleryItem data-testid={testId} {...rest}>
    <Stack hasGutter>
      <StackItem>
        <Flex
          gap={{ default: 'gapMd' }}
          direction={{ default: isOpen ? 'column' : 'row' }}
          alignItems={{ default: isOpen ? 'alignItemsFlexStart' : 'alignItemsCenter' }}
        >
          <FlexItem
            style={{
              display: 'inline-block',
              width: HEADER_ICON_SIZE,
              height: HEADER_ICON_SIZE,
              padding: HEADER_ICON_PADDING,
              borderRadius: HEADER_ICON_SIZE / 2,
              background: sectionTypeBackgroundColor(sectionType),
            }}
          >
            <img
              width={HEADER_ICON_SIZE - HEADER_ICON_PADDING * 2}
              height={HEADER_ICON_SIZE - HEADER_ICON_PADDING * 2}
              src={imgSrc}
              alt=""
            />
          </FlexItem>
          {onClick ? (
            <Button
              data-testid={testId ? `${testId}-button` : undefined}
              variant={ButtonVariant.link}
              isInline
              onClick={onClick}
              style={{
                fontSize: 'var(--pf-v6-global--FontSize--md)',
                fontWeight: 'var(--pf-v6-global--FontWeight--bold)',
              }}
            >
              {title}
            </Button>
          ) : (
            <FlexItem>{title}</FlexItem>
          )}
        </Flex>
      </StackItem>
      {isOpen ? <StackItem isFilled>{description}</StackItem> : null}
    </Stack>
  </DividedGalleryItem>
);

export default InfoGalleryItem;
