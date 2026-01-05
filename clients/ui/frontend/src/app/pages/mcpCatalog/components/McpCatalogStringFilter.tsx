import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';

const MAX_VISIBLE_FILTERS = 5;

type McpCatalogStringFilterProps = {
  title: string;
  values: string[];
  selectedValues: string[];
  onSelectionChange: (value: string, checked: boolean) => void;
};

const McpCatalogStringFilter: React.FC<McpCatalogStringFilterProps> = ({
  title,
  values,
  selectedValues,
  onSelectionChange,
}) => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');

  const valuesMatchingSearch = React.useMemo(
    () =>
      values.filter((value) => {
        const lowerValue = value.toLowerCase();
        const isSelected = selectedValues.includes(value);
        return lowerValue.includes(searchValue.trim().toLowerCase()) || isSelected;
      }),
    [values, selectedValues, searchValue],
  );

  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
  };

  const visibleValues = showMore
    ? valuesMatchingSearch
    : valuesMatchingSearch.slice(0, MAX_VISIBLE_FILTERS);

  if (values.length === 0) {
    return null;
  }

  return (
    <Content data-testid={`${title}-filter`}>
      <Content component={ContentVariants.h6}>{title}</Content>
      {values.length > MAX_VISIBLE_FILTERS && (
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
          label={checkbox}
          id={`${title}-${checkbox}`}
          key={checkbox}
          isChecked={selectedValues.includes(checkbox)}
          onChange={(_, checked) => onSelectionChange(checkbox, checked)}
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

export default McpCatalogStringFilter;
