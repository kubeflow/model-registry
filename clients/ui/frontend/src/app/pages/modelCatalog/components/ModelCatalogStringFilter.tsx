import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import { ModelCatalogFilterTypesByKey, FilterValue } from '~/app/modelCatalogTypes';
import { ModelCatalogFilterKeys } from '~/concepts/modelCatalog/const';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

const MAX_VISIBLE_FILTERS = 5;

type ModelCatalogStringFilterProps<K extends ModelCatalogFilterKeys> = {
  title: string;
  filterKey: K;
  filterToNameMapping: Partial<Record<FilterValue<K>, string>>;
  filters: ModelCatalogFilterTypesByKey[K];
};

const ModelCatalogStringFilter = <K extends ModelCatalogFilterKeys>({
  title,
  filterKey,
  filterToNameMapping,
  filters,
}: ModelCatalogStringFilterProps<K>): JSX.Element => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const [filteredValues, setFilteredValues] = React.useState<FilterValue<K>[]>(filters.values);
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);

  const getLabel = React.useCallback(
    (value: FilterValue<K>) => filterToNameMapping[value] ?? value,
    [filterToNameMapping],
  );

  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
    const lowerCaseValue = newValue.trim().toLowerCase();
    setFilteredValues(
      filters.values.filter((value) => {
        const label = getLabel(value).toLowerCase();
        return value.toLowerCase().includes(lowerCaseValue) || label.includes(lowerCaseValue);
      }),
    );
  };

  const onToggle = (checkbox: FilterValue<K>, checked: boolean) => {
    const nextState: Partial<Record<FilterValue<K>, boolean>> = {
      ...(filterData[filterKey] ?? {}),
      [checkbox]: checked,
    };
    setFilterData(filterKey, nextState);
  };

  const isChecked = (value: FilterValue<K>) => filterData[filterKey]?.[value] || false;

  const visibleValues = showMore ? filteredValues : filteredValues.slice(0, MAX_VISIBLE_FILTERS);

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
      {!showMore && filteredValues.length > MAX_VISIBLE_FILTERS && (
        <Button
          variant="link"
          onClick={() => setShowMore(true)}
          data-testid={`${title}-filter-show-more`}
        >
          Show more
        </Button>
      )}
      {showMore && filteredValues.length > MAX_VISIBLE_FILTERS && (
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
