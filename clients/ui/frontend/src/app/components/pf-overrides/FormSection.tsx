import * as React from 'react';
import { FormSection as PFFormSection, FormSectionProps, Content } from '@patternfly/react-core';

import './FormSection.scss';

type Props = FormSectionProps & {
  description?: React.ReactNode;
};

// Remove once https://github.com/patternfly/patternfly/issues/6663 is fixed
const FormSection: React.FC<Props> = ({
  description,
  title,
  titleElement: TitleElement = 'div',
  ...props
}) => (
  <PFFormSection
    {...props}
    titleElement={description ? 'div' : TitleElement}
    title={
      description ? (
        <>
          <TitleElement className="pf-v6-c-form__section-title">{title}</TitleElement>
          <Content component="p" className="kf-form-section__desc">
            {description}
          </Content>
        </>
      ) : (
        title
      )
    }
  />
);

export default FormSection;
