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
  FILTERED = 'purple',
  UNFILTERED = 'grey',
}

// Catalog source status values from the API
export enum CatalogSourceStatus {
  AVAILABLE = 'available',
  PARTIALLY_AVAILABLE = 'partially-available',
  ERROR = 'error',
  DISABLED = 'disabled',
}

/**
 * Checks whether a catalog source status indicates that models are available.
 * Sources with 'available' or 'partially-available' status have discoverable models.
 */
export const isSourceStatusWithModels = (status: string | undefined): boolean =>
  status === CatalogSourceStatus.AVAILABLE || status === CatalogSourceStatus.PARTIALLY_AVAILABLE;

// Type guard for Hugging Face sources
export const isHuggingFaceSource = (
  config: CatalogSourceConfig,
): config is HuggingFaceCatalogSourceConfig => config.type === CatalogSourceType.HUGGING_FACE;
