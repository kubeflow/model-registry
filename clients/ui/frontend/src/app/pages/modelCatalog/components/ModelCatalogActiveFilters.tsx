import React from 'react';
import {
  Button,
  Flex,
  FlexItem,
  Label,
  LabelGroup,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { TimesIcon } from '@patternfly/react-icons';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_PROVIDER_NAME_MAPPING,
  MODEL_CATALOG_LICENSE_NAME_MAPPING,
  MODEL_CATALOG_TASK_NAME_MAPPING,
  AllLanguageCodesMap,
} from '~/concepts/modelCatalog/const';
import type {
  ModelCatalogProvider,
  ModelCatalogLicense,
  ModelCatalogTask,
  AllLanguageCode,
} from '~/concepts/modelCatalog/const';
import { hasFiltersApplied } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

type ModelCatalogActiveFiltersProps = {
  searchTerm?: string;
  onResetAllFilters: () => void;
  // eslint-disable-next-line react/no-unused-prop-types
  onClearSearch?: () => void;
};

// Filter category display names
const FILTER_CATEGORY_NAMES: Record<ModelCatalogStringFilterKey, string> = {
  [ModelCatalogStringFilterKey.PROVIDER]: 'Provider',
  [ModelCatalogStringFilterKey.LICENSE]: 'License',
  [ModelCatalogStringFilterKey.TASK]: 'Task',
  [ModelCatalogStringFilterKey.LANGUAGE]: 'Language',
  [ModelCatalogStringFilterKey.HARDWARE_TYPE]: 'Hardware type',
  [ModelCatalogStringFilterKey.USE_CASE]: 'Use case',
};

const ModelCatalogActiveFilters: React.FC<ModelCatalogActiveFiltersProps> = ({
  searchTerm,
  onResetAllFilters,
}) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const filtersApplied = hasFiltersApplied(filterData);
  const hasSearch = Boolean(searchTerm);

  // Don't render if there are no active filters or search term
  if (!filtersApplied && !hasSearch) {
    return null;
  }

  const handleRemoveFilter = (filterKey: ModelCatalogStringFilterKey, value: string) => {
    const currentValues = filterData[filterKey];
    if (Array.isArray(currentValues)) {
      // Compare as strings to handle type mismatches
      const newValues = currentValues.filter((v) => String(v) !== String(value));
      setFilterData(filterKey, newValues);
    }
  };

  const handleClearCategory = (filterKey: ModelCatalogStringFilterKey) => {
    setFilterData(filterKey, []);
  };

  const getFilterLabel = (filterKey: ModelCatalogStringFilterKey, value: string): string => {
    try {
      switch (filterKey) {
        case ModelCatalogStringFilterKey.PROVIDER: {
          // Try to find in mapping, fallback to value
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          const providerValue = value as ModelCatalogProvider;
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          return MODEL_CATALOG_PROVIDER_NAME_MAPPING[providerValue] || value;
        }
        case ModelCatalogStringFilterKey.LICENSE: {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          const licenseValue = value as ModelCatalogLicense;
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          return MODEL_CATALOG_LICENSE_NAME_MAPPING[licenseValue] || value;
        }
        case ModelCatalogStringFilterKey.TASK: {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          const taskValue = value as ModelCatalogTask;
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          return MODEL_CATALOG_TASK_NAME_MAPPING[taskValue] || value;
        }
        case ModelCatalogStringFilterKey.LANGUAGE: {
          // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
          const languageValue = value as AllLanguageCode;
          // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
          return AllLanguageCodesMap[languageValue] || value;
        }
        default:
          return value;
      }
    } catch {
      return value;
    }
  };

  // Get active filters grouped by category
  const activeFiltersByCategory: Array<{
    category: ModelCatalogStringFilterKey;
    categoryName: string;
    values: string[];
  }> = [];

  // Process each filter category
  Object.values(ModelCatalogStringFilterKey).forEach((filterKey) => {
    if (filterKey === ModelCatalogStringFilterKey.USE_CASE) {
      // Skip USE_CASE as it's a single value, not an array
      return;
    }

    const filterValues = filterData[filterKey];
    if (Array.isArray(filterValues) && filterValues.length > 0) {
      activeFiltersByCategory.push({
        category: filterKey,
        categoryName: FILTER_CATEGORY_NAMES[filterKey],
        values: filterValues.map((v) => String(v)),
      });
    }
  });

  // Only show chips section if there are actual chips to display
  const hasChipsToShow = activeFiltersByCategory.length > 0;

  return (
    <Stack>
      {hasChipsToShow && (
        <StackItem>
          <Flex
            gap={{ default: 'gapMd' }}
            alignItems={{ default: 'alignItemsBaseline' }}
            wrap="wrap"
          >
            {activeFiltersByCategory.map(({ category, categoryName, values }) => (
              <FlexItem key={category}>
                <div
                  className="pf-v6-c-toolbar__group"
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: 'var(--pf-v6-global--spacer--xs)',
                    padding: 'var(--pf-v6-global--spacer--xs) var(--pf-v6-global--spacer--sm)',
                    border: '1px solid var(--pf-v6-global--BorderColor--200)',
                    borderRadius: 'var(--pf-v6-global--BorderRadius--sm)',
                    backgroundColor: 'var(--pf-v6-global--BackgroundColor--100)',
                    minHeight: '32px',
                    boxSizing: 'border-box',
                    width: 'fit-content',
                  }}
                  data-testid={`${category}-filter-container`}
                >
                  <LabelGroup categoryName={categoryName}>
                    {values.map((value) => {
                      const labelText = getFilterLabel(category, value);
                      return (
                        <Label
                          key={`${category}-${value}`}
                          variant="outline"
                          color="grey"
                          onClose={() => handleRemoveFilter(category, value)}
                          closeBtnProps={{
                            'aria-label': `Remove ${categoryName} filter ${labelText}`,
                          }}
                          data-testid={`${category}-filter-chip-${value}`}
                          style={{
                            backgroundColor: 'var(--pf-v6-global--BackgroundColor--200)',
                          }}
                        >
                          {labelText}
                        </Label>
                      );
                    })}
                  </LabelGroup>
                  {values.length > 0 && (
                    <Button
                      variant="plain"
                      onClick={() => handleClearCategory(category)}
                      aria-label={`Clear all ${categoryName} filters`}
                      icon={<TimesIcon />}
                      style={{
                        padding: 0,
                        marginLeft: 'var(--pf-v6-global--spacer--xs)',
                        color: 'var(--pf-v6-global--Color--100)',
                        minWidth: 'auto',
                      }}
                    />
                  )}
                </div>
              </FlexItem>
            ))}
          </Flex>
        </StackItem>
      )}
      {(filtersApplied || hasSearch) && (
        <StackItem className="pf-v6-u-pt-sm">
          <Button
            variant="link"
            isInline
            onClick={onResetAllFilters}
            data-testid="reset-all-filters-button"
            style={{ fontSize: 'var(--pf-v6-global--FontSize--sm)', padding: 0 }}
          >
            Reset all filters
          </Button>
        </StackItem>
      )}
    </Stack>
  );
};

export default ModelCatalogActiveFilters;
