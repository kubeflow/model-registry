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
import { getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelDetailsView from '~/app/pages/modelCatalog/screens/ModelDetailsView';
import { useCatalogModel } from '~/app/hooks/modelCatalog/useCatalogModel';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { getRegisterCatalogModelRoute } from '~/app/routes/modelCatalog/catalogModelRegister';

type RouteParams = {
  sourceId: string;
  modelName: string;
  repositoryName: string;
};

const ModelDetailsPage: React.FC = () => {
  const { sourceId, repositoryName, modelName } = useParams<RouteParams>();
  const navigate = useNavigate();

  const state = useCatalogModel(sourceId || '', `${repositoryName}/${modelName}` || '');
  const [model, modelLoaded, modelLoadError] = state;
  const { modelRegistries, modelRegistriesLoadError, modelRegistriesLoaded } = React.useContext(
    ModelRegistrySelectorContext,
  );

  // TODO: we don't have tags prop on models
  // const versionTag = extractVersionTag(model?.tags);

  const registerModelButton = () => {
    if (!modelRegistriesLoaded || modelRegistriesLoadError) {
      return null;
    }

    return modelRegistries.length === 0 ? (
      <Popover
        headerContent="Request access to a model registry"
        triggerAction="hover"
        data-testid="register-catalog-model-popover"
        bodyContent={
          <div>
            To request a new model registry, or to request permission to access an existing model
            registry, contact your administrator.
          </div>
        }
      >
        <Button variant="primary" isAriaDisabled data-testid="register-model-button">
          Register model
        </Button>
      </Popover>
    ) : (
      <Button
        data-testid="register-model-button"
        variant="primary"
        onClick={() => {
          if (sourceId) {
            navigate(getRegisterCatalogModelRoute(sourceId));
          }
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
      empty={modelLoaded && !modelLoadError && !model}
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
        !modelLoaded &&
        !modelLoadError &&
        model && (
          <ActionList>
            <ActionListGroup>{registerModelButton()}</ActionListGroup>
          </ActionList>
        )
      }
    >
      {model && <ModelDetailsView model={model} />}
    </ApplicationsPage>
  );
};

export default ModelDetailsPage;
