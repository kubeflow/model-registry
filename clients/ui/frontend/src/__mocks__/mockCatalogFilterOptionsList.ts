import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import {
  ModelCatalogStringFilterKey,
  ModelCatalogNumberFilterKey,
  ModelCatalogLicense,
  ModelCatalogProvider,
  ModelCatalogTask,
  AllLanguageCode,
} from '~/concepts/modelCatalog/const';

export const mockCatalogFilterOptionsList = (
  partial?: Partial<CatalogFilterOptionsList>,
): CatalogFilterOptionsList => ({
  filters: {
    [ModelCatalogStringFilterKey.PROVIDER]: {
      type: 'string',
      values: [ModelCatalogProvider.RED_HAT, ModelCatalogProvider.IBM, ModelCatalogProvider.GOOGLE],
    },
    [ModelCatalogStringFilterKey.LICENSE]: {
      type: 'string',
      values: [ModelCatalogLicense.APACHE_2_0, ModelCatalogLicense.MIT],
    },
    [ModelCatalogStringFilterKey.TASK]: {
      type: 'string',
      values: [
        ModelCatalogTask.TEXT_GENERATION,
        ModelCatalogTask.TEXT_TO_TEXT,
        ModelCatalogTask.IMAGE_TO_TEXT,
        ModelCatalogTask.IMAGE_TEXT_TO_TEXT,
        ModelCatalogTask.VIDEO_TO_TEXT,
        ModelCatalogTask.AUDIO_TO_TEXT,
      ],
    },
    [ModelCatalogStringFilterKey.LANGUAGE]: {
      type: 'string',
      values: [
        AllLanguageCode.AR,
        AllLanguageCode.CS,
        AllLanguageCode.DE,
        AllLanguageCode.EN,
        AllLanguageCode.ES,
        AllLanguageCode.FR,
        AllLanguageCode.IT,
        AllLanguageCode.JA,
        AllLanguageCode.KO,
        AllLanguageCode.NL,
        AllLanguageCode.PT,
        AllLanguageCode.ZH,
      ],
    },
    [ModelCatalogNumberFilterKey.MIN_RPS]: {
      type: 'number',
      range: {
        min: 1,
        max: 300,
      },
    },
    [ModelCatalogNumberFilterKey.MAX_LATENCY]: {
      type: 'number',
      range: {
        min: 20,
        max: 893,
      },
    },
    [ModelCatalogNumberFilterKey.WORKLOAD_TYPE]: {
      type: 'number',
      range: {
        min: 0,
        max: 5,
      },
    },
    [ModelCatalogNumberFilterKey.HARDWARE_TYPE]: {
      type: 'number',
      range: {
        min: 0,
        max: 10,
      },
    },
  },
  ...partial,
});
