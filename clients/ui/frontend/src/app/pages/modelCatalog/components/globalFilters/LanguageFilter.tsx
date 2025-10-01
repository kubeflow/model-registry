import * as React from 'react';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogFilterKeys,
  MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptionsList, ModelCatalogFilterTypesByKey } from '~/app/modelCatalogTypes';

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
  const language = filters?.[filterKey];

  if (!language) {
    return null;
  }

  return (
    <ModelCatalogStringFilter<ModelCatalogFilterKeys.LANGUAGE>
      title="Language"
      filterKey={filterKey}
      filterToNameMapping={LANGUAGE_NAME_MAPPING}
      filters={language}
    />
  );
};

export default LanguageFilter;
