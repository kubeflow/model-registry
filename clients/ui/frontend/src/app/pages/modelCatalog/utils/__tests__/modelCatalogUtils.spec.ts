// Disabling camelcase for this file because the inherent property names are not camelcase
/* eslint-disable camelcase */
import { CatalogFilterOptionsList, ModelCatalogFilterStates } from '~/app/modelCatalogTypes';
import {
  AllLanguageCode,
  ModelCatalogLicense,
  ModelCatalogNumberFilterKey,
  ModelCatalogProvider,
  ModelCatalogStringFilterKey,
  ModelCatalogTask,
  UseCaseOptionValue,
} from '~/concepts/modelCatalog/const';
import { filtersToFilterQuery } from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

// TODO: Implement performance filters.
describe('filtersToFilterQuery', () => {
  const mockFormData = ({
    tasks = [],
    license = [],
    provider = [],
    language = [],
    hardware_type = [],
    use_case = undefined,
    rps_mean = undefined,
    ttft_mean = undefined,
  }: Partial<ModelCatalogFilterStates>): ModelCatalogFilterStates => ({
    [ModelCatalogStringFilterKey.TASK]: tasks,
    [ModelCatalogStringFilterKey.PROVIDER]: provider,
    [ModelCatalogStringFilterKey.LICENSE]: license,
    [ModelCatalogStringFilterKey.LANGUAGE]: language,
    [ModelCatalogStringFilterKey.HARDWARE_TYPE]: hardware_type,
    [ModelCatalogStringFilterKey.USE_CASE]: use_case,
    [ModelCatalogNumberFilterKey.MIN_RPS]: rps_mean,
    [ModelCatalogNumberFilterKey.MAX_LATENCY]: ttft_mean,
  });

  const mockFilterOptions: CatalogFilterOptionsList = {
    filters: {
      [ModelCatalogStringFilterKey.TASK]: {
        type: 'string',
        values: [
          ModelCatalogTask.AUDIO_TO_TEXT,
          ModelCatalogTask.IMAGE_TEXT_TO_TEXT,
          ModelCatalogTask.IMAGE_TO_TEXT,
          ModelCatalogTask.TEXT_GENERATION,
          ModelCatalogTask.TEXT_TO_TEXT,
          ModelCatalogTask.VIDEO_TO_TEXT,
        ],
      },
      [ModelCatalogStringFilterKey.PROVIDER]: {
        type: 'string',
        values: [
          ModelCatalogProvider.ALIBABA_CLOUD,
          ModelCatalogProvider.DEEPSEEK,
          ModelCatalogProvider.GOOGLE,
          ModelCatalogProvider.IBM,
          ModelCatalogProvider.META,
          ModelCatalogProvider.MISTRAL_AI,
          ModelCatalogProvider.MOONSHOT_AI,
          ModelCatalogProvider.NEURAL_MAGIC,
          ModelCatalogProvider.NVIDIA,
          ModelCatalogProvider.NVIDIA_ALTERNATE,
          ModelCatalogProvider.RED_HAT,
        ],
      },
      [ModelCatalogStringFilterKey.LICENSE]: {
        type: 'string',
        values: [
          ModelCatalogLicense.APACHE_2_0,
          ModelCatalogLicense.GEMMA,
          ModelCatalogLicense.LLLAMA_3_3,
          ModelCatalogLicense.LLLAMA_3_1,
          ModelCatalogLicense.LLLAMA_3_3_ALTERNATE,
          ModelCatalogLicense.LLLAMA_4,
          ModelCatalogLicense.MIT,
          ModelCatalogLicense.MODIFIED_MIT,
        ],
      },
      [ModelCatalogStringFilterKey.LANGUAGE]: {
        type: 'string',
        values: [
          AllLanguageCode.BG,
          AllLanguageCode.CA,
          AllLanguageCode.CS,
          AllLanguageCode.DA,
          AllLanguageCode.DE,
          AllLanguageCode.EL,
          AllLanguageCode.EN,
          AllLanguageCode.ES,
          AllLanguageCode.FI,
          AllLanguageCode.FR,
          AllLanguageCode.HR,
          AllLanguageCode.HU,
          AllLanguageCode.IS,
          AllLanguageCode.IT,
          AllLanguageCode.NL,
          AllLanguageCode.NLD,
          AllLanguageCode.NO,
          AllLanguageCode.PL,
          AllLanguageCode.PT,
          AllLanguageCode.RO,
          AllLanguageCode.RU,
          AllLanguageCode.SK,
          AllLanguageCode.SL,
          AllLanguageCode.SR,
          AllLanguageCode.SV,
          AllLanguageCode.UK,
          AllLanguageCode.JA,
          AllLanguageCode.KO,
          AllLanguageCode.ZH,
          AllLanguageCode.HI,
          AllLanguageCode.TH,
          AllLanguageCode.VI,
          AllLanguageCode.ID,
          AllLanguageCode.MS,
          AllLanguageCode.ZSM,
          AllLanguageCode.AR,
          AllLanguageCode.FA,
          AllLanguageCode.HE,
          AllLanguageCode.TR,
          AllLanguageCode.UR,
          AllLanguageCode.TL,
        ],
      },
      [ModelCatalogStringFilterKey.HARDWARE_TYPE]: {
        type: 'string',
        values: ['GPU', 'CPU', 'TPU', 'FPGA'],
      },
      [ModelCatalogStringFilterKey.USE_CASE]: {
        type: 'string',
        values: [
          UseCaseOptionValue.CHATBOT,
          UseCaseOptionValue.CODE_FIXING,
          UseCaseOptionValue.LONG_RAG,
          UseCaseOptionValue.RAG,
        ],
      },
      [ModelCatalogNumberFilterKey.MIN_RPS]: {
        type: 'number',
        range: {
          min: 0,
          max: 300,
        },
      },
      [ModelCatalogNumberFilterKey.MAX_LATENCY]: {
        type: 'number',
        range: {
          min: 0,
          max: 1000,
        },
      },
    },
  };
  /* eslint-enable camelcase */

  describe('multi-selection values', () => {
    it('handles no data', () => {
      expect(filtersToFilterQuery(mockFormData({}), mockFilterOptions)).toBe('');
    });

    it('handles a single array of a single data point', () => {
      expect(
        filtersToFilterQuery(
          mockFormData({ tasks: [ModelCatalogTask.TEXT_TO_TEXT] }),
          mockFilterOptions,
        ),
      ).toBe('tasks LIKE \'%"text-to-text"%\'');
      expect(
        filtersToFilterQuery(
          mockFormData({ provider: [ModelCatalogProvider.GOOGLE] }),
          mockFilterOptions,
        ),
      ).toBe("provider='Google'");
      expect(
        filtersToFilterQuery(
          mockFormData({ license: [ModelCatalogLicense.APACHE_2_0] }),
          mockFilterOptions,
        ),
      ).toBe("license='apache-2.0'");
      expect(
        filtersToFilterQuery(mockFormData({ language: [AllLanguageCode.CA] }), mockFilterOptions),
      ).toBe('language LIKE \'%"ca"%\'');
    });

    it('handles multiple arrays of a single data point', () => {
      expect(
        filtersToFilterQuery(
          mockFormData({
            tasks: [ModelCatalogTask.TEXT_TO_TEXT],
            license: [ModelCatalogLicense.APACHE_2_0],
          }),
          mockFilterOptions,
        ),
      ).toBe("tasks LIKE '%\"text-to-text\"%' AND license='apache-2.0'");
      expect(
        filtersToFilterQuery(
          mockFormData({ provider: [ModelCatalogProvider.GOOGLE], language: [AllLanguageCode.CA] }),
          mockFilterOptions,
        ),
      ).toBe("provider='Google' AND language LIKE '%\"ca\"%'");
    });

    it('handles a single array with multiple data points', () => {
      expect(
        filtersToFilterQuery(
          mockFormData({ tasks: [ModelCatalogTask.TEXT_TO_TEXT, ModelCatalogTask.IMAGE_TO_TEXT] }),
          mockFilterOptions,
        ),
      ).toBe('(tasks LIKE \'%"text-to-text"%\' OR tasks LIKE \'%"image-to-text"%\')');
      expect(
        filtersToFilterQuery(
          mockFormData({ provider: [ModelCatalogProvider.GOOGLE, ModelCatalogProvider.DEEPSEEK] }),
          mockFilterOptions,
        ),
      ).toBe("provider IN ('Google','DeepSeek')");
      expect(
        filtersToFilterQuery(
          mockFormData({ license: [ModelCatalogLicense.APACHE_2_0, ModelCatalogLicense.MIT] }),
          mockFilterOptions,
        ),
      ).toBe("license IN ('apache-2.0','mit')");
      expect(
        filtersToFilterQuery(
          mockFormData({ language: [AllLanguageCode.CA, AllLanguageCode.PT] }),
          mockFilterOptions,
        ),
      ).toBe('(language LIKE \'%"ca"%\' OR language LIKE \'%"pt"%\')');
    });

    it('handles multiple arrays with mixed count of data points', () => {
      expect(
        filtersToFilterQuery(
          mockFormData({
            tasks: [ModelCatalogTask.TEXT_TO_TEXT, ModelCatalogTask.IMAGE_TO_TEXT],
            provider: [ModelCatalogProvider.GOOGLE],
            license: [ModelCatalogLicense.MIT],
            language: [
              AllLanguageCode.CA,
              AllLanguageCode.PT,
              AllLanguageCode.VI,
              AllLanguageCode.ZSM,
            ],
          }),
          mockFilterOptions,
        ),
      ).toBe(
        "(tasks LIKE '%\"text-to-text\"%' OR tasks LIKE '%\"image-to-text\"%') AND provider='Google' AND license='mit' AND (language LIKE '%\"ca\"%' OR language LIKE '%\"pt\"%' OR language LIKE '%\"vi\"%' OR language LIKE '%\"zsm\"%')",
      );
    });
  });

  // TODO: Implement performance filters.
  //   describe('less than values', () => {
  //     it('handles TimeToFirstToken - ttft', () => {
  //       // eslint-disable-next-line camelcase
  //       expect(filtersToFilterQuery(mockFormData({ ttft_mean: 100 }), mockFilterOptions)).toBe(
  //         'ttft_mean < 100',
  //       );
  //     });
  //   });

  //   describe('greater than values', () => {
  //     it('handles TimeToFirstToken - ttft', () => {
  //       // eslint-disable-next-line camelcase
  //       expect(filtersToFilterQuery(mockFormData({ rps_mean: 7 }), mockFilterOptions)).toBe(
  //         'rps_mean > 7',
  //       );
  //     });
  //   });

  //   describe('mixture of multiple types of values', () => {
  //     it('handles TimeToFirstToken - ttft', () => {
  //       expect(
  //         filtersToFilterQuery(
  //           mockFormData({
  //             // eslint-disable-next-line camelcase
  //             ttft_mean: 100,
  //             tasks: [ModelCatalogTask.TEXT_TO_TEXT],
  //             license: [ModelCatalogLicense.APACHE_2_0, ModelCatalogLicense.MIT],
  //             // eslint-disable-next-line camelcase
  //             rps_mean: 3,
  //           }),
  //           mockFilterOptions,
  //         ),
  //       ).toBe(
  //         "tasks='text-to-text' AND license IN ('apache-2.0','mit') AND ttft_mean < 100 AND rps_mean > 3",
  //       );
  //     });
  //   });
});
