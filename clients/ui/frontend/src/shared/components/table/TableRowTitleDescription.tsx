import * as React from 'react';
import { Content } from '@patternfly/react-core';
import MarkdownView from '~/shared/components/MarkdownView';

type TableRowTitleDescriptionProps = {
  title: React.ReactNode;
  // resource?: K8sResourceCommon; // TODO: Sort this out in refactor
  subtitle?: React.ReactNode;
  description?: string;
  descriptionAsMarkdown?: boolean;
  label?: React.ReactNode;
};

const TableRowTitleDescription: React.FC<TableRowTitleDescriptionProps> = ({
  title,
  description,
  // resource,
  subtitle,
  descriptionAsMarkdown,
  label,
}) => {
  let descriptionNode: React.ReactNode;
  if (description) {
    descriptionNode = descriptionAsMarkdown ? (
      <MarkdownView conciseDisplay markdown={description} />
    ) : (
      <Content
        component="p"
        data-testid="table-row-title-description"
        style={{ color: 'var(--pf-v6-global--Color--200)' }}
      >
        {description}
      </Content>
    );
  }

  return (
    <>
      <b data-testid="table-row-title">
        {/* {resource ? <ResourceNameTooltip resource={resource}>{title}</ResourceNameTooltip> : title} */}
        {title}
      </b>
      {subtitle}
      {descriptionNode}
      {label}
    </>
  );
};

export default TableRowTitleDescription;
