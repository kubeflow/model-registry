/* eslint-disable camelcase */
import type { CatalogModel } from '~/app/modelCatalogTypes';
import { CatalogModelCustomPropertyKey } from '~/concepts/modelCatalog/const';
import type { ModelRegistryCustomProperties } from '~/app/types';
import { ModelRegistryMetadataType } from '~/app/types';

export type HfAccessModelConfig = {
  name?: string;
  hfAccessType: string;
  hfGatedAccessGranted?: string;
};

export const buildHfAccessCustomProperties = (
  hfAccessType: string,
  hfGatedAccessGranted?: string,
): ModelRegistryCustomProperties => {
  const customProperties: ModelRegistryCustomProperties = {
    [CatalogModelCustomPropertyKey.HF_ACCESS_TYPE]: {
      string_value: hfAccessType,
      metadataType: ModelRegistryMetadataType.STRING,
    },
  };

  if (hfGatedAccessGranted !== undefined) {
    customProperties[CatalogModelCustomPropertyKey.HF_GATED_ACCESS_GRANTED] = {
      string_value: hfGatedAccessGranted,
      metadataType: ModelRegistryMetadataType.STRING,
    };
  }

  return customProperties;
};

export const createHfAccessCatalogModel = ({
  name = 'org/model',
  hfAccessType,
  hfGatedAccessGranted,
}: HfAccessModelConfig): CatalogModel => ({
  name,
  customProperties: buildHfAccessCustomProperties(hfAccessType, hfGatedAccessGranted),
});
