import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKey,
  MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, GlobalFilterTypes } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogFilterKey.LANGUAGE;

const LANGUAGE_NAME_MAPPING = {
  ...MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
};

type LanguageFilterProps = {
  filters?: Extract<CatalogFilterOptionsList['filters'], Partial<GlobalFilterTypes>>;
};

const LanguageFilter: React.FC<LanguageFilterProps> = ({ filters }) => {
  const language = filters?.[filterKey];

  if (!language) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKey.LANGUAGE>
      title="Language"
      filterKey={filterKey}
      filterToNameMapping={LANGUAGE_NAME_MAPPING}
      filters={language}
    />
  );
};

export default LanguageFilter;
