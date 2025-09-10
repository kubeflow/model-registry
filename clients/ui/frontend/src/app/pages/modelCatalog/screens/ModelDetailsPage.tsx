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
  Label,
  Stack,
  StackItem,
  Button,
  Popover,
  ActionListGroup,
  Skeleton,
} from '@patternfly/react-core';
import { TagIcon } from '@patternfly/react-icons';
import { ApplicationsPage } from 'mod-arch-shared';
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import { extractVersionTag } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
import { getRegisterCatalogModelRoute } from '~/app/routes/modelCatalog/catalogModelRegister';
import ModelDetailsView from './ModelDetailsView';

type RouteParams = {
  modelId: string;
};

const ModelDetailsPage: React.FC = () => {
  const { modelId } = useParams<RouteParams>();
  const navigate = useNavigate();
  const { sources, loading, error } = useModelCatalogSources();
  const { modelRegistries, modelRegistriesLoadError, modelRegistriesLoaded } = React.useContext(
    ModelRegistrySelectorContext,
  );

  const model = React.useMemo(() => {
    for (const source of sources) {
      const found = source.models?.find((m) => m.id === modelId);
      if (found) {
        return found;
      }
    }
    return undefined;
  }, [sources, modelId]);

  const versionTag = extractVersionTag(model?.tags);

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
          if (modelId) {
            navigate(getRegisterCatalogModelRoute(modelId));
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
          <BreadcrumbItem isActive>{model?.name || 'Details'}</BreadcrumbItem>
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
                  <FlexItem>{model.name}</FlexItem>
                  <Label variant="outline" icon={<TagIcon />}>
                    {versionTag || 'N/A'}
                  </Label>
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
      empty={!loading && !error && !model}
      emptyStatePage={
        !model ? (
          <div>
            Details not found. Return to <Link to="/model-catalog">Model catalog</Link>
          </div>
        ) : undefined
      }
      loadError={error}
      loaded={!loading}
      errorMessage="Unable to load model catalog"
      provideChildrenPadding
      headerAction={
        !loading &&
        !error &&
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
