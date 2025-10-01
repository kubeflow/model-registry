import { CatalogFilterOptionsList } from '~/app/modelCatalogTypes';
import {
  ModelCatalogFilterKeys,
  ModelCatalogLicenses,
  ModelCatalogProviders,
  ModelCatalogTasks,
  AllLanguageCodes,
} from '~/concepts/modelCatalog/const';

export const mockCatalogFilterOptionsList = (
  partial?: Partial<CatalogFilterOptionsList>,
): CatalogFilterOptionsList => ({
  filters: {
    [ModelCatalogFilterKeys.PROVIDER]: {
      type: 'string',
      values: [
        ModelCatalogProviders.RED_HAT,
        ModelCatalogProviders.IBM,
        ModelCatalogProviders.GOOGLE,
      ],
    },
    [ModelCatalogFilterKeys.LICENSE]: {
      type: 'string',
      values: [ModelCatalogLicenses.APACHE_2_0, ModelCatalogLicenses.MIT],
    },
    [ModelCatalogFilterKeys.TASK]: {
      type: 'string',
      values: [
        ModelCatalogTasks.TEXT_GENERATION,
        ModelCatalogTasks.TEXT_TO_TEXT,
        ModelCatalogTasks.IMAGE_TO_TEXT,
        ModelCatalogTasks.IMAGE_TEXT_TO_TEXT,
        ModelCatalogTasks.VIDEO_TO_TEXT,
        ModelCatalogTasks.AUDIO_TO_TEXT,
      ],
    },
    [ModelCatalogFilterKeys.LANGUAGE]: {
      type: 'string',
      values: [
        AllLanguageCodes.AR,
        AllLanguageCodes.CS,
        AllLanguageCodes.DE,
        AllLanguageCodes.EN,
        AllLanguageCodes.ES,
        AllLanguageCodes.FR,
        AllLanguageCodes.IT,
        AllLanguageCodes.JA,
        AllLanguageCodes.KO,
        AllLanguageCodes.NL,
        AllLanguageCodes.PT,
        AllLanguageCodes.ZH,
      ],
    },
  },
  ...partial,
});
