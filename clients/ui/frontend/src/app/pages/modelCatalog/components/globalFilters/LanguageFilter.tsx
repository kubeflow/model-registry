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
  const currentState = filterData[filterKey];

  React.useEffect(() => {
    if (!language) {
      return;
    }

    const filterKeys = language.values;
    const hasMatchingKeys =
      currentState !== undefined &&
      filterKeys.length === Object.keys(currentState).length &&
      filterKeys.every((key) => key in currentState);

    if (hasMatchingKeys) {
      return;
    }

    const nextState: ModelCatalogFilterStatesByKey[typeof filterKey] = {};
    filterKeys.forEach((key) => {
      nextState[key] = currentState?.[key] ?? false;
    });

    setFilterData(filterKey, nextState);
  }, [language, currentState, setFilterData]);

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
