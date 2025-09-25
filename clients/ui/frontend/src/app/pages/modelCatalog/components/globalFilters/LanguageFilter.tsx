import * as React from 'react';
import {
  ModelCatalogStringFilterStateType,
  ModelCatalogFilterResponseType,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';

type LanguageFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const LanguageFilter: React.FC<LanguageFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const { language } = filters;

  React.useEffect(() => {
    if (language && !('language' in filterData)) {
      const state: Record<string, boolean> = {};
      language.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData('language', state);
    }
  }, [language, filterData, setFilterData]);

  if (!language) {
    return null;
  }

  return (
    <ModelCatalogStringFilter
      title="Language"
      filterKey="language"
      filters={language}
      data={filterData}
      setData={(state: ModelCatalogStringFilterStateType) => setFilterData('language', state)}
    />
  );
};

export default LanguageFilter;
