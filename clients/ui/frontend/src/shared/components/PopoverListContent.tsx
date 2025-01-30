import * as React from 'react';
import { Content, ContentProps } from '@patternfly/react-core';

type PopoverListContentProps = ContentProps & {
  leadText?: React.ReactNode;
  listHeading?: React.ReactNode;
  listItems: React.ReactNode[];
};

const ContentText: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <Content component="small" style={{ color: 'var(--pf-t--global--text--color--regular)' }}>
    {children}
  </Content>
);

const PopoverListContent: React.FC<PopoverListContentProps> = ({
  leadText,
  listHeading,
  listItems,
  ...props
}) => (
  <Content {...props}>
    {leadText ? <ContentText>{leadText}</ContentText> : null}
    {listHeading ? <Content component="h4">{listHeading}</Content> : null}
    <Content component="ul" style={{ margin: 0 }}>
      {listItems.map((item, index) => (
        <Content component="li" key={index}>
          <ContentText>{item}</ContentText>
        </Content>
      ))}
    </Content>
  </Content>
);

export default PopoverListContent;
