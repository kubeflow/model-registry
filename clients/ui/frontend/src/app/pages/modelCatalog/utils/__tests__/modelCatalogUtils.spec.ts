// Disabling camelcase for this file because the inherent property names are not camelcase
/* eslint-disable camelcase */
import {
  CatalogFilterOptionsList,
  CatalogSource,
  CatalogSourceList,
  ModelCatalogFilterStates,
} from '~/app/modelCatalogTypes';
import {
  AllLanguageCode,
  ModelCatalogLicense,
  ModelCatalogNumberFilterKey,
  ModelCatalogProvider,
  ModelCatalogStringFilterKey,
  ModelCatalogTask,
  UseCaseOptionValue,
} from '~/concepts/modelCatalog/const';
import { CatalogSourceStatus } from '~/concepts/modelCatalogSettings/const';
import {
  filtersToFilterQuery,
  filterEnabledCatalogSources,
  filterSourcesWithModels,
  hasSourcesWithModels,
  getUniqueSourceLabels,
  hasSourcesWithoutLabels,
} from '~/app/pages/modelCatalog/utils/modelCatalogUtils';

// TODO: Implement performance filters.
describe('filtersToFilterQuery', () => {
  const mockFormData = ({
    tasks = [],
    license = [],
    provider = [],
    language = [],
    hardware_type = [],
    use_case = [],
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
    ttft_mean,
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
      ttft_mean: {
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
      ).toBe("tasks='text-to-text'");
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
      ).toBe("language='ca'");
      expect(
        filtersToFilterQuery(
          // eslint-disable-next-line camelcase
          mockFormData({ use_case: [UseCaseOptionValue.CHATBOT] }),
          mockFilterOptions,
        ),
      ).toBe("use_case='chatbot'");
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
      ).toBe("tasks='text-to-text' AND license='apache-2.0'");
      expect(
        filtersToFilterQuery(
          mockFormData({ provider: [ModelCatalogProvider.GOOGLE], language: [AllLanguageCode.CA] }),
          mockFilterOptions,
        ),
      ).toBe("provider='Google' AND language='ca'");
    });

    it('handles a single array with multiple data points', () => {
      expect(
        filtersToFilterQuery(
          mockFormData({ tasks: [ModelCatalogTask.TEXT_TO_TEXT, ModelCatalogTask.IMAGE_TO_TEXT] }),
          mockFilterOptions,
        ),
      ).toBe("tasks IN ('text-to-text','image-to-text')");
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
      ).toBe("language IN ('ca','pt')");
      expect(
        filtersToFilterQuery(
          // eslint-disable-next-line camelcase
          mockFormData({ use_case: [UseCaseOptionValue.CHATBOT, UseCaseOptionValue.RAG] }),
          mockFilterOptions,
        ),
      ).toBe("use_case IN ('chatbot','rag')");
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
        "tasks IN ('text-to-text','image-to-text') AND provider='Google' AND license='mit' AND language IN ('ca','pt','vi','zsm')",
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

describe('catalog source filtering utilities', () => {
  const createMockSource = (overrides: Partial<CatalogSource> = {}): CatalogSource => ({
    id: 'source-1',
    name: 'Test Source',
    labels: ['Red Hat'],
    enabled: true,
    status: CatalogSourceStatus.AVAILABLE,
    ...overrides,
  });

  const createMockSourceList = (items: CatalogSource[] = []): CatalogSourceList => ({
    items,
    size: items.length,
    pageSize: 10,
    nextPageToken: '',
  });

  describe('filterEnabledCatalogSources', () => {
    it('returns null when catalogSources is null', () => {
      expect(filterEnabledCatalogSources(null)).toBeNull();
    });

    it('returns empty list when no sources are enabled and available', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', enabled: false, status: CatalogSourceStatus.DISABLED }),
        createMockSource({ id: '2', enabled: true, status: CatalogSourceStatus.ERROR }),
      ]);
      const result = filterEnabledCatalogSources(sources);
      expect(result?.items).toHaveLength(0);
    });

    it('filters out disabled sources', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', enabled: true, status: CatalogSourceStatus.AVAILABLE }),
        createMockSource({ id: '2', enabled: false, status: CatalogSourceStatus.AVAILABLE }),
      ]);
      const result = filterEnabledCatalogSources(sources);
      expect(result?.items).toHaveLength(1);
      expect(result?.items?.[0].id).toBe('1');
    });

    it('filters out sources without available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', enabled: true, status: CatalogSourceStatus.AVAILABLE }),
        createMockSource({ id: '2', enabled: true, status: CatalogSourceStatus.ERROR }),
        createMockSource({ id: '3', enabled: true, status: CatalogSourceStatus.DISABLED }),
      ]);
      const result = filterEnabledCatalogSources(sources);
      expect(result?.items).toHaveLength(1);
      expect(result?.items?.[0].id).toBe('1');
    });

    it('returns all sources when all are enabled and available', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', enabled: true, status: CatalogSourceStatus.AVAILABLE }),
        createMockSource({ id: '2', enabled: true, status: CatalogSourceStatus.AVAILABLE }),
      ]);
      const result = filterEnabledCatalogSources(sources);
      expect(result?.items).toHaveLength(2);
    });
  });

  describe('filterSourcesWithModels', () => {
    it('returns null when catalogSources is null', () => {
      expect(filterSourcesWithModels(null)).toBeNull();
    });

    it('returns only sources with available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', status: CatalogSourceStatus.AVAILABLE }),
        createMockSource({ id: '2', status: CatalogSourceStatus.ERROR }),
        createMockSource({ id: '3', status: CatalogSourceStatus.DISABLED }),
      ]);
      const result = filterSourcesWithModels(sources);
      expect(result?.items).toHaveLength(1);
      expect(result?.items?.[0].id).toBe('1');
    });

    it('returns empty list when no sources have available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', status: CatalogSourceStatus.ERROR }),
        createMockSource({ id: '2', status: CatalogSourceStatus.DISABLED }),
      ]);
      const result = filterSourcesWithModels(sources);
      expect(result?.items).toHaveLength(0);
    });
  });

  describe('hasSourcesWithModels', () => {
    it('returns false when catalogSources is null', () => {
      expect(hasSourcesWithModels(null)).toBe(false);
    });

    it('returns false when catalogSources has no items', () => {
      expect(hasSourcesWithModels(createMockSourceList([]))).toBe(false);
    });

    it('returns false when no sources have available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', status: CatalogSourceStatus.ERROR }),
        createMockSource({ id: '2', status: CatalogSourceStatus.DISABLED }),
      ]);
      expect(hasSourcesWithModels(sources)).toBe(false);
    });

    it('returns true when at least one source has available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', status: CatalogSourceStatus.ERROR }),
        createMockSource({ id: '2', status: CatalogSourceStatus.AVAILABLE }),
      ]);
      expect(hasSourcesWithModels(sources)).toBe(true);
    });

    it('returns true when all sources have available status', () => {
      const sources = createMockSourceList([
        createMockSource({ id: '1', status: CatalogSourceStatus.AVAILABLE }),
        createMockSource({ id: '2', status: CatalogSourceStatus.AVAILABLE }),
      ]);
      expect(hasSourcesWithModels(sources)).toBe(true);
    });
  });

  describe('getUniqueSourceLabels', () => {
    it('returns empty array when catalogSources is null', () => {
      expect(getUniqueSourceLabels(null)).toEqual([]);
    });

    it('returns empty array when catalogSources has no items', () => {
      expect(getUniqueSourceLabels(createMockSourceList([]))).toEqual([]);
    });

    it('returns unique labels from enabled and available sources', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat', 'Enterprise'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: ['Red Hat', 'Community'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      const labels = getUniqueSourceLabels(sources);
      expect(labels).toHaveLength(3);
      expect(labels).toContain('Red Hat');
      expect(labels).toContain('Enterprise');
      expect(labels).toContain('Community');
    });

    it('excludes labels from disabled sources', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: ['Excluded'],
          enabled: false,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      const labels = getUniqueSourceLabels(sources);
      expect(labels).toEqual(['Red Hat']);
    });

    it('excludes labels from sources without available status', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: ['Error Source'],
          enabled: true,
          status: CatalogSourceStatus.ERROR,
        }),
      ]);
      const labels = getUniqueSourceLabels(sources);
      expect(labels).toEqual(['Red Hat']);
    });

    it('trims whitespace from labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['  Red Hat  ', 'Enterprise'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      const labels = getUniqueSourceLabels(sources);
      expect(labels).toContain('Red Hat');
    });

    it('excludes empty or whitespace-only labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat', '', '   '],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      const labels = getUniqueSourceLabels(sources);
      expect(labels).toEqual(['Red Hat']);
    });
  });

  describe('hasSourcesWithoutLabels', () => {
    it('returns false when catalogSources is null', () => {
      expect(hasSourcesWithoutLabels(null)).toBe(false);
    });

    it('returns false when catalogSources has no items', () => {
      expect(hasSourcesWithoutLabels(createMockSourceList([]))).toBe(false);
    });

    it('returns true when an enabled and available source has no labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: [],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      expect(hasSourcesWithoutLabels(sources)).toBe(true);
    });

    it('returns true when an enabled and available source has only whitespace labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['   ', ''],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      expect(hasSourcesWithoutLabels(sources)).toBe(true);
    });

    it('returns false when all enabled and available sources have labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: ['Community'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      expect(hasSourcesWithoutLabels(sources)).toBe(false);
    });

    it('ignores disabled sources without labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: [],
          enabled: false,
          status: CatalogSourceStatus.AVAILABLE,
        }),
      ]);
      expect(hasSourcesWithoutLabels(sources)).toBe(false);
    });

    it('ignores sources without available status that have no labels', () => {
      const sources = createMockSourceList([
        createMockSource({
          id: '1',
          labels: ['Red Hat'],
          enabled: true,
          status: CatalogSourceStatus.AVAILABLE,
        }),
        createMockSource({
          id: '2',
          labels: [],
          enabled: true,
          status: CatalogSourceStatus.ERROR,
        }),
      ]);
      expect(hasSourcesWithoutLabels(sources)).toBe(false);
    });
  });
});
