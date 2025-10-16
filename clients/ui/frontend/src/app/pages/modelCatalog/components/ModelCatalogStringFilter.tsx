import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import {
  ModelCatalogStringFilterOptions,
  ModelCatalogStringFilterValueType,
} from '~/app/modelCatalogTypes';
import { ModelCatalogStringFilterKey } from '~/concepts/modelCatalog/const';
import { useCatalogStringFilterState } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

const MAX_VISIBLE_FILTERS = 5;

type ArrayFilterKey = Exclude<ModelCatalogStringFilterKey, ModelCatalogStringFilterKey.USE_CASE>;

type ModelCatalogStringFilterProps<K extends ArrayFilterKey> = {
  title: string;
  filterKey: K;
  filterToNameMapping: Partial<Record<ModelCatalogStringFilterValueType[K], string>>;
  filters: ModelCatalogStringFilterOptions[K];
};

const ModelCatalogStringFilter = <K extends ArrayFilterKey>({
  title,
  filterKey,
  filterToNameMapping,
  filters,
}: ModelCatalogStringFilterProps<K>): JSX.Element => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const { isSelected, setSelected } = useCatalogStringFilterState(filterKey);

  const getLabel = React.useCallback(
    (value: ModelCatalogStringFilterValueType[K]) => filterToNameMapping[value] ?? value,
    [filterToNameMapping],
  );

  const valuesMatchingSearch = React.useMemo(
    () =>
      filters.values.filter((value) => {
        const label = getLabel(value).toLowerCase();
        return (
          value.toLowerCase().includes(searchValue.trim().toLowerCase()) ||
          label.includes(searchValue.trim().toLowerCase()) ||
          isSelected(value)
        );
      }),
    [filters.values, getLabel, isSelected, searchValue],
  );

  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
  };

  const visibleValues = showMore
    ? valuesMatchingSearch
    : valuesMatchingSearch.slice(0, MAX_VISIBLE_FILTERS);

  return (
    <Content data-testid={`${title}-filter`}>
      <Content component={ContentVariants.h6}>{title}</Content>
      {filters.values.length > MAX_VISIBLE_FILTERS && (
        <SearchInput
          placeholder={`Search ${title.toLowerCase()}`}
          data-testid={`${title}-filter-search`}
          className="pf-v6-u-mb-sm"
          value={searchValue}
          onChange={(_event, newValue) => onSearchChange(newValue)}
        />
      )}
      {visibleValues.length === 0 && (
        <div data-testid={`${title}-filter-empty`}>No results found</div>
      )}
      {visibleValues.map((checkbox) => (
        <Checkbox
          data-testid={`${title}-${checkbox}-checkbox`}
          label={getLabel(checkbox)}
          id={checkbox}
          key={checkbox}
          isChecked={isSelected(checkbox)}
          onChange={(_, checked) => setSelected(checkbox, checked)}
        />
      ))}
      {!showMore && valuesMatchingSearch.length > MAX_VISIBLE_FILTERS && (
        <Button
          variant="link"
          onClick={() => setShowMore(true)}
          data-testid={`${title}-filter-show-more`}
        >
          Show more
        </Button>
      )}
      {showMore && valuesMatchingSearch.length > MAX_VISIBLE_FILTERS && (
        <Button
          variant="link"
          onClick={() => setShowMore(false)}
          data-testid={`${title}-filter-show-less`}
        >
          Show less
        </Button>
      )}
    </Content>
  );
};

export default ModelCatalogStringFilter;
