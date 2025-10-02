import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import {
  ModelCatalogFilterStates,
  GlobalFilterTypes,
  StringFilterValue,
} from '~/app/modelCatalogTypes';
import { ModelCatalogFilterKey } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const MAX_VISIBLE_FILTERS = 5;

type ModelCatalogStringFilterProps<K extends ModelCatalogFilterKey> = {
  title: string;
  filterKey: K;
  filterToNameMapping: Partial<Record<StringFilterValue<K>, string>>;
  filters: GlobalFilterTypes[K];
};

const ModelCatalogStringFilter = <K extends ModelCatalogFilterKey>({
  title,
  filterKey,
  filterToNameMapping,
  filters,
}: ModelCatalogStringFilterProps<K>): JSX.Element => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  const getLabel = React.useCallback(
    (value: StringFilterValue<K>) => filterToNameMapping[value] ?? value,
    [filterToNameMapping],
  );

  const valuesMatchingSearch = React.useMemo(
    () =>
      filters.values.filter((value) => {
        const label = getLabel(value).toLowerCase();
        return (
          value.toLowerCase().includes(searchValue.trim().toLowerCase()) ||
          label.includes(searchValue.trim().toLowerCase())
        );
      }),
    [filters.values, getLabel, searchValue],
  );

  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
  };

  const onToggle = (checkbox: StringFilterValue<K>, checked: boolean) => {
    const currentState = filterData[filterKey];
    const nextState: StringFilterValue<K>[] = checked
      ? [...currentState, checkbox]
      : currentState.filter((item) => item !== checkbox);
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    setFilterData(filterKey, nextState as ModelCatalogFilterStates[K]);
  };

  const isChecked = (value: StringFilterValue<K>) => {
    // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
    const currentValues = filterData[filterKey] as StringFilterValue<K>[];
    return currentValues.includes(value);
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
      {visibleValues.map((checkbox) => (
        <Checkbox
          data-testid={`${title}-${checkbox}-checkbox`}
          label={getLabel(checkbox)}
          id={checkbox}
          key={checkbox}
          isChecked={isChecked(checkbox)}
          onChange={(_, checked) => onToggle(checkbox, checked)}
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
