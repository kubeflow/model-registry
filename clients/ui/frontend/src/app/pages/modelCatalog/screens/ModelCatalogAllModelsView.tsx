import React from 'react';
import {
  Stack,
  StackItem,
  Title,
  Gallery,
  Button,
  Skeleton,
  Alert,
  Flex,
  FlexItem,
  Bullseye,
} from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import {
  CategoryModels,
  useCatalogAllCategoriesModels,
} from '~/app/hooks/modelCatalog/useAllCategoriesModels';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';
import { CatalogModel } from '~/app/modelCatalogTypes';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import { getSourceFromSourceId } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type ModelCatalogAllModelsViewProps = {
  searchTerm: string;
};

const ModelCatalogAllModelsView: React.FC<ModelCatalogAllModelsViewProps> = ({ searchTerm }) => {
  const { catalogSources, updateSelectedSourceLabel } = React.useContext(ModelCatalogContext);

  const { categoriesData, allCategoriesLoaded, isAnyLoading } = useCatalogAllCategoriesModels(
    searchTerm,
    4,
  );

  const handleShowMoreCategory = (categoryLabel: string) => {
    updateSelectedSourceLabel(categoryLabel);
  };

  const hasMoreModels = (categoryData: CategoryModels) => {
    if (!categoryData.models) {
      return false;
    }
    return categoryData.models.items.length >= 4;
  };

  return (
    <Stack hasGutter>
      {Object.entries(categoriesData).map(([categoryLabel, categoryData]) => (
        <React.Fragment key={categoryLabel}>
          <StackItem>
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              justifyContent={{ default: 'justifyContentSpaceBetween' }}
            >
              <FlexItem>
                <Title headingLevel="h3" size="lg">
                  {`${categoryLabel} models`}
                </Title>
              </FlexItem>

              {categoryData.loaded && hasMoreModels(categoryData) && (
                <FlexItem>
                  <Button
                    variant="link"
                    size="sm"
                    onClick={() => handleShowMoreCategory(categoryLabel)}
                  >
                    Show all {categoryLabel} models →
                  </Button>
                </FlexItem>
              )}
            </Flex>
          </StackItem>

          <StackItem style={{ paddingBottom: '20px' }}>
            {categoryData.error ? (
              <Alert variant="warning" title={`Failed to load ${categoryLabel} models`} isInline>
                {categoryData.error.message}
              </Alert>
            ) : !categoryData.loaded ? (
              <Gallery hasGutter minWidths={{ default: '300px' }}>
                {Array.from({ length: 4 }).map((_, index) => (
                  <Skeleton
                    key={index}
                    height="280px"
                    width="100%"
                    screenreaderText={`Loading ${categoryLabel} models`}
                  />
                ))}
              </Gallery>
            ) : categoryData.models?.items.length === 0 ? (
              <EmptyModelCatalogState
                testid="empty-model-catalog-state"
                title="No result found"
                headerIcon={SearchIcon}
                description={<>Adjust your filters and try again.</>}
              />
            ) : (
              <Gallery hasGutter minWidths={{ default: '300px' }}>
                {categoryData.models?.items.slice(0, 4).map((model: CatalogModel) => (
                  <ModelCatalogCard
                    key={`${model.name}/${model.source_id}`}
                    model={model}
                    source={
                      model.source_id
                        ? getSourceFromSourceId(model.source_id || '', catalogSources)
                        : undefined
                    }
                  />
                ))}
              </Gallery>
            )}
          </StackItem>
        </React.Fragment>
      ))}

      {!allCategoriesLoaded && isAnyLoading && (
        <StackItem>
          <Bullseye>Loading more categories...</Bullseye>
        </StackItem>
      )}
    </Stack>
  );
};

export default ModelCatalogAllModelsView;
