/* eslint-disable camelcase */
import { ModelCatalogStringFilterKey, UseCaseOptionValue } from '~/concepts/modelCatalog/const';
import {
  ModelCatalogFilterStates,
  CatalogPerformanceMetricsArtifact,
  CatalogArtifactType,
  MetricsType,
} from '~/app/modelCatalogTypes';
import { ModelRegistryMetadataType } from '~/app/types';
import { filterHardwareConfigurationArtifacts } from '~/app/pages/modelCatalog/utils/hardwareConfigurationFilterUtils';

describe('hardwareConfigurationFilterUtils', () => {
  describe('filterHardwareConfigurationArtifacts', () => {
    const createMockArtifact = (useCase: string): CatalogPerformanceMetricsArtifact => ({
      artifactType: CatalogArtifactType.metricsArtifact,
      metricsType: MetricsType.performanceMetrics,
      createTimeSinceEpoch: '1739210683000',
      lastUpdateTimeSinceEpoch: '1739210683000',
      customProperties: {
        use_case: {
          metadataType: ModelRegistryMetadataType.STRING,
          string_value: useCase,
        },
        requests_per_second: {
          metadataType: ModelRegistryMetadataType.DOUBLE,
          double_value: 10,
        },
      },
    });

    const createFilterState = (useCases: UseCaseOptionValue[]): ModelCatalogFilterStates => ({
      [ModelCatalogStringFilterKey.TASK]: [],
      [ModelCatalogStringFilterKey.PROVIDER]: [],
      [ModelCatalogStringFilterKey.LICENSE]: [],
      [ModelCatalogStringFilterKey.LANGUAGE]: [],
      [ModelCatalogStringFilterKey.HARDWARE_TYPE]: [],
      [ModelCatalogStringFilterKey.USE_CASE]: useCases,
      rps_mean: undefined,
      ttft_mean: undefined,
    });

    it('should filter artifacts by exact use case match', () => {
      const artifacts = [
        createMockArtifact('rag'),
        createMockArtifact('long_rag'),
        createMockArtifact('chatbot'),
      ];

      // Filter for only 'rag' - should NOT match 'long_rag'
      const ragFilterState = createFilterState([UseCaseOptionValue.RAG]);
      const ragResults = filterHardwareConfigurationArtifacts(artifacts, ragFilterState);

      expect(ragResults).toHaveLength(1);
      expect(ragResults[0].customProperties.use_case?.string_value).toBe('rag');
    });

    it('should filter artifacts for long_rag independently from rag', () => {
      const artifacts = [
        createMockArtifact('rag'),
        createMockArtifact('long_rag'),
        createMockArtifact('chatbot'),
      ];

      // Filter for only 'long_rag' - should NOT match 'rag'
      const longRagFilterState = createFilterState([UseCaseOptionValue.LONG_RAG]);
      const longRagResults = filterHardwareConfigurationArtifacts(artifacts, longRagFilterState);

      expect(longRagResults).toHaveLength(1);
      expect(longRagResults[0].customProperties.use_case?.string_value).toBe('long_rag');
    });

    it('should filter artifacts for multiple use cases', () => {
      const artifacts = [
        createMockArtifact('rag'),
        createMockArtifact('long_rag'),
        createMockArtifact('chatbot'),
        createMockArtifact('code_fixing'),
      ];

      // Filter for both 'rag' and 'long_rag'
      const multiFilterState = createFilterState([
        UseCaseOptionValue.RAG,
        UseCaseOptionValue.LONG_RAG,
      ]);
      const multiResults = filterHardwareConfigurationArtifacts(artifacts, multiFilterState);

      expect(multiResults).toHaveLength(2);
      const useCases = multiResults.map(
        (artifact) => artifact.customProperties.use_case?.string_value,
      );
      expect(useCases).toContain('rag');
      expect(useCases).toContain('long_rag');
      expect(useCases).not.toContain('chatbot');
      expect(useCases).not.toContain('code_fixing');
    });

    it('should return all artifacts when no use case filter is applied', () => {
      const artifacts = [
        createMockArtifact('rag'),
        createMockArtifact('long_rag'),
        createMockArtifact('chatbot'),
      ];

      const noFilterState = createFilterState([]);
      const results = filterHardwareConfigurationArtifacts(artifacts, noFilterState);

      expect(results).toHaveLength(3);
    });

    it('should return empty array when no artifacts match the filter', () => {
      const artifacts = [createMockArtifact('rag'), createMockArtifact('long_rag')];

      const chatbotFilterState = createFilterState([UseCaseOptionValue.CHATBOT]);
      const results = filterHardwareConfigurationArtifacts(artifacts, chatbotFilterState);

      expect(results).toHaveLength(0);
    });

    it('should handle artifacts without use_case property', () => {
      const artifactWithoutUseCase: CatalogPerformanceMetricsArtifact = {
        artifactType: CatalogArtifactType.metricsArtifact,
        metricsType: MetricsType.performanceMetrics,
        createTimeSinceEpoch: '1739210683000',
        lastUpdateTimeSinceEpoch: '1739210683000',
        customProperties: {
          requests_per_second: {
            metadataType: ModelRegistryMetadataType.DOUBLE,
            double_value: 10,
          },
        },
      };

      const artifacts = [artifactWithoutUseCase, createMockArtifact('rag')];

      const ragFilterState = createFilterState([UseCaseOptionValue.RAG]);
      const results = filterHardwareConfigurationArtifacts(artifacts, ragFilterState);

      // Should only return the artifact with matching use_case
      expect(results).toHaveLength(1);
      expect(results[0].customProperties.use_case?.string_value).toBe('rag');
    });

    describe('exact matching verification', () => {
      it('should NOT match rag when filtering for long_rag', () => {
        const artifacts = [createMockArtifact('rag')];
        const longRagFilterState = createFilterState([UseCaseOptionValue.LONG_RAG]);
        const results = filterHardwareConfigurationArtifacts(artifacts, longRagFilterState);

        expect(results).toHaveLength(0);
      });

      it('should NOT match long_rag when filtering for rag', () => {
        const artifacts = [createMockArtifact('long_rag')];
        const ragFilterState = createFilterState([UseCaseOptionValue.RAG]);
        const results = filterHardwareConfigurationArtifacts(artifacts, ragFilterState);

        expect(results).toHaveLength(0);
      });

      it('should match exact string values only', () => {
        const artifacts = [
          createMockArtifact('rag'),
          createMockArtifact('long_rag'),
          createMockArtifact('rag_extended'), // This should not exist in real data, but testing edge case
          createMockArtifact('chatbot'),
        ];

        const ragFilterState = createFilterState([UseCaseOptionValue.RAG]);
        const results = filterHardwareConfigurationArtifacts(artifacts, ragFilterState);

        // Should only match exact 'rag', not 'long_rag' or 'rag_extended'
        expect(results).toHaveLength(1);
        expect(results[0].customProperties.use_case?.string_value).toBe('rag');
      });
    });
  });
});
