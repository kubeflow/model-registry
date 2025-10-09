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
  onShowMore: (label: string, replace: boolean) => void;
  displayName?: string;
};

const CatalogCategorySection: React.FC<CategorySectionProps> = ({
  label,
  searchTerm,
  pageSize,
  catalogSources,
  onShowMore,
  displayName,
}) => {
  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    undefined,
    label,
    pageSize,
    searchTerm,
  );

  const handleShowMoreCategory = (categoryLabel: string) => {
    onShowMore(categoryLabel, true);
  };

  return (
    <>
      <StackItem>
        <Flex
          alignItems={{ default: 'alignItemsCenter' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
        >
          <FlexItem>
            <Title headingLevel="h3" size="lg" data-testid={label}>
              {`${displayName ?? label} models`}
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
                data-testid={`show-more-button-${label.toLowerCase().replace(/\s+/g, '-')}`}
                onClick={() => handleShowMoreCategory(label)}
              >
                Show all {displayName ?? label} models
              </Button>
            </FlexItem>
          )}
        </Flex>
      </StackItem>
      <StackItem className="pf-v6-u-pb-xl">
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
                data-testid={`category-skeleton-${label.toLowerCase().replace(/\s+/g, '-')}-${index}`}
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
