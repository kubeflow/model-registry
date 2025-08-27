import * as React from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  Breadcrumb,
  BreadcrumbItem,
  Content,
  ContentVariants,
  Flex,
  FlexItem,
  Label,
  Skeleton,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { TagIcon } from '@patternfly/react-icons';
import { ApplicationsPage } from 'mod-arch-shared';
import { useModelCatalogSources } from '~/app/hooks/modelCatalog/useModelCatalogSources';
import { extractVersionTag } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelDetailsView from '~/app/pages/modelCatalog/screens/ModelDetailsView';

type RouteParams = {
  modelId: string;
};

const ModelDetailsRoute: React.FC = () => {
  const { modelId } = useParams<RouteParams>();
  const { sources, loading, error } = useModelCatalogSources();

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
              <img src={model.logo} alt="model logo" style={{ height: '40px', width: '40px' }} />
            ) : (
              <Skeleton
                shape="square"
                width="40px"
                height="40px"
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
    >
      {model && <ModelDetailsView model={model} />}
    </ApplicationsPage>
  );
};

export default ModelDetailsRoute;
