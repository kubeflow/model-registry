import React from 'react';
import { useParams, useNavigate } from 'react-router';
import { Link } from 'react-router-dom';
import {
  ActionList,
  Breadcrumb,
  BreadcrumbItem,
  Content,
  ContentVariants,
  Flex,
  FlexItem,
  Stack,
  StackItem,
  Button,
  Popover,
  ActionListGroup,
  Skeleton,
} from '@patternfly/react-core';
import { ApplicationsPage } from 'mod-arch-shared';
import { decodeParams, getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelDetailsTabs from '~/app/pages/modelCatalog/screens/ModelDetailsTabs';
import { useCatalogModel } from '~/app/hooks/modelCatalog/useCatalogModel';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { getRegisterCatalogModelRoute } from '~/app/routes/modelCatalog/catalogModelRegister';
import { CatalogModelDetailsParams } from '~/app/modelCatalogTypes';
import { useCatalogModelArtifacts } from '~/app/hooks/modelCatalog/useCatalogModelArtifacts';

const ModelDetailsPage: React.FC = () => {
  const params = useParams<CatalogModelDetailsParams>();
  const decodedParams = decodeParams(params);
  const navigate = useNavigate();

  const state = useCatalogModel(
    decodedParams.sourceId || '',
    encodeURIComponent(`${decodedParams.modelName}`),
  );
  const [model, modelLoaded, modelLoadError] = state;
  const { modelRegistries, modelRegistriesLoadError, modelRegistriesLoaded } = React.useContext(
    ModelRegistrySelectorContext,
  );

  const [artifacts, artifactLoaded, artifactsLoadError] = useCatalogModelArtifacts(
    decodedParams.sourceId || '',
    encodeURIComponent(encodeURIComponent(`${decodedParams.modelName}`)) || '',
  );

  const registerButtonPopover = (headerContent: string, bodyContent: string) => (
    <Popover
      headerContent={headerContent}
      triggerAction="hover"
      data-testid="register-catalog-model-popover"
      bodyContent={<div>{bodyContent}</div>}
    >
      <Button variant="primary" isAriaDisabled data-testid="register-model-button">
        Register model
      </Button>
    </Popover>
  );

  const registerModelButton = () => {
    if (!modelRegistriesLoaded || modelRegistriesLoadError) {
      return null;
    }

    if (artifactsLoadError) {
      return registerButtonPopover(
        'Unable to load model artifacts',
        'Model registration is unavailable due to an error loading model artifacts. Please try again later.',
      );
    }

    if (!artifactLoaded) {
      return (
        <Button variant="primary" data-testid="register-model-button" isLoading>
          Register model
        </Button>
      );
    }

    return modelRegistries.length === 0 ? (
      registerButtonPopover(
        'Request access to a model registry',
        'To request a new model registry, or to request permission to access an existing model registry, contact your administrator.',
      )
    ) : artifacts.items.length === 0 ? (
      registerButtonPopover('', 'Model location is unavailable')
    ) : (
      <Button
        data-testid="register-model-button"
        variant="primary"
        onClick={() => {
          navigate(getRegisterCatalogModelRoute(decodedParams.sourceId, decodedParams.modelName));
        }}
      >
        Register model
      </Button>
    );
  };

  return (
    <ApplicationsPage
      breadcrumb={
        <Breadcrumb>
          <BreadcrumbItem>
            <Link to="/model-catalog">Model catalog</Link>
          </BreadcrumbItem>
          <BreadcrumbItem isActive>{getModelName(model?.name || '') || 'Details'}</BreadcrumbItem>
        </Breadcrumb>
      }
      title={
        model ? (
          <Flex
            spaceItems={{ default: 'spaceItemsMd' }}
            alignItems={{ default: 'alignItemsCenter' }}
          >
            {model.logo ? (
              <img src={model.logo} alt="model logo" style={{ height: '56px', width: '56px' }} />
            ) : (
              <Skeleton
                shape="square"
                width="56px"
                height="56px"
                screenreaderText="Brand image loading"
              />
            )}
            <Stack>
              <StackItem>
                <Flex
                  spaceItems={{ default: 'spaceItemsSm' }}
                  alignItems={{ default: 'alignItemsCenter' }}
                >
                  <FlexItem>{getModelName(model.name)}</FlexItem>
                </Flex>
              </StackItem>
              <StackItem>
                <Content component={ContentVariants.small}>Provided by {model.provider}</Content>
              </StackItem>
            </Stack>
          </Flex>
        ) : (
          'Model details'
        )
      }
      empty={!model}
      emptyStatePage={
        !model ? (
          <div>
            Details not found. Return to <Link to="/model-catalog">Model catalog</Link>
          </div>
        ) : undefined
      }
      loadError={modelLoadError}
      loaded={modelLoaded}
      errorMessage="Unable to load model catalog"
      provideChildrenPadding
      headerAction={
        modelLoaded &&
        !modelLoadError &&
        model && (
          <ActionList>
            <ActionListGroup>{registerModelButton()}</ActionListGroup>
          </ActionList>
        )
      }
    >
      {model && <ModelDetailsTabs model={model} decodedParams={decodedParams} />}
    </ApplicationsPage>
  );
};

export default ModelDetailsPage;
