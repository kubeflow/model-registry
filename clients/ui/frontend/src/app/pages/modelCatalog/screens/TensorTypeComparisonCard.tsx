import {
  Alert,
  Card,
  CardBody,
  Divider,
  Flex,
  FlexItem,
  Label,
  Skeleton,
  Spinner,
  Title,
  Truncate,
} from '@patternfly/react-core';
import React from 'react';
import { Link } from 'react-router';
import { useCatalogModelsBySources } from '~/app/hooks/modelCatalog/useCatalogModelsBySource';
import { CatalogModel } from '~/app/modelCatalogTypes';
import { catalogModelDetailsFromModel } from '~/app/routes/modelCatalog/catalogModel';
import { getStringValue } from '~/app/utils';
import { getModelName } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type TensorTypeComparisonCardProps = {
  model: CatalogModel;
};
const TensorTypeComparisonCard: React.FC<TensorTypeComparisonCardProps> = ({ model }) => {
  const variantGroupId = getStringValue(model.customProperties ?? {}, 'variant_group_id');

  const variantFilterQuery = `variant_group_id.string_value="${variantGroupId}"`;

  const { catalogModels, catalogModelsLoaded, catalogModelsLoadError } = useCatalogModelsBySources(
    model.source_id || '',
    undefined,
    10,
    '',
    undefined,
    undefined,
    variantFilterQuery,
  );

  const sortedModels = React.useMemo(() => {
    if (!catalogModels.items.length) {
      return [];
    }

    return [...catalogModels.items]
      .toSorted((a, b) => {
        if (a.name === model.name) {
          return -1;
        }
        if (b.name === model.name) {
          return 1;
        }
        return 0;
      })
      .slice(0, 4);
  }, [catalogModels.items, model.name]);

  if (!variantGroupId) {
    return null;
  }

  return (
    <Card data-testid="compression-comparison-card">
      <CardBody>
        <Flex direction={{ default: 'column' }} gap={{ default: 'gapLg' }}>
          <FlexItem>
            <Flex direction={{ default: 'column' }} gap={{ default: 'gapSm' }}>
              <FlexItem>
                <Title headingLevel="h2" size="lg">
                  Compression level comparison
                </Title>
              </FlexItem>
              <FlexItem>
                <p>View benchmark performance of this model&apos;s available compression levels.</p>
              </FlexItem>
            </Flex>
          </FlexItem>
          <FlexItem style={{ overflowX: 'auto' }}>
            {catalogModelsLoadError ? (
              <Alert
                variant="danger"
                isInline
                title="Error loading performance data"
                data-testid="compression-comparison-error"
              >
                {catalogModelsLoadError.message || 'An error occurred'}
              </Alert>
            ) : !catalogModelsLoaded ? (
              <Spinner size="lg" data-testid="compression-comparison-loading" />
            ) : sortedModels.length === 0 ? (
              <Alert
                variant="info"
                isInline
                title="No compression variants found"
                data-testid="compression-comparison-empty"
              />
            ) : (
              <Flex
                gap={{ default: 'gapMd' }}
                flexWrap={{ default: 'nowrap' }}
                justifyContent={{ default: 'justifyContentSpaceEvenly' }}
              >
                {sortedModels.map((variant, index) => {
                  const tensorType = getStringValue(variant.customProperties ?? {}, 'tensor_type');
                  const isCurrent = variant.name === model.name;
                  const modelDisplayName = getModelName(variant.name || '');

                  return (
                    <React.Fragment key={`${variant.name}-${index}`}>
                      {index > 0 && (
                        <Divider
                          orientation={{ default: 'vertical' }}
                          data-testid={`compression-divider-${index}`}
                        />
                      )}
                      <FlexItem
                        data-testid={`compression-variant-${index}`}
                        style={{
                          minWidth: '180px',
                        }}
                      >
                        <Flex
                          alignItems={{ default: 'alignItemsFlexStart' }}
                          gap={{ default: 'gapSm' }}
                          flexWrap={{ default: 'nowrap' }}
                        >
                          <FlexItem flex={{ default: 'flexNone' }}>
                            {variant.logo ? (
                              <img
                                src={variant.logo}
                                alt="model logo"
                                style={{ height: '56px', width: '56px' }}
                                data-testid={`compression-logo-${index}`}
                              />
                            ) : (
                              <Skeleton
                                shape="square"
                                width="56px"
                                height="56px"
                                screenreaderText="Brand image loading"
                                data-testid={`compression-variant-skeleton-${index}`}
                              />
                            )}
                          </FlexItem>
                          <FlexItem>
                            <Flex
                              direction={{ default: 'column' }}
                              spaceItems={{ default: 'spaceItemsXs' }}
                            >
                              <FlexItem>
                                {isCurrent ? (
                                  <Truncate
                                    content={modelDisplayName}
                                    position="middle"
                                    maxCharsDisplayed={20}
                                    tooltipPosition="top"
                                    data-testid="compression-current-model-name"
                                  />
                                ) : (
                                  <Link
                                    to={catalogModelDetailsFromModel(
                                      encodeURIComponent(variant.name || ''),
                                      variant.source_id,
                                    )}
                                    data-testid={`compression-link-${index}`}
                                  >
                                    <Truncate
                                      content={modelDisplayName}
                                      position="middle"
                                      maxCharsDisplayed={20}
                                      tooltipPosition="top"
                                      style={{ textDecoration: 'underline' }}
                                    />
                                  </Link>
                                )}
                              </FlexItem>
                              <FlexItem>
                                <Flex
                                  spaceItems={{ default: 'spaceItemsXs' }}
                                  alignItems={{ default: 'alignItemsCenter' }}
                                >
                                  <FlexItem>
                                    {tensorType && tensorType !== '-' && (
                                      <Label
                                        isCompact
                                        data-testid={`compression-tensor-type-${index}`}
                                      >
                                        {tensorType}
                                      </Label>
                                    )}
                                  </FlexItem>
                                  {isCurrent && (
                                    <FlexItem>
                                      <Label
                                        isCompact
                                        variant="outline"
                                        data-testid="compression-current-label"
                                      >
                                        Current model
                                      </Label>
                                    </FlexItem>
                                  )}
                                </Flex>
                              </FlexItem>
                            </Flex>
                          </FlexItem>
                        </Flex>
                      </FlexItem>
                    </React.Fragment>
                  );
                })}
              </Flex>
            )}
          </FlexItem>
        </Flex>
      </CardBody>
    </Card>
  );
};

export default TensorTypeComparisonCard;
