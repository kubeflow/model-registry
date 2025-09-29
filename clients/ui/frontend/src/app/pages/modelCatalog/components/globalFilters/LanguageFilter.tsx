import * as React from 'react';
import {
  ModelCatalogFilterResponseType,
  ModelCatalogFilterStatesByKey,
} from '~/app/pages/modelCatalog/types';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
} from '~/concepts/modelCatalog/const';

const filterKey = ModelCatalogFilterKeys.LANGUAGE;

const LANGUAGE_NAME_MAPPING = {
  ...MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
};

type LanguageFilterProps = {
  filters: ModelCatalogFilterResponseType['filters'];
};

const LanguageFilter: React.FC<LanguageFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const language = filters[filterKey];

  React.useEffect(() => {
    if (language && !(filterKey in filterData)) {
      const state: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
      language.values.forEach((key) => {
        state[key] = false;
      });
      setFilterData(filterKey, state);
    }
  }, [language, filterData, setFilterData]);

  if (!language) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.LANGUAGE>
      title="Language"
      filterToNameMapping={LANGUAGE_NAME_MAPPING}
      filters={language}
      data={filterData[filterKey]}
      setData={(state) => setFilterData(filterKey, state)}
    />
  );
};

export default LanguageFilter;
