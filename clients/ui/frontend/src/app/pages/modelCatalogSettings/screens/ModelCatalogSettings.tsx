import * as React from 'react';
import { Button, Flex, FlexItem } from '@patternfly/react-core';
import { PlusCircleIcon } from '@patternfly/react-icons';
import { useNavigate } from 'react-router-dom';
import { ProjectObjectType, TitleWithIcon, ApplicationsPage } from 'mod-arch-shared';
import {
  CATALOG_SETTINGS_PAGE_TITLE,
  CATALOG_SETTINGS_DESCRIPTION,
  addSourceUrl,
} from '~/app/routes/modelCatalogSettings/modelCatalogSettings';

const ModelCatalogSettings: React.FC = () => {
  const navigate = useNavigate();

  return (
    <ApplicationsPage
      title={
        <TitleWithIcon
          title={CATALOG_SETTINGS_PAGE_TITLE}
          objectType={ProjectObjectType.modelCatalog}
        />
      }
      description={CATALOG_SETTINGS_DESCRIPTION}
      empty={false}
      loaded
      provideChildrenPadding
    >
      <Flex>
        <FlexItem>
          <Button
            variant="primary"
            icon={<PlusCircleIcon />}
            onClick={() => navigate(addSourceUrl())}
            data-testid="add-source-button"
          >
            Add a source
          </Button>
        </FlexItem>
      </Flex>
    </ApplicationsPage>
  );
};

export default ModelCatalogSettings;
