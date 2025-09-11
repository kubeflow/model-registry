import React from 'react';
import { Link } from 'react-router-dom';
import { DashboardDescriptionListGroup } from 'mod-arch-shared';
import { ModelArtifact } from '~/app/types';
import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';
import { modelSourcePropertiesToCatalogParams } from '~/concepts/modelRegistry/utils';

type ModelVersionRegisteredFromLinkProps = {
  modelArtifact: ModelArtifact;
  isModelCatalogAvailable: boolean;
};

const ModelVersionRegisteredFromLink: React.FC<ModelVersionRegisteredFromLinkProps> = ({
  modelArtifact,
  isModelCatalogAvailable,
}) => {
  const registeredFromCatalogDetails = modelSourcePropertiesToCatalogParams(modelArtifact);

  if (!registeredFromCatalogDetails) {
    return null;
  }

  const registeredfromText = (
    <span className="pf-v6-u-font-weight-bold" data-testid="registered-from-catalog">
      {registeredFromCatalogDetails.modelName} ({registeredFromCatalogDetails.tag})
    </span>
  );

  const renderContent = () => {
    const catalogModelUrl = getCatalogModelDetailsRoute({
      modelName: registeredFromCatalogDetails.modelName || '',
      tag: registeredFromCatalogDetails.tag || '',
      sourceName: registeredFromCatalogDetails.sourceName,
      repositoryName: registeredFromCatalogDetails.repositoryName,
    });
    return (
      <>
        {isModelCatalogAvailable ? (
          <Link to={catalogModelUrl}>{registeredfromText}</Link>
        ) : (
          registeredfromText
        )}{' '}
        in Model catalog
      </>
    );
  };

  const content = renderContent();

  return (
    <DashboardDescriptionListGroup title="Registered from" groupTestId="registered-from-title">
      {content}
    </DashboardDescriptionListGroup>
  );
};

export default ModelVersionRegisteredFromLink;
