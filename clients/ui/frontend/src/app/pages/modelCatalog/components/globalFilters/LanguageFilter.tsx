import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import { ModelCatalogContext } from '~/app/context/modelCatalog/ModelCatalogContext';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
} from '~/concepts/modelCatalog/const';
import {
  CatalogFilterOptionsList,
  ModelCatalogFilterTypesByKey,
  ModelCatalogLanguagesFilterStateType,
} from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKeys.LANGUAGE;

const LANGUAGE_NAME_MAPPING = {
  ...MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
};

type LanguageFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<ModelCatalogFilterTypesByKey>>;
};

const LanguageFilter: React.FC<LanguageFilterProps> = ({ filters }) => {
  const { filterData, setFilterData } = React.useContext(ModelCatalogContext);
  const language = filters?.[filterKey];
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

    const nextState: ModelCatalogLanguagesFilterStateType = {};
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
