import React from 'react';
import { Link } from 'react-router-dom';
import { DashboardDescriptionListGroup } from 'mod-arch-shared';
import { ModelArtifact } from '~/app/types';
import { modelSourcePropertiesToCatalogParams } from '~/concepts/modelRegistry/utils';
import { getCatalogModelDetailsRoute } from '~/app/routes/modelCatalog/catalogModelDetails';

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
      {registeredFromCatalogDetails.modelName}
    </span>
  );

  const renderContent = () => {
    const catalogModelUrl = getCatalogModelDetailsRoute({
      modelName: registeredFromCatalogDetails.modelName,
      sourceId: registeredFromCatalogDetails.sourceId,
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
