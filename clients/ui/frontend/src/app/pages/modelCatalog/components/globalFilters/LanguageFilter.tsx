import * as React from 'react';
import { StackItem } from '@patternfly/react-core';
import ModelCatalogStringFilter from '~/app/pages/modelCatalog/components/ModelCatalogStringFilter';
import {
  ModelCatalogStringFilterKey,
  MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
} from '~/concepts/modelCatalog/const';
import { CatalogFilterOptions, ModelCatalogStringFilterOptions } from '~/app/modelCatalogTypes';

const filterKey = ModelCatalogStringFilterKey.LANGUAGE;

const LANGUAGE_NAME_MAPPING = {
  ...MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
};

type LanguageFilterProps = {
  filters?: Extract<CatalogFilterOptions, Partial<ModelCatalogStringFilterOptions>>;
};

const LanguageFilter: React.FC<LanguageFilterProps> = ({ filters }) => {
  const language = filters?.[filterKey];

  if (!language) {
    return null;
  }

  return (
    <StackItem>
      <ModelCatalogStringFilter<ModelCatalogStringFilterKey.LANGUAGE>
        title="Language"
        filterKey={filterKey}
        filterToNameMapping={LANGUAGE_NAME_MAPPING}
        filters={language}
      />
    </StackItem>
  );
};

export default LanguageFilter;
