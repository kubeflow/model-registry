import * as React from 'react';
import { Button, Flex, FlexItem, Content, ContentVariants } from '@patternfly/react-core';
import { AngleDownIcon, AngleRightIcon } from '@patternfly/react-icons';

interface CollapsibleSectionProps {
  open?: boolean;
  setOpen?: (update: boolean) => void;
  title: string;
  titleVariant?: ContentVariants.h1 | ContentVariants.h2;
  children?: React.ReactNode;
  id?: string;
  showChildrenWhenClosed?: boolean;
}

const CollapsibleSection: React.FC<CollapsibleSectionProps> = ({
  open,
  setOpen,
  title,
  titleVariant = ContentVariants.h2,
  children,
  id,
  showChildrenWhenClosed,
}) => {
  const [innerOpen, setInnerOpen] = React.useState<boolean>(true);
  const localId = id || title.replace(/ /g, '-');
  const titleId = `${localId}-title`;

  return (
    <>
      <Flex
        gap={{ default: 'gapMd' }}
        alignItems={{ default: 'alignItemsCenter' }}
        style={
          (open ?? innerOpen) || showChildrenWhenClosed
            ? {
                marginBottom: 'var(--pf-t--global--spacer--md)',
              }
            : undefined
        }
      >
        <FlexItem>
          <Button
            icon={(open ?? innerOpen) ? <AngleDownIcon /> : <AngleRightIcon />}
            aria-labelledby={titleId}
            aria-expanded={open}
            variant="plain"
            style={{
              paddingLeft: 0,
              paddingRight: 0,
              fontSize:
                titleVariant === ContentVariants.h2
                  ? 'var(--pf-v6-global--FontSize--xl)'
                  : 'var(--pf-v6-global--FontSize--2xl)',
            }}
            onClick={() => (setOpen ? setOpen(!open) : setInnerOpen((prev) => !prev))}
          />
        </FlexItem>
        <FlexItem>
          <Content>
            <Content id={titleId} component={titleVariant}>
              {title}
            </Content>
          </Content>
        </FlexItem>
      </Flex>
      {(open ?? innerOpen) || showChildrenWhenClosed ? children : null}
    </>
  );
};

export default CollapsibleSection;
