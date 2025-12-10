import {
  Alert,
  Button,
  Flex,
  FlexItem,
  Grid,
  GridItem,
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

  const itemsToDisplay = catalogModels.items.slice(0, pageSize);

  return (
    <>
      <StackItem className="pf-v6-u-pb-xl">
        <Flex
          alignItems={{ default: 'alignItemsCenter' }}
          justifyContent={{ default: 'justifyContentSpaceBetween' }}
          className="pf-v6-u-mb-md"
        >
          <FlexItem>
            <Title headingLevel="h3" size="lg" data-testid={`title ${label}`}>
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
                data-testid={`show-more-button ${label.toLowerCase().replace(/\s+/g, '-')}`}
                onClick={() => onShowMore(label)}
              >
                Show all {displayName ?? label} models
              </Button>
            </FlexItem>
          )}
        </Flex>

        {catalogModelsLoadError ? (
          <Alert
            variant="danger"
            title={`Failed to load ${displayName ?? label} models`}
            data-testid={`error-state ${label}`}
          >
            {catalogModelsLoadError.message}
          </Alert>
        ) : !catalogModelsLoaded ? (
          <Grid hasGutter>
            {Array.from({ length: 4 }).map((_, index) => (
              <GridItem key={index} sm={6} md={6} lg={6} xl={6} xl2={3}>
                <Skeleton
                  height="280px"
                  width="100%"
                  screenreaderText={`Loading ${label} models`}
                  data-testid={`category-skeleton-${label.toLowerCase().replace(/\s+/g, '-')}-${index}`}
                />
              </GridItem>
            ))}
          </Grid>
        ) : catalogModels.items.length === 0 ? (
          <EmptyModelCatalogState
            testid={`empty-model-catalog-state ${label}`}
            title="No result found"
            headerIcon={SearchIcon}
            description="Adjust your filters and try again."
          />
        ) : (
          <Grid hasGutter>
            {itemsToDisplay.map((model) => (
              <GridItem
                key={`${model.name}/${model.source_id}`}
                sm={6}
                md={6}
                lg={6}
                xl={6}
                xl2={3}
              >
                <ModelCatalogCard
                  model={model}
                  source={getSourceFromSourceId(model.source_id || '', catalogSources)}
                />
              </GridItem>
            ))}
          </Grid>
        )}
      </StackItem>
    </>
  );
};
export default CatalogCategorySection;
