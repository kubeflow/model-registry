import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import {
  ModelCatalogFilterTypesByKey,
  ModelCatalogFilterDataType,
  ModelCatalogFilterState,
  FilterValue,
} from '~/app/modelCatalogTypes';
import { ModelCatalogFilterKeys } from '~/concepts/modelCatalog/const';

const MAX_VISIBLE_FILTERS = 5;

type ModelCatalogStringFilterProps<K extends ModelCatalogFilterKeys> = {
  title: string;
  filterToNameMapping: Partial<Record<FilterValue<K>, string>>;
  filters: ModelCatalogFilterTypesByKey[K];
  data?: ModelCatalogFilterDataType[K];
  setData: (state: ModelCatalogFilterState<K>) => void;
};

const ModelCatalogStringFilter = <K extends ModelCatalogFilterKeys>({
  title,
  filterToNameMapping,
  filters,
  data,
  setData,
}: ModelCatalogStringFilterProps<K>): JSX.Element => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const [filteredValues, setFilteredValues] = React.useState<FilterValue<K>[]>(filters.values);

  const getLabel = React.useCallback(
    (value: FilterValue<K>) => filterToNameMapping[value] ?? value,
    [filterToNameMapping],
  );

  React.useEffect(() => {
    const lowerCaseValue = searchValue.trim().toLowerCase();

    if (!lowerCaseValue) {
      setFilteredValues(filters.values);
      return;
    }

    setFilteredValues(
      filters.values.filter((value) => {
        const label = getLabel(value).toLowerCase();

        return value.toLowerCase().includes(lowerCaseValue) || label.includes(lowerCaseValue);
      }),
    );
  }, [filters.values, searchValue, getLabel]);

  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
  };

  const onToggle = (checkbox: FilterValue<K>, checked: boolean) => {
    const nextState: Partial<Record<FilterValue<K>, boolean>> = {
      ...(data ?? {}),
      [checkbox]: checked,
    };
    setData(nextState);
  };

  const isChecked = (value: FilterValue<K>) => Boolean(data?.[value]);

  const visibleValues = showMore ? filteredValues : filteredValues.slice(0, MAX_VISIBLE_FILTERS);

  return (
    <Content>
      <Content component={ContentVariants.h6}>{title}</Content>
      {filters.values.length > MAX_VISIBLE_FILTERS && (
        <SearchInput
          value={searchValue}
          onChange={(_event, newValue) => onSearchChange(newValue)}
        />
      )}
      {visibleValues.map((checkbox) => (
        <Checkbox
          label={getLabel(checkbox)}
          id={checkbox}
          key={checkbox}
          isChecked={isChecked(checkbox)}
          onChange={(_, checked) => onToggle(checkbox, checked)}
        />
      ))}
      {!showMore && filteredValues.length > MAX_VISIBLE_FILTERS && (
        <Button variant="link" onClick={() => setShowMore(true)}>
          Show more
        </Button>
      )}
      {showMore && filteredValues.length > MAX_VISIBLE_FILTERS && (
        <Button variant="link" onClick={() => setShowMore(false)}>
          Show less
        </Button>
      )}
    </Content>
  );
};

export default ModelCatalogStringFilter;
