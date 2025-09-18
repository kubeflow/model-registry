import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { PageSection, Tab, Tabs, TabTitleText } from '@patternfly/react-core';
import { ModelVersion, ModelArtifactList, RegisteredModel } from '~/app/types';
import { ModelVersionDetailsTabTitle, ModelVersionDetailsTab } from './const';
import ModelVersionDetailsView from './ModelVersionDetailsView';

type ModelVersionDetailTabsProps = {
  tab: ModelVersionDetailsTab;
  registeredModel: RegisteredModel | null;
  modelVersion: ModelVersion;
  isArchiveVersion?: boolean;
  refresh: () => void;
  modelArtifacts: ModelArtifactList;
  modelArtifactsLoaded: boolean;
  modelArtifactsLoadError: Error | undefined;
};

const ModelVersionDetailsTabs: React.FC<ModelVersionDetailTabsProps> = ({
  tab,
  registeredModel,
  modelVersion: mv,
  isArchiveVersion,
  refresh,
  modelArtifacts,
  modelArtifactsLoaded,
  modelArtifactsLoadError,
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
            registeredModel={registeredModel}
            modelVersion={mv}
            refresh={refresh}
            isArchiveVersion={isArchiveVersion}
            modelArtifacts={modelArtifacts}
            modelArtifactsLoaded={modelArtifactsLoaded}
            modelArtifactsLoadError={modelArtifactsLoadError}
          />
        </PageSection>
      </Tab>
    </Tabs>
  );
};

export default ModelVersionDetailsTabs;
