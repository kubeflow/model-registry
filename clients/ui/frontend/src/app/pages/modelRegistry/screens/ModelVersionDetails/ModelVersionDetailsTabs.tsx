import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { PageSection, Tab, Tabs, TabTitleText } from '@patternfly/react-core';
import { ModelVersion, ModelArtifactList } from '~/app/types';
import { ModelVersionDetailsTabTitle, ModelVersionDetailsTab } from './const';
import ModelVersionDetailsView from './ModelVersionDetailsView';

type ModelVersionDetailTabsProps = {
  tab: ModelVersionDetailsTab;
  modelVersion: ModelVersion;
  isArchiveVersion?: boolean;
  refresh: () => void;
  modelArtifacts: ModelArtifactList;
};

const ModelVersionDetailsTabs: React.FC<ModelVersionDetailTabsProps> = ({
  tab,
  modelVersion: mv,
  isArchiveVersion,
  refresh,
  modelArtifacts,
}) => {
  const navigate = useNavigate();
  return (
    <Tabs
      activeKey={tab}
      aria-label="Model versions details page tabs"
      role="region"
      data-testid="model-versions-details-page-tabs"
      onSelect={(_event, eventKey) => navigate(`../${eventKey}`, { relative: 'path' })}
    >
      <Tab
        eventKey={ModelVersionDetailsTab.DETAILS}
        title={<TabTitleText>{ModelVersionDetailsTabTitle.DETAILS}</TabTitleText>}
        aria-label="Model versions details tab"
        data-testid="model-versions-details-tab"
      >
        <PageSection
          hasBodyWrapper={false}
          isFilled
          data-testid="model-versions-details-tab-content"
        >
          <ModelVersionDetailsView
            modelVersion={mv}
            refresh={refresh}
            isArchiveVersion={isArchiveVersion}
            modelArtifacts={modelArtifacts}
          />
        </PageSection>
      </Tab>
    </Tabs>
  );
};

export default ModelVersionDetailsTabs;
