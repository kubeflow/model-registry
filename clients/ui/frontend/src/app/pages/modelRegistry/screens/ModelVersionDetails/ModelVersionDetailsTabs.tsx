import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { PageSection, Tab, Tabs, TabTitleText } from '@patternfly/react-core';
import {
  DeploymentMode,
  FetchStateObject,
  InferenceServiceKind,
  ServingRuntimeKind,
  useModularArchContext,
} from 'mod-arch-shared';
import { ModelVersion } from '~/app/types';
import { ModelVersionDetailsTabTitle, ModelVersionDetailsTab } from './const';
import ModelVersionDetailsView from './ModelVersionDetailsView';
import ModelVersionRegisteredDeploymentsView from './ModelVersionRegisteredDeploymentsView';

type ModelVersionDetailTabsProps = {
  tab: ModelVersionDetailsTab;
  modelVersion: ModelVersion;
  inferenceServices: FetchStateObject<InferenceServiceKind[]>;
  servingRuntimes: FetchStateObject<ServingRuntimeKind[]>;
  isArchiveVersion?: boolean;
  refresh: () => void;
};

const ModelVersionDetailsTabs: React.FC<ModelVersionDetailTabsProps> = ({
  tab,
  modelVersion: mv,
  inferenceServices,
  servingRuntimes,
  isArchiveVersion,
  refresh,
}) => {
  const navigate = useNavigate();
  const { config } = useModularArchContext();
  const { deploymentMode } = config;
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
          />
        </PageSection>
      </Tab>
      {!isArchiveVersion &&
        (deploymentMode === DeploymentMode.Standalone ||
          deploymentMode === DeploymentMode.Federated) && (
          <Tab
            eventKey={ModelVersionDetailsTab.DEPLOYMENTS}
            title={<TabTitleText>{ModelVersionDetailsTabTitle.DEPLOYMENTS}</TabTitleText>}
            aria-label="Deployments tab"
            data-testid="deployments-tab"
          >
            <PageSection hasBodyWrapper={false} isFilled data-testid="deployments-tab-content">
              <ModelVersionRegisteredDeploymentsView
                inferenceServices={inferenceServices}
                servingRuntimes={servingRuntimes}
                refresh={refresh}
              />
            </PageSection>
          </Tab>
        )}
    </Tabs>
  );
};

export default ModelVersionDetailsTabs;
