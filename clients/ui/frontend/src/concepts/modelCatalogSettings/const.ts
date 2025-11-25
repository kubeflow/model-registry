import {
  CatalogSourceConfig,
  HuggingFaceCatalogSourceConfig,
  CatalogSourceType,
} from '~/app/modelCatalogTypes';

export const CATALOG_SOURCE_TYPE_LABELS: Record<CatalogSourceType, string> = {
  [CatalogSourceType.YAML]: 'YAML file',
  [CatalogSourceType.HUGGING_FACE]: 'Hugging Face',
};

export enum ModelVisibilityBadgeColor {
  FILTERED = 'blue',
  UNFILTERED = 'grey',
}

// Catalog source status values from the API
export enum CatalogSourceStatus {
  AVAILABLE = 'available',
  ERROR = 'error',
  DISABLED = 'disabled',
}

// Maximum length for error message display in table before truncation
export const ERROR_MESSAGE_TRUNCATE_LENGTH = 50;

// Type guard for Hugging Face sources
export const isHuggingFaceSource = (
  config: CatalogSourceConfig,
): config is HuggingFaceCatalogSourceConfig => config.type === CatalogSourceType.HUGGING_FACE;
