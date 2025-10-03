import {
  Alert,
  Button,
  Flex,
  FlexItem,
  Gallery,
  Skeleton,
  StackItem,
  Title,
} from '@patternfly/react-core';
import React from 'react';
import { ArrowRightIcon, SearchIcon } from '@patternfly/react-icons';
import { CatalogSourceList } from '~/app/modelCatalogTypes';
import { useCatalogModelsBySources } from '~/app/hooks/modelCatalog/useCatalogModelsBySource';
import EmptyModelCatalogState from '~/app/pages/modelCatalog/EmptyModelCatalogState';
import { getSourceFromSourceId } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';
import ModelCatalogCard from '~/app/pages/modelCatalog/components/ModelCatalogCard';

type CategorySectionProps = {
  label: string;
  searchTerm: string;
  pageSize: number;
  catalogSources: CatalogSourceList | null;
  onShowMore: (label: string) => void;
};

const CatalogCategorySection: React.FC<CategorySectionProps> = ({
  label,
  searchTerm,
  pageSize,
  catalogSources,
  onShowMore,
}) => {
  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    undefined,
    label,
    pageSize,
    searchTerm,
  );

  const handleShowMoreCategory = (categoryLabel: string) => {
    onShowMore(categoryLabel);
  };

  return (
    <>
      <StackItem>
        <Flex
          alignItems={{ default: 'alignItemsCenter' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
        >
          <FlexItem>
            <Title headingLevel="h3" size="lg">
              {`${label} models`}
            </Title>
          </FlexItem>
          {catalogModels.items.length >= 4 && (
            <FlexItem>
              <Button
                variant="link"
                size="sm"
                isInline
                icon={<ArrowRightIcon />}
                iconPosition="right"
                onClick={() => handleShowMoreCategory(label)}
              >
                Show all {label} models
              </Button>
            </FlexItem>
          )}
        </Flex>
      </StackItem>
      <StackItem style={{ paddingBottom: '20px' }}>
        {catalogModelsLoadError ? (
          <Alert variant="warning" title={`Failed to load ${label} models`}>
            {catalogModelsLoadError.message}
          </Alert>
        ) : !catalogModelsLoaded ? (
          <Gallery hasGutter minWidths={{ default: '300px' }}>
            {Array.from({ length: 4 }).map((_, index) => (
              <Skeleton
                key={index}
                height="280px"
                width="100%"
                screenreaderText={`Loading ${label} models`}
              />
            ))}
          </Gallery>
        ) : catalogModels.items.length === 0 ? (
          <EmptyModelCatalogState
            testid="empty-model-catalog-state"
            title="No result found"
            headerIcon={SearchIcon}
            description={<>Adjust your filters and try again.</>}
          />
        ) : (
          <>
            <Gallery hasGutter minWidths={{ default: '300px' }}>
              {catalogModels.items.slice(0, pageSize).map((model) => (
                <ModelCatalogCard
                  key={model.name}
                  model={model}
                  source={getSourceFromSourceId(model.source_id || '', catalogSources)}
                />
              ))}
            </Gallery>
          </>
        )}
      </StackItem>
    </>
  );
};
export default CatalogCategorySection;
