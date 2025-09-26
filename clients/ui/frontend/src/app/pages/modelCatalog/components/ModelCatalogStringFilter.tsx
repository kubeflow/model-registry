import { Button, Checkbox, Content, ContentVariants, SearchInput } from '@patternfly/react-core';
import * as React from 'react';
import {
  ModelCatalogFilterDataType,
  ModelCatalogFilterCategoryType,
  ModelCatalogStringFilterStateType,
} from '~/app/pages/modelCatalog/types';

type ModelCatalogStringFilterProps = {
  title: string;
  filterKey: string;
  filterToNameMapping: Record<string, string>;
  filters: ModelCatalogFilterCategoryType;
  data?: ModelCatalogFilterDataType;
  setData: (state: ModelCatalogStringFilterStateType) => void;
};

const ModelCatalogStringFilter: React.FC<ModelCatalogStringFilterProps> = ({
  title,
  filterKey,
  filterToNameMapping,
  data,
  filters,
  setData,
}) => {
  const [showMore, setShowMore] = React.useState(false);
  const [searchValue, setSearchValue] = React.useState('');
  const [filteredValues, setFilteredValues] = React.useState(filters.values);
  const onSearchChange = (newValue: string) => {
    setSearchValue(newValue);
    setFilteredValues(
      filters.values.filter((value) => value.toLowerCase().includes(newValue.toLowerCase())),
    );
  };

  return (
    <Content>
      <Content component={ContentVariants.h6}>{title}</Content>
      {filters.values.length > 5 && (
        <SearchInput
          value={searchValue}
          onChange={(_event, newValue) => onSearchChange(newValue)}
        />
      )}
      {filteredValues.slice(0, 6).map((checkbox) => (
        <Checkbox
          label={checkbox in filterToNameMapping ? filterToNameMapping[checkbox] : checkbox}
          id={checkbox}
          key={checkbox}
          isChecked={data?.[filterKey]?.[checkbox] || false}
          onChange={(_, checked) => {
            if (data?.[filterKey]) {
              setData({ ...data[filterKey], [checkbox]: checked });
            }
          }}
        />
      ))}
      {showMore &&
        filteredValues
          .slice(6)
          .map((checkbox) => <Checkbox label={checkbox} id={checkbox} key={checkbox} />)}
      {filteredValues.length > 5 && (
        <Button variant="link" onClick={() => setShowMore(!showMore)}>
          {showMore ? 'Show less' : 'Show more'}
        </Button>
      )}
    </Content>
  );
};

export default ModelCatalogStringFilter;
